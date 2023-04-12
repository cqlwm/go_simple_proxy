package main

import (
	"http_forwarder_go/util"
	"net/http"
	"os"
	"strings"
)

const (
	cycleCheckKey   string = "Hf-Cycle-Check"
	reDomain        string = "Re-Domain"
	accessKey       string = "Hf-Access-Key"
	filterHeaderKey string = "Filter-Headers"
	keepHeaderKey   string = "Keep-Headers"
)

var accessSecret = os.Getenv("HF_ACCESS_SECRET")

func main() {
	httpHandlerWrapper("/", util.ShiftRequest)
	_ = http.ListenAndServe(":80", nil)
}

type ShiftFunc0 func(*http.Request, string) *util.SimpleHttpResponse

func setHeader(w http.ResponseWriter) {
	headers := map[string]string{
		"Access-Control-Allow-Origin":      "*",
		"Access-Control-Allow-Headers":     "Content-Type,AccessToken,X-CSRF-Token,Authorization,Token,Hf-Access-Key,Re-Domain",
		"Access-Control-Allow-Credentials": "true",
		"Access-Control-Allow-Methods":     "POST, GET, OPTIONS, PUT, DELETE, HEAD, PATCH",
		"content-type":                     "application/json;charset=UTF-8",
	}
	for k, v := range headers {
		w.Header().Set(k, v)
	}
}

func httpHandlerWrapper(pattern string, f ShiftFunc0) {

	handler := func(w http.ResponseWriter, r *http.Request) {
		setHeader(w)

		rOrigin := r.Header.Get("Origin")
		if rOrigin != "" {
			// Prevents the front end from submitting Origin.Header repeatedly
			origins := strings.Split(rOrigin, ",")
			origin := strings.TrimSpace(origins[len(origins)-1])
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusAccepted)
			return
		}

		// check header
		var errmsg string
		if util.MapDelGet(r.Header, accessKey) != accessSecret {
			errmsg = "Prohibit unauthorized users from accessing the system"
		}

		domain := util.MapDelGet(r.Header, reDomain)
		if domain == "" {
			errmsg = "Re-Domain is empty"
		}

		if util.MapDelGet(r.Header, cycleCheckKey) != "" {
			errmsg = "The request is stuck in a loop; Verify the [Re-Domain] field in the request header."
		} else {
			r.Header.Set(cycleCheckKey, "1")
		}

		if errmsg != "" {
			w.Header().Set("content-type", "text/plain;charset=UTF-8")
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(errmsg))
			return
		}

		// keep header
		keeps := util.MapValues(r.Header, keepHeaderKey)
		if keeps != nil {
			util.FilterMapKeys(r.Header, keeps, true)
		} else {
			filters := util.MapValues(r.Header, filterHeaderKey)
			if filters != nil {
				util.FilterMapKeys(r.Header, filters, false)
			}
		}

		simpleResponse := f(r, domain)
		delete(simpleResponse.Header, "Access-Control-Allow-Origin")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		mergeHandler := util.HttpHeadMergeHandler{
			Target: w.Header(),
			Heads:  []http.Header{simpleResponse.Header},
		}
		mergeHandler.Invoke()
		w.WriteHeader(simpleResponse.State)
		_, _ = w.Write(simpleResponse.Body)
		return
	}

	http.HandleFunc(pattern, handler)
}
