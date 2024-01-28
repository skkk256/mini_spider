package request

import (
	"errors"
	"net/http"
	"strings"
	"net/url"
)

type Request struct {
	Req *http.Request
	RetryTime int
	Depth int
	IsDownload bool
}

func NewRequest(method string, url string, data url.Values, depth int, RetryTime int, IsDownload bool) (*Request, error) {
	if url == "" {
		return nil, errors.New("url is empty")
	}
	req, err := http.NewRequest(method, url, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, errors.New("new request error: " + err.Error())
	}

	return &Request{Req: req, 
		RetryTime: RetryTime,
		IsDownload: IsDownload,
		Depth: depth}, nil
}