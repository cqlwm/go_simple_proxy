package main

import (
	"encoding/json"
	"http_forwarder_go/util"
	"io/ioutil"
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

func rewriteHttp(w http.ResponseWriter, r *http.Request) *util.SimpleHttpResponse {
	method := r.Method
	path := r.RequestURI

	header := r.Header

	// check header
	if delGet(header, accessKey) != accessSecret {
		returnResponse(w, 500, "Prohibit unauthorized users from accessing the system")
		return nil
	}

	domain := delGet(header, reDomain)
	if domain == "" {
		returnResponse(w, 500, "Re-Domain is empty")
		return nil
	}

	if delGet(header, cycleCheckKey) != "" {
		returnResponse(w, 500, "The request is stuck in a loop; Verify the [Re-Domain] field in the request header.")
		return nil
	} else {
		header.Set(cycleCheckKey, "1")
	}

	// request
	data, _ := ioutil.ReadAll(r.Body)
	res := *util.DoRequest(method, domain+path, header, data)
	if res.Err != nil {
		returnResponse(w, 500, res.Err.Error())
		return nil
	}

	return &res
}

type HandlerFunc0 func(http.ResponseWriter, *http.Request) *util.SimpleHttpResponse

func main() {
	http.HandleFunc("/", cors(rewriteHttp))
	_ = http.ListenAndServe(":80", nil)
}

func cors(f HandlerFunc0) http.HandlerFunc {
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

		res := f(w, r)
		delete(res.Header, "Access-Control-Allow-Origin")

		// write response
		mergeHandler := util.HttpHeadMergeHandler{
			Target: w.Header(),
			Heads:  []http.Header{res.Header},
		}
		mergeHandler.Invoke()

		w.WriteHeader(res.State)
		w.Write(res.Body)
		return
	}
}
