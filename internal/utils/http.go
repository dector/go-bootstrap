package utils

import (
	"mime"
	"net/http"
	"strconv"
	"strings"
)

// AcceptsJSON reports whether the request's Accept header explicitly accepts
// application/json. It inspects every Accept field value, handles
// comma-separated media ranges, parameters and q-values, and ignores entries
// with q=0.
func AcceptsJSON(r *http.Request) bool {
	for _, value := range r.Header.Values("Accept") {
		for _, part := range strings.Split(value, ",") {
			mediaType, params, err := mime.ParseMediaType(part)
			if err != nil || mediaType != "application/json" {
				continue
			}

			if q, ok := params["q"]; ok {
				if v, err := strconv.ParseFloat(q, 64); err == nil && v <= 0 {
					continue
				}
			}

			return true
		}
	}

	return false
}
