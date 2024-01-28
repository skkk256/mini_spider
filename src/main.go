package main

import (
    "flag"
    "fmt"
    "io"
    "net/http"
)

import "shengkuan/mini_spider/utils"

var (
	displayVersion bool
	confFile       string
	logDir         string
)

func init() {
    flag.Usage = utils.DisplayHelpMenu
	flag.BoolVar(&displayVersion, "v", false, "display spider version then exit")
	flag.StringVar(&confFile, "c", "../conf/spider.conf", "set config file path")
	flag.StringVar(&logDir, "l", "../log/", "set log directory")
}

func main() {
    flag.Parse()
    if displayVersion {
        fmt.Println("mSpider version: 1.0.0")
        return
    }
    fmt.Println("confFile:", confFile)
    fmt.Println("logDir:", logDir)

    resp, err := http.Get("http://www.baidu.com/")
    if err != nil {
        fmt.Println("http get error", err)
        return
    }
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        fmt.Println("read error", err)
        return
    }
    fmt.Println(string(body))
}