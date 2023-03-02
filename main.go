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

func responseReturn(w http.ResponseWriter) {

}

func rewriteHttp(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	if method == "OPTIONS" {
		w.WriteHeader(200)
		return
	}

	path := r.RequestURI
	header := r.Header

	if strings.TrimSpace(header.Get(accessKey)) != accessSecret {
		w.WriteHeader(500)
		r, _ := json.Marshal(map[string]string{"msg": "Prohibit unauthorized users from accessing the system"})
		w.Write(r)
		return
		//panic("Prohibit unauthorized users from accessing the system")
	} else {
		header.Del(accessKey)
	}

	domain := strings.TrimSpace(header.Get(reDomain))
	if domain == "" {
		w.WriteHeader(500)
		r, _ := json.Marshal(map[string]string{"msg": "domain is empty"})
		w.Write(r)
		return
		//panic("domain is empty")
	} else {
		header.Del(reDomain)
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
	if res.err != nil {
		w.WriteHeader(500)
		r, _ := json.Marshal(map[string]string{"msg": "request error" + res.err.Error()})
		w.Write(r)
		return
	}
	fmt.Println(string(res.Body), res.Header)
	hh := w.Header()
	mapToHeader(&res.Header, &hh)

	w.WriteHeader(200)
	w.Write(res.Body)
}

type SimpleResponse struct {
	Body   []byte
	Header map[string][]string
	err    error
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
		return &SimpleResponse{err: err}
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	responseHeader := resp.Header
	body, err := ioutil.ReadAll(resp.Body)

	return &SimpleResponse{
		Body:   body,
		Header: responseHeader,
		err:    err,
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
	http.HandleFunc("/", cors(rewriteHttp))
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		panic(err)
	}
}

func cors(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token,Authorization,Token,Hf-Access-Key,Re-Domain") //header的类型
		w.Header().Add("Access-Control-Allow-Credentials", "true")                                                                          //设置为true，允许ajax异步请求带cookie信息
		w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, HEAD, PATCH")
		w.Header().Set("content-type", "application/json;charset=UTF-8")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		f(w, r)
	}
}
