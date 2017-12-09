package main

import (
	"crypto"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"hash"
	"io"
	"os"
	"s3hash-go/driver"
	"s3hash-go/s3driver"
	"strconv"
	"time"

	"github.com/codegangsta/cli"
)

// HashInfo ...
type HashInfo struct {
	DateTime time.Time `json:"datetime"`
	Path     string    `json:"path"`
	Hash     string    `json:"hash"`
	Binary   string    `json:"binary"`
	Base64   string    `json:"base64"`
	Seconds  string    `json:"seconds"`
}

var debug bool
var input string
var filename string

func main() {
	debug = false
	app := cli.NewApp()
	app.Name = "compute hash"
	app.Usage = "usage compute hash"
	app.Version = "0.1.0"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "debug",
			Usage:       "enable debugging",
			Destination: &debug,
		},
	}

	app.Action = func(c *cli.Context) {
		//log.Println("called app.Action")
	}

	app.Commands = []cli.Command{
		{
			Name:    "md5",
			Aliases: []string{"m"},
			Usage:   "compute hash md5",
			Action:  cmdMd5,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "input",
					Usage:       "target file",
					Destination: &input,
				},
				cli.StringFlag{
					Name:        "output",
					Usage:       "output json file",
					Destination: &filename,
				},
			},
		},
		{
			Name:    "sha1",
			Aliases: []string{"s"},
			Usage:   "compute hash sha1",
			Action:  cmdSha1,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "input",
					Usage:       "target file",
					Destination: &input,
				},
				cli.StringFlag{
					Name:        "output",
					Usage:       "output json file",
					Destination: &filename,
				},
			},
		},
		{
			Name:    "sha224",
			Aliases: []string{"s"},
			Usage:   "compute hash sha224",
			Action:  cmdSha224,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "input",
					Usage:       "target file",
					Destination: &input,
				},
				cli.StringFlag{
					Name:        "output",
					Usage:       "output json file",
					Destination: &filename,
				},
			},
		},
		{
			Name:    "sha256",
			Aliases: []string{"s"},
			Usage:   "compute hash sha256",
			Action:  cmdSha256,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "input",
					Usage:       "target file",
					Destination: &input,
				},
				cli.StringFlag{
					Name:        "output",
					Usage:       "output json file",
					Destination: &filename,
				},
			},
		},
		{
			Name:    "sha384",
			Aliases: []string{"s"},
			Usage:   "compute hash sha384",
			Action:  cmdSha384,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "input",
					Usage:       "target file",
					Destination: &input,
				},
				cli.StringFlag{
					Name:        "output",
					Usage:       "output json file",
					Destination: &filename,
				},
			},
		},
		{
			Name:    "sha512",
			Aliases: []string{"s"},
			Usage:   "compute hash sha512",
			Action:  cmdSha512,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "input",
					Usage:       "target file",
					Destination: &input,
				},
				cli.StringFlag{
					Name:        "output",
					Usage:       "output json file",
					Destination: &filename,
				},
			},
		},
		{
			Name:    "sha512_224",
			Aliases: []string{"s"},
			Usage:   "compute hash sha512_224",
			Action:  cmdSha512_224,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "input",
					Usage:       "target file",
					Destination: &input,
				},
				cli.StringFlag{
					Name:        "output",
					Usage:       "output json file",
					Destination: &filename,
				},
			},
		},
		{
			Name:    "sha512_256",
			Aliases: []string{"s"},
			Usage:   "compute hash sha512_256",
			Action:  cmdSha512_256,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "input",
					Usage:       "target file",
					Destination: &input,
				},
				cli.StringFlag{
					Name:        "output",
					Usage:       "output json file",
					Destination: &filename,
				},
			},
		},
	}

	app.Run(os.Args)
}

func cmdMd5(c *cli.Context) {
	data, err := start(crypto.MD5, md5.New(), input)
	if err != nil {
		fmt.Println(err)
		return
	}

	//ioutil.WriteFile(filename, data, 0644)
	if filename != "" {
		writeFile(filename, data)
	} else {
		fmt.Println(string(data))
	}
}

func cmdSha1(c *cli.Context) {
	data, err := start(crypto.SHA1, sha1.New(), input)
	if err != nil {
		fmt.Println(err)
		return
	}

	if filename != "" {
		writeFile(filename, data)
	} else {
		fmt.Println(string(data))
	}
}

func cmdSha224(c *cli.Context) {
	data, err := start(crypto.SHA224, sha256.New224(), input)
	if err != nil {
		fmt.Println(err)
		return
	}

	if filename != "" {
		writeFile(filename, data)
	} else {
		fmt.Println(string(data))
	}
}

func cmdSha256(c *cli.Context) {
	data, err := start(crypto.SHA256, sha256.New(), input)
	if err != nil {
		fmt.Println(err)
		return
	}

	if filename != "" {
		writeFile(filename, data)
	} else {
		fmt.Println(string(data))
	}
}

func cmdSha384(c *cli.Context) {
	data, err := start(crypto.SHA384, sha512.New384(), input)
	if err != nil {
		fmt.Println(err)
		return
	}

	if filename != "" {
		writeFile(filename, data)
	} else {
		fmt.Println(string(data))
	}
}

func cmdSha512(c *cli.Context) {
	data, err := start(crypto.SHA512, sha512.New(), input)
	if err != nil {
		fmt.Println(err)
		return
	}

	if filename != "" {
		writeFile(filename, data)
	} else {
		fmt.Println(string(data))
	}
}

func cmdSha512_224(c *cli.Context) {
	data, err := start(crypto.SHA512_224, sha512.New512_224(), input)
	if err != nil {
		fmt.Println(err)
		return
	}

	if filename != "" {
		writeFile(filename, data)
	} else {
		fmt.Println(string(data))
	}
}

func cmdSha512_256(c *cli.Context) {
	data, err := start(crypto.SHA512_256, sha512.New512_256(), input)
	if err != nil {
		fmt.Println(err)
		return
	}

	if filename != "" {
		writeFile(filename, data)
	} else {
		fmt.Println(string(data))
	}
}

func start(h crypto.Hash, crypto hash.Hash, path string) ([]byte, error) {
	var factory iodriver.DriverFactory
	factory = &s3driver.S3Driver{
		Profile: "default",
		Debug:   false,
	}

	driver := factory.NewDriver()
	start := time.Now()
	buf, err := compute(driver, crypto, path)
	if err != nil {
		return nil, err
	}

	val := base64.StdEncoding.EncodeToString(buf)
	sec := (time.Now().Sub(start)).Seconds()
	hashinfo := HashInfo{
		DateTime: start,
		Path:     path,
		Hash:     strconv.Itoa(int(h)),
		Binary:   fmt.Sprintf("%x", buf),
		Base64:   val,
		Seconds:  fmt.Sprintf("%f", sec),
	}

	data, err := json.MarshalIndent(hashinfo, "", " ")
	//data, err := json.Marshal(hashinfo)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func compute(driver iodriver.Driver, crypto hash.Hash, path string) ([]byte, error) {
	file, err := driver.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	//buf := make([]byte, 1 * 1024 * 1024)
	buf := make([]byte, 4096)
	for {
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			return nil, err
		}
		if n == 0 {
			break
		}

		if _, err := crypto.Write(buf[:n]); err != nil {
			return nil, err
		}
	}

	return crypto.Sum(nil), nil
}

func writeFile(filename string, data []byte) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	n, err := file.Write(data)
	if err == nil && n < len(data) {
		return io.ErrShortWrite
	}
	return err
}
