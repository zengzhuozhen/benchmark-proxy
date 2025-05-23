package core

import (
	"bytes"
	"io"
	"net/http"
)

type TagReplaceHandler interface {
	Replace(req *http.Request)
}

type URLQueryTagHandler struct {
	parser *TagCompoundParser
}

func (h *URLQueryTagHandler) Replace(req *http.Request) {
	query := req.URL.Query()
	for key, values := range query {
		for i, val := range values {
			values[i] = h.parser.ParseCustomizeTag(val)
		}
		query[key] = values
	}
	req.URL.RawQuery = query.Encode()
}

type BodyTagHandler struct {
	parser *TagCompoundParser
}

func (h *BodyTagHandler) Replace(req *http.Request) {
	bodyBytes, err := io.ReadAll(req.Body)
	if err != nil {
		bodyBytes = nil
	}
	defer func() {
		req.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	}()
	if len(bodyBytes) > 0 {
		parsed := h.parser.ParseCustomizeTag(string(bodyBytes))
		bodyBytes = []byte(parsed)
	}
}
