package downloader

import (
	"shengkuan/mini_spider/request"
	"shengkuan/mini_spider/utils"
	"strings"
	"testing"
)

func TestDownload(t *testing.T) {

	url :=  "http://ir.baidu.com"
	req, _ := request.NewRequest("GET", url, nil, utils.SEED_START_DEPTH, 0, false)

	var timeout uint
	timeout = 2
	outputDir := ""
	d := NewDownloader(int(timeout), outputDir)
	p, err := d.Download(req)
	if err != nil {
		t.Error("http download failed!")
	}

	strBody := p.StringRespBody()
	if !strings.Contains(strBody, "baidu") {
		t.Error("download html page failed!")
	}
}
