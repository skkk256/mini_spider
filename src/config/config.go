package config

import (
	"errors"
	"path/filepath"
)

import (
	"github.com/go-gcfg/gcfg"
)

import (
	"shengkuan/mini_spider/utils"
)

type Config struct {
	Spider struct {
		UrlListFile     string 
		OutputDirectory string
		MaxDepth        int
		CrawlInterval   int
		CrawlTimeout    int
		TargetUrl       string
		ThreadCount     int
	}
}

func (c *Config) Check() (bool, error) {
	if !utils.IsFileExist(c.Spider.UrlListFile) {
		return false, errors.New("UrlListFile not exist: " + c.Spider.UrlListFile)
	}
	if !utils.IsDirExist(c.Spider.OutputDirectory) {
		return false, errors.New("OutputDirectory not exist: " + c.Spider.OutputDirectory)
	}
	if c.Spider.CrawlInterval < 0 {
		return false, errors.New("CrawlInterval must greater than zero")
	}
	if c.Spider.CrawlTimeout <= 0 {
		return false, errors.New("CrawlTimeout must greater than zero")
	}
	if c.Spider.TargetUrl == "" {
		return false, errors.New("TargetUrl empty")
	}
	if c.Spider.ThreadCount <= 0 {
		return false, errors.New("ThreadCount must greater than zero")
	}
	return true, nil
}

func LoadConfigFromFile(filePath string) (*Config, error) {
	var conf Config
	err := gcfg.ReadFileInto(&conf, filePath)
	if err != nil {
		return nil, err
	}

	configDir := filepath.Dir(filePath)
	conf.Spider.UrlListFile = resolvePath(configDir, conf.Spider.UrlListFile)
	conf.Spider.OutputDirectory = resolvePath(configDir, conf.Spider.OutputDirectory)

	if _, err := conf.Check(); err != nil {
		return nil, err
	}

	return &conf, nil
}

func resolvePath(basePath, targetPath string) string {
	if filepath.IsAbs(targetPath) {
		return targetPath
	}
	return filepath.Join(basePath, targetPath)
}
