package main

import (
    "flag"
    "fmt"
	"path"
	"runtime"
	"os"
)

import (
	"shengkuan/mini_spider/utils"
	"shengkuan/mini_spider/crawler"
	"shengkuan/mini_spider/config"
	"shengkuan/mini_spider/downloader"
	"shengkuan/mini_spider/parser"
)
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
	confFile = path.Clean(confFile)
	cfg, err := config.LoadConfigFromFile(confFile)
	if err != nil {
		fmt.Println("Load Config File Err: " + err.Error())
		os.Exit(1)
	}

	logDir = path.Clean(logDir)
	if !utils.IsDirExist(logDir) {
		if ok, err := utils.Mkdir(logDir); !ok {
			fmt.Println("Create Log Dir Err: ", err.Error())
			os.Exit(1)
		}
	}
	utils.InitialLogger("mSpider", "INFO", logDir, true, "midnight", 3)

	runtime.GOMAXPROCS(runtime.NumCPU())
	mainSpider := crawler.NewCrawler(cfg.Spider.ThreadCount, cfg.Spider.MaxDepth, cfg.Spider.CrawlInterval)

	var origins map[string]bool
	origins, err = utils.GetSeedFromFile(cfg.Spider.UrlListFile, false)
	if err != nil {
		utils.Logger.Error(err)
		os.Exit(1)
	}
	if len(origins) <= 0 {
		utils.Logger.Warn("urls empty: ", cfg.Spider.UrlListFile)
		os.Exit(1)
	}
	mainSpider.AppendOrigins(origins)

	p := parser.NewParser(cfg.Spider.TargetUrl)
	mainSpider.SetParser(p)

	d := downloader.NewDownloader(cfg.Spider.CrawlTimeout, cfg.Spider.OutputDirectory)
	mainSpider.AddDownloader(d)

	mainSpider.Start()
}

// func download(url string, queue chan string) {
// 	visited[url] = true
// 	client := &http.Client{}
// 	req, _ := http.NewRequest("GET", url, nil)
// 	// 自定义Header
// 	req.Header.Set("User-Agent", "Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1)")

// 	resp, err := client.Do(req)
// 	if err != nil {
// 		fmt.Println("http get error", err)
// 		return
// 	}
// 	defer resp.Body.Close()

// 	links := collectlinks.All(resp.Body)
// 	for _, link := range links {
// 		absolute := urlJoin(link, url)
// 		if url != " " {
// 			if !visited[absolute] {
// 				fmt.Println("parse url", absolute)
// 				go func() {
// 					queue <- absolute
// 				}()
// 			}
// 		}
// 	}
// }

// func urlJoin(href, base string) string {
// 	uri, err := url.Parse(href)
// 	if err != nil {
// 		return " "
// 	}
// 	baseUrl, err := url.Parse(base)
// 	if err != nil {
// 		return " "
// 	}
// 	return baseUrl.ResolveReference(uri).String()
// }