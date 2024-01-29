package downloader

import (
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

import (
	"shengkuan/mini_spider/request"
	"shengkuan/mini_spider/response"
	"shengkuan/mini_spider/utils"
)

//https://developer.mozilla.org/zh-CN/docs/Web/HTTP/Basics_of_HTTP/MIME_types
//if page's content-type contained in the map,then skipped(don't parse it)
//using map[string]interface{} to get better performance(because the value is useless)
var skippedContentType = map[string]interface{}{
	"video": nil,
	"audio": nil,
	"image": nil,
	"messa": nil,
	"appli": nil,
}

type Downloader struct {
	client       *http.Client
	crawlTimeout time.Duration
	resDir       string
}

func NewDownloader(timeout int, resDir string) *Downloader {
	crawlTimeout := time.Duration(timeout) * time.Second
	client := &http.Client{
		Timeout: crawlTimeout,
	}
	return &Downloader{
		client:       client,
		crawlTimeout: crawlTimeout,
		resDir:       resDir,
	}
}

func (d *Downloader) SetClient(client *http.Client) {
	d.client = client
}

func (d *Downloader) SetTimeout(timeout int) *Downloader {
	d.crawlTimeout = time.Second * time.Duration(timeout)
	d.client.Timeout = d.crawlTimeout
	return d
}

func (d *Downloader) SetResDir(dir string) *Downloader {
	if dir != "" {
		d.resDir = dir
	}
	return d
}

//Download function means that we create a real client to make a real request
//then we will get a real resp
//at last, we apply all the attributes of the bellow resp to our "Page"(response.NewResp)
//in one word, "Page" is not directly downloaded by a client, it is just a dummy copy.

//this function will save the returned pages to file(on disk) and return an *response.Response object.

func (d *Downloader) Download(req *request.Request) (*response.Response, error) {
	utils.Logger.Info(">>>>>> Downloading <<<<<< ", req)
	req.Header.Set("User-Agent", "Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1)")
	resp, err := d.client.Do(req.Request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	//apply resp's attributions to page
	page := response.NewResp(req)
	page.SetRespHeader(resp.Header)
	page.SetRespBody(content)
	page.SetRespStatus(resp.StatusCode)
	page.SetRespCookies(resp.Cookies())

	//don't parse page when it's content-type is contained in the map
	contentType := http.DetectContentType(content)
	realType := strings.ToLower(contentType[:5])
	if _, exist := skippedContentType[realType]; exist {
		page.SetIsParse(false)
	} else {
		page.SetIsParse(true)
	}

	if req.IsDownload() && page.Success() {
		fileSave := url.QueryEscape(req.Url())
		res, err := os.Create(d.resDir + "/" + fileSave)
		if err != nil {
			utils.Logger.Warn("downloader create file failed ", fileSave, err)
			return page, nil
		}
		defer res.Close()
		_, err = res.Write(content)
		if err != nil {
			utils.Logger.Warn("[warn] downloader save resp to file failed ", err)
		}
	}
	return page, nil

}
