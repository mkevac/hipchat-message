package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	Token   string
	BaseURL string
	baseURL *url.URL
}

func (c *Config) SetBaseURL(u string) error {
	ur, err := url.Parse(u)
	if err != nil {
		return err
	}
	c.baseURL = ur
	c.BaseURL = u
	return nil
}

func (c *Config) SetToken(t string) error {
	c.Token = t
	return nil
}

func (c *Config) GetBaseURL() *url.URL {
	return c.baseURL
}

func (c *Config) Save() error {
	conf := xdgApp.ConfigPath(configFileName)

	buf, err := json.Marshal(c)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(conf), 0700); err != nil {
		return err
	}

	if err = ioutil.WriteFile(conf, buf, 0600); err != nil {
		return err
	}

	return nil
}

func (c *Config) Load() error {
	conf := xdgApp.ConfigPath(configFileName)

	buf, err := ioutil.ReadFile(conf)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(buf, c); err != nil {
		return err
	}

	if err := c.SetBaseURL(c.BaseURL); err != nil {
		return err
	}

	return nil
}

func (c *Config) Check() error {
	if !strings.HasPrefix(c.BaseURL, "http://") && !strings.HasPrefix(c.BaseURL, "https://") {
		return fmt.Errorf("BaseURL is not an URL")
	}
	return nil
}
