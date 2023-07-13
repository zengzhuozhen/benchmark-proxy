package core

import (
	"net/http"
	"regexp"
)

type ResponseChecker struct {
	status           int
	body             string
	isCustomReg      bool
	customRegPattern *regexp.Regexp
}

func (c *ResponseChecker) CheckStatus(responseStatus int) bool {
	return responseStatus == c.status
}

func (c *ResponseChecker) CheckBody(responseMessage string) bool {
	if c.isCustomReg {
		return c.customRegPattern.MatchString(responseMessage)
	}
	return responseMessage == c.body

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
		pattern := `@Reg\[(.+?)\]`
		re := regexp.MustCompile(pattern)
		match := re.FindStringSubmatch(body)
		if len(match) > 1 {
			checker.isCustomReg = true
			checker.customRegPattern = regexp.MustCompile(match[1])
		}
	}
}
