package main

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v1"
)

type Config struct {
	Auth     []string
	Channels []string

	Nick, User, Server string
}

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

func (c *Config) init() {
	c.Nick = "noye"
	c.User = "museun"
	c.Server = "localhost:6667"
	c.Channels = []string{"#noye"}
	c.Auth = []string{"museun"}
}

func (c *Config) toMap() map[string]string {
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
