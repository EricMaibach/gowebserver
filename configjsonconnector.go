package main

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type configJSONConnector struct {
	data ConfigData
}

func (c *configJSONConnector) GetConfig() *ConfigData {
	return c.data
}

func (c *configJSONConnector) SetConfig(data ConfigData) error {
	ex, err := os.Executable()
	if err != nil {
		return err
	}
	dir, _ := filepath.Split(ex)
	path := filepath.Join(dir, "gowebserver.config.json")
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	e := json.NewEncoder(f)
	err = e.Encode(data)
	if err != nil {
		return err
	}
	c.data, _ = c.getConfig()
	return nil
}

func (c *configJSONConnector) getConfig() (*ConfigData, error) {
	ex, err := os.Executable()
	if err != nil {
		return nil, err
	}
	slogger.Info("Got file path: ", ex)
	dir, _ := filepath.Split(ex)
	slogger.Info("Got folder: ", dir)
	path := filepath.Join(dir, "gowebserver.config.json")
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	conf := &ConfigData{}

	r := json.NewDecoder(f)
	err = r.Decode(&conf)
	if err != nil {
		return nil, err
	}

	if conf.WWWPort == "" {
		conf.HTTPPort = "80"
	}

	if conf.WWWFolder == "" {
		conf.WWWFolder = "~/ui"
	}

	return conf, nil
}

func NewConfigJSONConnector() *configJSONConnector {
	cc := new(configJSONConnector)
	var err error
	cc.data, err = cc.getConfig()
	if err != nil {
		cc.data = new(ConfigData)
		cc.SetConfig(*cc.data)
	}
	return cc
}
