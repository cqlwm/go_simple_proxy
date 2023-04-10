package util

import "strings"

func MapValues(header map[string][]string, key string) []string {
	return header[key]
}

func FilterMapKeys(headers map[string][]string, sub []string, keep bool) {
	for key := range headers {
		if keep && Contains(sub, key) == false {
			delete(headers, key)
		}
		if keep == false && Contains(sub, key) {
			delete(headers, key)
		}
	}
}

func Contains(arr []string, val string) bool {
	for _, item := range arr {
		if item == val {
			return true
		}
	}
	return false
}

func MapDelGet(m map[string][]string, key string) string {
	v := m[key]
	if v == nil {
		return ""
	}
	delete(m, key)
	return strings.TrimSpace(v[0])
}
