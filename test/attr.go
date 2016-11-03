package test

import (
	"golang.org/x/net/html"
)

func extractAttrs(attrs []html.Attribute, keys []string) map[string]string {
	attrMap := make(map[string]string)
	for _, attr := range attrs {
		for i, key := range keys {
			if attr.Key == key {
				attrMap[key] = attr.Val
				// Delete the key from the keys slice as found
				keys = append(keys[:i], keys[i+1:]...)
			}
		}
	}
	return attrMap
}

func attrPresent(attrs []html.Attribute, key string) bool {
	for _, attr := range attrs {
		if attr.Key == key {
			return true
		}
	}
	return false
}
