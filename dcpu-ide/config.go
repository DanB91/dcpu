// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import "os"

type Config struct {
	Address  string // Listen address for server.
	DcpuPath string // Include path for dcpu library code.
	Timeout  uint   // Shut the server down after X seconds of idleness.
}

func NewConfig() *Config {
	c := new(Config)
	c.Timeout = 10
	c.Address = os.Getenv("DCPU_IDE_ADDRESS")
	c.DcpuPath = os.Getenv("DCPU_PATH")

	if len(c.Address) == 0 {
		c.Address = ":7070"
	}

	return c
}
