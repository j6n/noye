package config

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/j6n/noye/logger"
	"gopkg.in/yaml.v1"
)

var log = logger.Get()

// Config represents an irc configuration
type Config struct {
	Auth     []string
	Channels []string

	Nick, User, Server string
}

// NewConfig tries to load config.yaml, or a default
func NewConfig() *Config {
	conf := &Config{}

	_, err := os.Stat("config.yaml")
	if !os.IsNotExist(err) {
		data, err := ioutil.ReadFile("config.yaml")
		if err != nil {
			goto newConfig
		}

		if err := yaml.Unmarshal(data, conf); err != nil {
			log.Errorf("unable to laod yaml: %s\n%s\n", err, data)
			goto newConfig
		}

		return conf
	}

newConfig:
	conf.init()
	data, err := yaml.Marshal(conf)
	if err != nil {
		log.Fatalln(err)
	}

	if err := ioutil.WriteFile("config.yaml", data, 0666); err != nil {
		log.Fatalln(err)
	}

	log.Println("wrote default config.yaml, you should probably change this")
	return conf
}

// ToMap shittily converts the config to a string map
func (c *Config) ToMap() map[string]string {
	data, _ := json.Marshal(c)
	m := make(map[string]interface{})
	out := make(map[string]string)

	json.Unmarshal(data, &m)

	for k, v := range m {
		b, _ := json.Marshal(v)
		out[k] = string(b)
	}

	return out
}

func (c *Config) init() {
	c.Nick = "noye"
	c.User = "museun"
	c.Server = "localhost:6667"
	c.Channels = []string{"#noye"}
	c.Auth = []string{"museun"}
}
