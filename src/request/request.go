package request

import (
	"errors"
	"net/http"
	"strings"
	"net/url"
	"fmt"
)

type Request struct {
	*http.Request
	retryTimes int
	depth int
	isDownload bool
}

func NewRequest(method string, url string, data url.Values, depth int, retryTimes int, isDownload bool) (*Request, error) {
	if url == "" {
		return nil, errors.New("url is empty")
	}
	req, err := http.NewRequest(method, url, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, errors.New("new request error: " + err.Error())
	}

	return &Request{
		Request: req, 
		retryTimes: retryTimes,
		isDownload: isDownload,
		depth: depth}, nil
}

func (r *Request) SetIsDownload(isDownload bool) {
	if r == nil {
		return
	}
	r.isDownload = isDownload
}

func (r *Request) Valid() bool {
	return r.Url() != ""
}

func (r *Request) Url() string {
	if r == nil {
		return ""
	}
	return r.URL.String()
}

func (r *Request) Depth() int {
	if r == nil {
		return 0
	}
	return r.depth
}

func (r *Request) IsDownload() bool {
	if r == nil {
		return false
	}
	return r.isDownload
}

func (r *Request) RetryTimes() int {
	if r == nil {
		return 0
	}
	return r.retryTimes
}

func (r *Request) SetRetryTimes(times int) {
	if r == nil {
		return
	}
	r.retryTimes = times
}

func (r *Request) String() string {
	return fmt.Sprintf(
		"Request Method: %s, Url: %s, Retry Times: %d, Depth: %d, Download: %t",
		r.Request.Method,
		r.Request.URL.String(),
		r.retryTimes,
		r.depth,
		r.isDownload)
}