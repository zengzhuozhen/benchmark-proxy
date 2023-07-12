package core

import "net/http"

type ResponseChecker struct {
	status int
	body   string
}

type ResponseCheckOption func(checker *ResponseChecker)

func NewResponseChecker(options ...ResponseCheckOption) *ResponseChecker {
	r := &ResponseChecker{}
	ResponseCheckerStatusRule(http.StatusOK)(r) // default:check response status 200
	for _, i := range options {
		i(r)
	}
	return r
}

func ResponseCheckerStatusRule(status int) ResponseCheckOption {
	return func(checker *ResponseChecker) {
		checker.status = status
	}
}

func ResponseCheckerBodyRule(body string) ResponseCheckOption {
	return func(checker *ResponseChecker) {
		checker.body = body
	}
}
