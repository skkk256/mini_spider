package parser

import (
	"regexp"
	"strings"
)

import (
	"github.com/jackdanger/collectlinks"
)

import (
	"shengkuan/mini_spider/utils"
	"shengkuan/mini_spider/response"
	"shengkuan/mini_spider/request"

)

type Parser struct {
	target *regexp.Regexp
}

func NewParser(pattern string) *Parser {
	return &Parser{target: regexp.MustCompile(pattern)}
}

func (p *Parser) Parse(resp *response.Response) (error) {	
	if resp == nil || !resp.IsParse() {
		return nil
	}

	bodyString := resp.StringRespBody()
	if bodyString == "" {
		return nil
	}
	utils.Logger.Info(">>>>>> Parsing <<<<<< ", resp)

	bodyReader := strings.NewReader(bodyString)

	strUTFBody, err := utils.TransCharsetUTF8(resp.RespHeader("Content-Type"), bodyReader)
	if err == nil {
		bodyReader = strings.NewReader(strUTFBody)
	}

	subUrl := collectlinks.All(bodyReader)
	for _, url := range subUrl {
		absUrl, err := utils.TransUrlFromRelToAbs(resp.URL(), url)
		if err != nil {
			utils.Logger.Error("TransUrlFromRelToAbs error: ", err)
			continue
		}
		// if _, ok := visited[absUrl]; ok {
		// 	continue
		// }
		req, err := request.NewRequest("GET", absUrl, nil, resp.Depth()+1, 0, false)
		if err != nil {
			utils.Logger.Error("NewRequest error: ", err)
			continue
		}
		if p.target.MatchString(absUrl) {
			req.SetIsDownload(true)
		}
		resp.AddComingCrawlReq(req)
	}
	return nil
}

