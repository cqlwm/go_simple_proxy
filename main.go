package main

import (
	"encoding/json"
	"http_forwarder_go/util"
	"io"
	"net/http"
	"os"
	"strings"
)

const (
	cycleCheckKey string = "Hf-Cycle-Check"
	reDomain      string = "Re-Domain"
	accessKey     string = "Hf-Access-Key"
)

var accessSecret = os.Getenv("HF_ACCESS_SECRET")

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

func failed(msg string) *util.SimpleHttpResponse {
	b, _ := json.Marshal(map[string]string{"msg": msg})
	return &util.SimpleHttpResponse{
		State: 500,
		Body:  b,
	}
}

func rewriteHttp(r *http.Request) *util.SimpleHttpResponse {
	method := r.Method
	path := r.RequestURI

	header := r.Header

	// check header
	if delGet(header, accessKey) != accessSecret {
		return failed("Prohibit unauthorized users from accessing the system")
	}

	domain := delGet(header, reDomain)
	if domain == "" {
		return failed("Re-Domain is empty")
	}

	if delGet(header, cycleCheckKey) != "" {
		return failed("The request is stuck in a loop; Verify the [Re-Domain] field in the request header.")
	} else {
		header.Set(cycleCheckKey, "1")
	}

	// request
	data, _ := io.ReadAll(r.Body)
	res := *util.DoRequest(method, domain+path, header, data)
	if res.Err != nil {
		return failed(res.Err.Error())
	}

	return &res
}

type HandlerFunc0 func(*http.Request) *util.SimpleHttpResponse

func main() {
	http.HandleFunc("/", cors(rewriteHttp))
	_ = http.ListenAndServe(":80", nil)
}

func cors(f HandlerFunc0) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rOrigin := r.Header.Get("Origin")
		if rOrigin != "" {
			// Prevents the front end from submitting Origin.Header repeatedly
			origins := strings.Split(r.Header.Get("Origin"), ",")
			origin := strings.TrimSpace(origins[len(origins)-1])
			w.Header().Set("Access-Control-Allow-Origin", origin)
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

		simpleResponse := f(r)
		delete(simpleResponse.Header, "Access-Control-Allow-Origin")

		mergeHandler := util.HttpHeadMergeHandler{
			Target: w.Header(),
			Heads:  []http.Header{simpleResponse.Header},
		}
		mergeHandler.Invoke()

		w.WriteHeader(simpleResponse.State)
		_, _ = w.Write(simpleResponse.Body)
		return
	}
}
