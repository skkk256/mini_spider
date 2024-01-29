package parser

import (
	"io"
	"net/http"
	"net/url"
	"testing"
)

import (
	"shengkuan/mini_spider/request"
	"shengkuan/mini_spider/response"
)

func TestParser(t *testing.T) {
	data := make(url.Values)
	url := "http://www.baidu.com"
	req, _ := request.NewRequest("GET", url, data, 0, 0, true)
	client := &http.Client{}

	page := response.NewResp(req)


	resp, err := client.Do(req.Request)
	if err != nil {
		t.Error(err)
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}

	//apply all the attrs to page
	page.SetRespBody(content)
	page.SetRespStatus(resp.StatusCode)
	page.SetRespHeader(resp.Header)
	page.SetIsParse(true)

	parser := NewParser(".*.(htm|html)$")
	err = parser.Parse(page)
	if err != nil {
		t.Error(err)
	}
	var urlLists []string
	for _, r := range page.GetComingCrawlReq() {
		if r.Url() == "http://map.baidu.com" || r.Url() == "http://news.baidu.com" {
			urlLists = append(urlLists, r.Url())
		}
	}
	if len(urlLists) != 2 {
		t.Error("test Parse() error")
	}

}
