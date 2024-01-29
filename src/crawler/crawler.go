package crawler

import (
	// "fmt"
	// "os"
	// "strings"
	"time"
	// "github.com/baidu/go-lib/queue"
)

import (
	"shengkuan/mini_spider/downloader"
	"shengkuan/mini_spider/manager"
	"shengkuan/mini_spider/parser"
	"shengkuan/mini_spider/request"
	"shengkuan/mini_spider/utils"
	"shengkuan/mini_spider/queue"
)

type Crawler struct {
	concurrency   int
	crawlMaxDepth int
	origins       []string
	crawlInterval time.Duration
	spiderManager *manager.RoutineManager
	taskQueue     *queue.TaskQueue
	parser		  *parser.Parser
	downloader    *downloader.Downloader
}


func NewCrawler(concurrency int, maxDepth int, interval int) *Crawler {
	crawInterval := time.Duration(interval) * time.Second
	spiderManager := manager.NewRoutineManager(concurrency)
	taskQueue := queue.NewTaskQueue()
	return &Crawler{
		concurrency:   concurrency,
		crawlMaxDepth: maxDepth,
		crawlInterval: crawInterval,
		spiderManager: spiderManager,
		taskQueue:     taskQueue,
	}
}

func (c *Crawler) SetParser(parser *parser.Parser) {
	c.parser = parser
}

// func (c *Crawler) addStartReq(req *request.Request) bool {
// 	if req == nil || !req.Valid() {
// 		return false
// 	}
// 	c.taskQueue.Push(req)
// 	c.origins = append(c.origins, req.Url())
// 	return true
// }

// func (c *Crawler) addReq(req *request.Request) bool {
// 	if req == nil || !req.Valid() {
// 		return false
// 	}
// 	c.taskQueue.Push(req)
// 	return true
// }

func (c *Crawler) AddDownloader(downloader *downloader.Downloader) {
	c.downloader = downloader
}


func (c *Crawler) AppendOrigins(urls map[string]bool) {
	for url, isDownload := range urls {
		req, error := request.NewRequest("GET", url, nil, utils.SEED_START_DEPTH, 0, isDownload)
		if error != nil {
			utils.Logger.Error("NewRequest error: ", error)
			continue
		}
		if req == nil || !req.Valid() {
			continue
		}
		c.taskQueue.Push(req)
		c.origins = append(c.origins, req.Url())

	}
}

// func (c *Crawler) getSubUrls(res request.Response) []string {
// 	resp, err := utils.GetResp(url)
// 	if err != nil {
// 		return nil
// 	}
// 	defer resp.Body.Close()
// 	links := collectlinks.All(resp.Body)
// 	return links
// }

func (c *Crawler) Start() {
	for {
		if c.isFinish() {
			break
		}

		if c.taskQueue.Empty() {
			continue
		}

		sub := c.spiderManager.GetOne()
		if !sub {
			continue
		}

		go func() {
			defer c.spiderManager.FreeOne()
			for {
				req, err := c.taskQueue.Pop()
				if err != nil {
					utils.Logger.Error("get req from queue failed: ", err)
					continue
				}
				if req == nil {
					break
				}
				utils.Logger.Info(">>>>>> popping request <<<<<< ", req.Url())
				c.process(req)
				c.sleepInterval()
			}
		}()

	}
}

func (c *Crawler) process(req *request.Request) {
	page, err := c.downloader.Download(req)
	if err != nil || page == nil || !page.Success() {
		utils.Logger.Error("download req err: ", req, err)
		return
	}

	if page.Depth() >= c.crawlMaxDepth {
		page.SetIsParse(false)
	}
	err = c.parser.Parse(page)
	if err != nil {
		utils.Logger.Error("parse resp err: ", page, err)
		return
	}

	for _, req := range page.GetComingCrawlReq() {
		if req.Depth() > c.crawlMaxDepth {
			continue
		}
		c.taskQueue.Push(req)
	}
}

func (c *Crawler) sleepInterval() {
	time.Sleep(c.crawlInterval)
}

func (c *Crawler) isFinish() bool {
	if c.taskQueue.Count() == 0 && c.spiderManager.Used() == 0 {
		utils.Logger.Info("request queue nil, all sub spider are idle, Spider exit.")
		time.Sleep(30 * time.Second)
		return true
	}
	return false
}