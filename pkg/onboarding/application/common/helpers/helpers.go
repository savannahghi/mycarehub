package helpers

import (
	"net/url"
)

// MergeURLValues merges > 1 url.Values into one
func MergeURLValues(values ...url.Values) url.Values {
	merged := url.Values{}
	for _, value := range values {
		for k, v := range value {
			merged[k] = v
		}
	}
	return merged
}
