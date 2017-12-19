package driver

import "io"

// Driver ...
type Driver interface {
	Open(string) (io.ReadCloser, error)
}
