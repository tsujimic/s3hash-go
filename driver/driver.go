package iodriver

import "io"

// DriverFactory ...
type DriverFactory interface {
	NewDriver() Driver
}

// Driver ...
type Driver interface {
	Open(string) (io.ReadCloser, error)
}
