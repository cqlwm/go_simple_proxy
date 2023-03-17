package util

import (
	"net/http"
	"testing"
)

func TestHttpHeadMergeHandler(t *testing.T) {

	m := map[string][]string{}
	m["a"] = []string{"1", "2"}

	n := map[string][]string{}
	n["b"] = []string{"1", "2"}

	h := HttpHeadMergeHandler{
		m,
		[]http.Header{n},
	}

	h.Invoke()

	for s, _ := range m {
		println(s)
	}

}