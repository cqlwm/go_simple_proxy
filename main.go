package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

const cycleCheckKey = "Hf-Cycle-Check"
const reDomain = "Re-Domain"

const accessKey = "Hf-Access-Key"
const accessSecret = "UgFUrwVGktW9XbkozneV"

func rewriteHttp(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	path := r.RequestURI
	header := r.Header

	if strings.TrimSpace(header.Get(accessKey)) != accessSecret {
		panic("Prohibit unauthorized users from accessing the system")
	} else {
		header.Del(accessKey)
	}

	domain := strings.TrimSpace(header.Get(reDomain))
	if domain == "" {
		panic("domain is empty")
	}

	fallCycle := strings.TrimSpace(header.Get(cycleCheckKey))
	if fallCycle != "" {
		hh := w.Header()
		hh.Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(500)
		errMsg := "The request is stuck in a loop; Verify the [Re-Domain] field in the request header."
		w.Write([]byte(errMsg))
		return
	}

	data, _ := ioutil.ReadAll(r.Body)

	header.Set(cycleCheckKey, "1")
	res := *doRequest(method, domain+path, header, data)
	fmt.Println(string(res.Body), res.Header)

	hh := w.Header()
	mapToHeader(&res.Header, &hh)
	hh.Set("Content-Type", "text/plain; charset=utf-8")

	w.WriteHeader(200)
	w.Write(res.Body)

}

type SimpleResponse struct {
	Body   []byte
	Header map[string][]string
}

func doRequest(method string, url string, header http.Header, data []byte) *SimpleResponse {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		panic(err)
	}

	mapToHeader(headerToMap(header), &req.Header)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)

	responseHeader := resp.Header
	body, _ := ioutil.ReadAll(resp.Body)

	return &SimpleResponse{
		Body:   body,
		Header: responseHeader,
	}
}

func headerToMap(header http.Header) *map[string][]string {
	bs, _ := json.Marshal(header)
	var data map[string][]string
	err := json.Unmarshal(bs, &data)
	if err != nil {
		fmt.Println("Error:", err)
		panic("error Unmarshal failed")
	}
	return &data
}

func mapToHeader(hmap *map[string][]string, header *http.Header) {
	for k, vs := range *hmap {
		for _, v := range vs {
			header.Add(k, v)
		}
	}
}

func main() {
	http.HandleFunc("/", rewriteHttp)
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		panic(err)
	}
}
