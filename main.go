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

func returnResponse(w http.ResponseWriter, status int, msg string) {
	w.WriteHeader(status)
	r, _ := json.Marshal(map[string]string{"msg": msg})
	w.Write(r)
}

func delGet(h http.Header, key string) string {
	v := strings.TrimSpace(h.Get(key))
	h.Del(key)
	return v
}

func rewriteHttp(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	path := r.RequestURI

	header := r.Header

	// check header
	if delGet(header, accessKey) != accessSecret {
		returnResponse(w, 500, "Prohibit unauthorized users from accessing the system")
		return
	}

	domain := delGet(header, reDomain)
	if domain == "" {
		returnResponse(w, 500, "Re-Domain is empty")
		return
	}

	if delGet(header, cycleCheckKey) != "" {
		returnResponse(w, 500, "The request is stuck in a loop; Verify the [Re-Domain] field in the request header.")
		return
	} else {
		header.Set(cycleCheckKey, "1")
	}

	// request
	data, _ := ioutil.ReadAll(r.Body)
	res := *doRequest(method, domain+path, header, data)
	if res.err != nil {
		returnResponse(w, 500, res.err.Error())
		return
	}

	// write response
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
		return &SimpleResponse{err: err}
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
		rOrigin := r.Header.Get("Origin")
		if rOrigin != "" {
			w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		} else {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}
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
