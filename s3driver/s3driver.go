package s3driver

import (
	"errors"
	"io"
	"net/http"
	"os"
	"s3hash-go/driver"
	"s3hash-go/pkg/fpath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// DefaultTimeout ...
const DefaultTimeout time.Duration = 30 * time.Second

// S3Driver ...
type S3Driver struct {
	MaxRetries      int
	Profile         string
	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string
	Region          string
	Debug           bool
	Timeout         time.Duration
}

// NewDriver ...
func (driver *S3Driver) NewDriver() iodriver.Driver {
	return &S3Driver{
		MaxRetries: 3,
		Profile:    driver.Profile,
		Region:     driver.Region,
		Debug:      driver.Debug,
		Timeout:    DefaultTimeout,
	}
}

// Open ...
func (driver *S3Driver) Open(path string) (io.ReadCloser, error) {
	var err error
	bucket, key := fpath.SplitPath(path)
	svc, err := driver.newClientWithBucket(bucket)
	if err != nil {
		return nil, os.ErrNotExist
	}

	u, err := NewDownloader(svc, bucket, key)
	if err != nil {
		return nil, err
	}

	return u, nil
}

// Copy copies from src to dst until either EOF is reached
// on src or an error occurs. It returns the number of bytes
// copied and the first error encountered while copying, if any.
//
// A successful Copy returns err == nil, not err == EOF.
// Because Copy is defined to read from src until EOF, it does
// not treat an EOF from Read as an error to be reported.
//
// If src implements the WriterTo interface,
// the copy is implemented by calling src.WriteTo(dst).
// Otherwise, if dst implements the ReaderFrom interface,
// the copy is implemented by calling dst.ReadFrom(src).
func (driver *S3Driver) Copy(dst io.Writer, src io.Reader) (written int64, err error) {
	return copyBuffer(dst, src, nil)
}

// CopyBuffer is identical to Copy except that it stages through the
// provided buffer (if one is required) rather than allocating a
// temporary one. If buf is nil, one is allocated; otherwise if it has
// zero length, CopyBuffer panics.
func (driver *S3Driver) CopyBuffer(dst io.Writer, src io.Reader, buf []byte) (written int64, err error) {
	if buf != nil && len(buf) == 0 {
		panic("empty buffer in io.CopyBuffer")
	}
	return copyBuffer(dst, src, buf)
}

// copyBuffer is the actual implementation of Copy and CopyBuffer.
// if buf is nil, one is allocated.
func copyBuffer(dst io.Writer, src io.Reader, buf []byte) (written int64, err error) {
	// If the reader has a WriteTo method, use it to do the copy.
	// Avoids an allocation and a copy.
	// if wt, ok := src.(io.WriterTo); ok {
	// 	return wt.WriteTo(dst)
	// }

	// Similarly, if the writer has a ReadFrom method, use it to do the copy.
	// if rt, ok := dst.(io.ReaderFrom); ok {
	// 	return rt.ReadFrom(src)
	// }

	if buf == nil {
		buf = make([]byte, 32*1024)
		//buf = make([]byte, 5*1024*1024)
	}

	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er == io.EOF {
			break
		}
		if er != nil {
			err = er
			break
		}
	}
	return written, err
}

func normalizeError(err error) error {
	s := err.Error()
	return errors.New(strings.Replace(s, "\n", "", -1))
}

func (driver *S3Driver) newClient() (*s3.S3, error) {
	level := aws.LogOff
	if driver.Debug {
		level = aws.LogDebug
	}

	if driver.AccessKeyID == "" {
		sess, err := session.NewSessionWithOptions(session.Options{
			Profile:           driver.Profile,
			SharedConfigState: session.SharedConfigEnable,
		})
		if err != nil {
			return nil, err
		}

		cfg := aws.NewConfig().
			WithLogLevel(level).
			WithMaxRetries(driver.MaxRetries).
			WithHTTPClient(&http.Client{Timeout: driver.Timeout})

		return s3.New(sess, cfg), nil
	}

	cfg := aws.NewConfig().
		WithCredentials(credentials.NewStaticCredentials(driver.AccessKeyID, driver.SecretAccessKey, driver.SessionToken)).
		WithLogLevel(level).
		WithRegion(driver.Region).
		WithMaxRetries(driver.MaxRetries).
		WithHTTPClient(&http.Client{Timeout: driver.Timeout})

	return s3.New(session.New(), cfg), nil
}

func (driver *S3Driver) newClientWithBucket(bucket string) (*s3.S3, error) {
	svc, err := driver.newClient()
	if err != nil {
		return nil, err
	}

	req, result := svc.GetBucketLocationRequest(&s3.GetBucketLocationInput{
		Bucket: aws.String(bucket),
	})
	req.Handlers.Unmarshal.PushBackNamed(s3.NormalizeBucketLocationHandler)
	err = req.Send()

	// result, err := svc.GetBucketLocation(&s3.GetBucketLocationInput{Bucket: aws.String(bucket)})
	if err != nil {
		return nil, err
	}

	var creds *credentials.Credentials
	if driver.AccessKeyID == "" {
		creds = credentials.NewSharedCredentials("", driver.Profile)
	} else {
		creds = credentials.NewStaticCredentials(driver.AccessKeyID, driver.SecretAccessKey, driver.SessionToken)
	}

	level := aws.LogOff
	if driver.Debug {
		level = aws.LogDebug
	}

	cfg := aws.NewConfig().
		WithCredentials(creds).
		WithLogLevel(level).
		//WithRegion(normalizeRegionValue(result.LocationConstraint)).
		WithRegion(aws.StringValue(result.LocationConstraint)).
		WithMaxRetries(driver.MaxRetries).
		WithHTTPClient(&http.Client{Timeout: driver.Timeout})

	svc = s3.New(session.New(), cfg)
	acc, err := svc.GetBucketAccelerateConfiguration(&s3.GetBucketAccelerateConfigurationInput{Bucket: aws.String(bucket)})
	if err != nil {
		return nil, err
	}

	cfg = aws.NewConfig().
		WithCredentials(creds).
		WithLogLevel(level).
		//WithRegion(normalizeRegionValue(result.LocationConstraint)).
		WithRegion(aws.StringValue(result.LocationConstraint)).
		WithMaxRetries(driver.MaxRetries).
		WithHTTPClient(&http.Client{Timeout: driver.Timeout}).
		WithS3UseAccelerate(aws.StringValue(acc.Status) == s3.BucketAccelerateStatusEnabled)

	return s3.New(session.New(), cfg), nil
}
