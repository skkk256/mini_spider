package request

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"
)

func TestNewRequest(t *testing.T) {
	data := url.Values{}
	req, err := NewRequest("GET", "http://www.baidu.com", data, 0, 0, true)
	if err != nil {
		t.Error(err)
	}

	client := &http.Client{}
	resp, _ := client.Do(req.Request)
	fmt.Println(resp)
	fmt.Println()
}