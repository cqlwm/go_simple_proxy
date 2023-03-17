package util

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
)

type SimpleHttpResponse struct {
	State  int
	Body   []byte
	Header map[string][]string
	Err    error
}

type HttpHeadMergeHandler struct {
	Target http.Header
	Heads  []http.Header
}

func (h *HttpHeadMergeHandler) Invoke() {
	if h.Target == nil {
		h.Target = http.Header{}
	}
	var r map[string][]string = h.Target

	for _, head := range h.Heads {
		var item map[string][]string = head
		for k, v := range item {
			r[k] = v
		}
	}
}

func DoRequest(method string, url string, header http.Header, data []byte) *SimpleHttpResponse {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		return &SimpleHttpResponse{Err: err}
	}

	merge := HttpHeadMergeHandler{
		Target: req.Header,
		Heads:  []http.Header{header},
	}
	merge.Invoke()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return &SimpleHttpResponse{Err: err}
	}

	responseHeader := resp.Header
	body, err := ioutil.ReadAll(resp.Body)

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	return &SimpleHttpResponse{
		State:  resp.StatusCode,
		Body:   body,
		Header: responseHeader,
		Err:    err,
	}
}
