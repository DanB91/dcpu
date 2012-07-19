// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Config struct {
	Address     string `json:"-"` // Listen address for server.
	Timeout     uint   // Shut the server down after X seconds of idleness.
	IncludePath string // DCPU path, pointing to standard library.
	ProjectPath string // Base directory for code projects.
}

func NewConfig() *Config {
	c := new(Config)
	c.Timeout = 10
	c.Address = os.Getenv("DCPU_IDE_ADDRESS")
	c.IncludePath = os.Getenv("DCPU_PATH")

	if len(c.Address) == 0 {
		c.Address = ":7070"
	}

	c.ProjectPath = getHomeDir()
	return c
}

// Load loads configuration data from a file.
func (c *Config) Load(file string) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}

	json.Unmarshal(data, c)
}

// Load saves configuration data from to a file.
func (c *Config) Save(file string) {
	// Cherry-pick the fields we want to store.
	// 
	// We can't use struct field tags here because the config
	// struct is also used directly in the apiConfig call.
	// It has different requirements for which tags are in place.

	cfg := struct {
		Address     string
		ProjectPath string
		Timeout     uint
	}{
		config.Address,
		config.ProjectPath,
		config.Timeout,
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return
	}

	ioutil.WriteFile(file, data, 0600)
}
