package config

import (
	"fmt"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

// Load returns Configuration struct
func Load(env string) (*Configuration, error) {
	bytes, err := ioutil.ReadFile("config." + env + ".yaml")
	if err != nil {
		return nil, fmt.Errorf("error reading config file, %s", err)
	}
	var cfg = new(Configuration)
	if err := yaml.Unmarshal(bytes, cfg); err != nil {
		return nil, fmt.Errorf("unable to decode into struct, %v", err)
	}
	return cfg, nil
}

// Configuration holds data necessery for configuring application
type Configuration struct {
	Server *Server `yaml:"server"`
}

// Server holds data necessery for server configuration
type Server struct {
	Port         string
	Debug        bool
	ReadTimeout  int
	WriteTimeout int
}
