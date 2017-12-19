package s3driver

import (
	"context"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

// DefaultDownloadConcurrency ...
const DefaultDownloadConcurrency = 5

// DefaultDownloadPartSize ...
const DefaultDownloadPartSize = 1024 * 1024 * 5

// DefaultReadTimeout ...
const DefaultReadTimeout time.Duration = 30 * time.Second

// Downloader ...
type Downloader struct {
	ctx         aws.Context
	cancel      context.CancelFunc
	S3          s3iface.S3API
	Bucket      string
	Key         string
	PartSize    int64
	Concurrency int
	Timeout     time.Duration

	id int64

	wg  sync.WaitGroup
	m   sync.Mutex
	err error

	pos        int64
	totalBytes int64
	readBytes  int64

	partBodyMaxRetries int

	ch      chan int64
	readBuf []byte
	offset  int
	length  int
	queue   chan struct{}
	done    chan struct{}
}

// NewDownloader ...
func NewDownloader(svc s3iface.S3API, bucket, key string, options ...func(*Downloader)) (*Downloader, error) {
	return NewDownloaderWithContext(aws.BackgroundContext(), svc, bucket, key, options...)
}

// NewDownloaderWithContext ...
func NewDownloaderWithContext(ctx context.Context, svc s3iface.S3API, bucket, key string, options ...func(*Downloader)) (*Downloader, error) {
	ctx, cancel := context.WithCancel(ctx)
	d := &Downloader{
		ctx:                ctx,
		cancel:             cancel,
		S3:                 svc,
		Bucket:             bucket,
		Key:                key,
		PartSize:           DefaultDownloadPartSize,
		Concurrency:        DefaultDownloadConcurrency,
		Timeout:            DefaultReadTimeout,
		id:                 0,
		readBytes:          0,
		partBodyMaxRetries: 3,
		readBuf:            make([]byte, DefaultDownloadPartSize),
		offset:             0,
		length:             0,
		queue:              make(chan struct{}),
		done:               make(chan struct{}),
	}

	for _, option := range options {
		option(d)
	}

	d.ch = make(chan int64, d.Concurrency)

	output, err := svc.HeadObjectWithContext(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}

	contentLength := aws.Int64Value(output.ContentLength)
	d.totalBytes = contentLength

	for i := 0; i < d.Concurrency; i++ {
		d.wg.Add(1)
		go d.downloadPart(d.ch)
	}

	d.wg.Add(1)
	go d.queuingChunks(contentLength / d.PartSize)

	return d, nil
}

// Read ...
func (d *Downloader) Read(p []byte) (int, error) {
	if err := d.geterr(); err != nil {
		return 0, err
	}

	if d.offset == 0 {
		select {
		case <-d.queue:
		case <-time.After(d.Timeout):
			d.cancel()
			return 0, io.ErrNoProgress
		}
	}

	n := copy(p, d.readBuf[d.offset:d.length])
	d.offset += n
	d.readBytes += int64(n)

	if d.offset >= d.length {
		d.offset = 0
		d.id++
	}

	if d.readBytes >= d.totalBytes {
		d.seterr(io.EOF)
	}

	return n, nil
}

// Close ...
func (d *Downloader) Close() error {
	//fmt.Println("----->called Close start")
	// Wait for completion
	d.cancel()
	close(d.done)
	close(d.queue)
	d.wg.Wait()
	//fmt.Println("----->called Close finish")
	return nil
}

func (d *Downloader) queuingChunks(total int64) {
	defer d.wg.Done()

	var n int64
	for n <= total {
		select {
		case <-d.done:
			return
		case d.ch <- n:
			n++
		default:
			time.Sleep(10 * time.Millisecond)
		}
	}
	close(d.ch)
}

func (d *Downloader) downloadPart(ch chan int64) {
	defer d.wg.Done()

	partSize := d.PartSize
	partBuf := make([]byte, partSize)

Loop:
	for {
		id, ok := <-ch
		if !ok || d.geterr() != nil {
			break
		}

		chunk := &dlchunk{buf: partBuf, start: id * partSize, size: partSize}
		n, err := d.downloadChunk(chunk)
		if err != nil {
			d.seterr(err)
			break
		}

		for {
			select {
			case <-d.done:
				return
			default:
				if id == d.id {
					copy(d.readBuf, partBuf)
					d.length = n
					d.queue <- struct{}{}
					continue Loop
				}
				time.Sleep(10 * time.Millisecond)
			}
		}
	}
}

func (d *Downloader) downloadChunk(chunk *dlchunk) (int, error) {
	rng := fmt.Sprintf("bytes=%d-%d", chunk.start, chunk.start+chunk.size-1)
	in := &s3.GetObjectInput{
		Bucket: aws.String(d.Bucket),
		Key:    aws.String(d.Key),
		Range:  aws.String(rng),
	}

	var n int
	//var err error
	for retry := 0; retry <= d.partBodyMaxRetries; retry++ {
		out, err := d.S3.GetObjectWithContext(d.ctx, in)
		if err != nil {
			return 0, err
		}

		d.setTotalBytes(out) // Set total if not yet set.

		n, err = io.ReadFull(out.Body, chunk.buf)
		out.Body.Close()
		if err == nil || err == io.ErrUnexpectedEOF {
			break
		}
	}

	//d.incrWritten(int64(n))
	return n, nil
}

func (d *Downloader) setTotalBytes(resp *s3.GetObjectOutput) {
	d.m.Lock()
	defer d.m.Unlock()

	if d.totalBytes >= 0 {
		return
	}

	if resp.ContentRange == nil {
		// ContentRange is nil when the full file contents is provied, and
		// is not chunked. Use ContentLength instead.
		if resp.ContentLength != nil {
			d.totalBytes = *resp.ContentLength
			return
		}
	} else {
		parts := strings.Split(*resp.ContentRange, "/")

		total := int64(-1)
		var err error
		// Checking for whether or not a numbered total exists
		// If one does not exist, we will assume the total to be -1, undefined,
		// and sequentially download each chunk until hitting a 416 error
		totalStr := parts[len(parts)-1]
		if totalStr != "*" {
			total, err = strconv.ParseInt(totalStr, 10, 64)
			if err != nil {
				d.err = err
				return
			}
		}

		d.totalBytes = total
	}
}

func (d *Downloader) getTotalBytes() int64 {
	d.m.Lock()
	defer d.m.Unlock()

	return d.totalBytes
}

func (d *Downloader) geterr() error {
	d.m.Lock()
	defer d.m.Unlock()

	return d.err
}

func (d *Downloader) seterr(e error) {
	d.m.Lock()
	defer d.m.Unlock()

	d.err = e
}

type dlchunk struct {
	start int64
	size  int64
	buf   []byte
}
