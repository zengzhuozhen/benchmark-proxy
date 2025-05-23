package core

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type BodyCheckType int

const (
	BodyCheckEqual BodyCheckType = iota
	BodyCheckContains
	BodyCheckRegex
)

type StatusCheckType int

const (
	StatusCheckEqual StatusCheckType = iota
	StatusCheckMulti
	StatusCheckRange
)

type ResponseChecker struct {
	// status check
	status          int             // 用于 StatusCheckEqual，要求响应状态码等于该值
	statusList      []int           // 用于 StatusCheckMulti，要求响应状态码等于列表中任一值
	statusRange     [2]int          // 用于 StatusCheckRange，要求响应状态码在该区间内（闭区间）
	statusCheckType StatusCheckType // 状态码校验类型（等值、多值、区间）
	// body check
	body             string         // 用于等值、包含、正则三种 body 校验，具体取决于 bodyCheckType
	bodyCheckType    BodyCheckType  // body 校验类型（等值、包含、正则）
	customRegPattern *regexp.Regexp // 用于 BodyCheckRegex，body 校验为正则时的编译表达式
}

func (c *ResponseChecker) Check(responseStatus int, responseMessage string) bool {
	return c.checkStatus(responseStatus) && c.checkBodyFlexible(responseMessage)
}

func (c *ResponseChecker) checkStatus(responseStatus int) bool {
	switch c.statusCheckType {
	case StatusCheckEqual:
		return responseStatus == c.status
	case StatusCheckMulti:
		for _, s := range c.statusList {
			if responseStatus == s {
				return true
			}
		}
		return false
	case StatusCheckRange:
		return responseStatus >= c.statusRange[0] && responseStatus <= c.statusRange[1]
	default:
		return false
	}
}

func (c *ResponseChecker) checkBodyFlexible(responseMessage string) bool {
	switch c.bodyCheckType {
	case BodyCheckEqual:
		return c.checkBodyEqual(responseMessage)
	case BodyCheckContains:
		return c.checkBodyContains(responseMessage)
	case BodyCheckRegex:
		return c.checkBodyRegex(responseMessage)
	default:
		return c.checkBodyEqual(responseMessage)
	}
}

func (c *ResponseChecker) checkBodyEqual(responseMessage string) bool {
	return c.body == "" || responseMessage == c.body
}

func (c *ResponseChecker) checkBodyContains(responseMessage string) bool {
	return c.body == "" || (strings.Contains(responseMessage, c.body))
}

func (c *ResponseChecker) checkBodyRegex(responseMessage string) bool {
	return c.customRegPattern != nil && c.customRegPattern.MatchString(responseMessage)
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
			checker.customRegPattern = regexp.MustCompile(match[1])
		}
	}
}

func NewResponseCheckerWithBodyType(status int, body string, bodyType BodyCheckType) *ResponseChecker {
	checker := &ResponseChecker{
		status:        status,
		body:          body,
		bodyCheckType: bodyType,
	}
	if bodyType == BodyCheckRegex {
		checker.customRegPattern = regexp.MustCompile(body)
	}
	return checker
}

// 自动推断工厂，支持 status/body 多种表达式
func SmartResponseChecker(statusExpr string, body string) *ResponseChecker {
	checker := &ResponseChecker{}
	// --- status check ---
	if strings.Contains(statusExpr, "-") {
		parts := strings.Split(statusExpr, "-")
		if len(parts) == 2 {
			min, _ := strconv.Atoi(strings.TrimSpace(parts[0]))
			max, _ := strconv.Atoi(strings.TrimSpace(parts[1]))
			checker.statusRange = [2]int{min, max}
			checker.statusCheckType = StatusCheckRange
		}
	} else if strings.Contains(statusExpr, ",") {
		vals := strings.Split(statusExpr, ",")
		for _, v := range vals {
			n, _ := strconv.Atoi(strings.TrimSpace(v))
			checker.statusList = append(checker.statusList, n)
		}
		checker.statusCheckType = StatusCheckMulti
	} else if statusExpr != "" {
		code, _ := strconv.Atoi(statusExpr)
		checker.status = code
		checker.statusCheckType = StatusCheckEqual
	} else {
		checker.status = http.StatusOK
		checker.statusCheckType = StatusCheckEqual
	}
	// --- body check ---
	if strings.HasPrefix(body, "@Reg[") && strings.HasSuffix(body, "]") {
		pattern := body[5 : len(body)-1]
		checker.body = pattern
		checker.bodyCheckType = BodyCheckRegex
		checker.customRegPattern = regexp.MustCompile(pattern)
	} else if strings.HasPrefix(body, "@Contains[") && strings.HasSuffix(body, "]") {
		substr := body[10 : len(body)-1]
		checker.body = substr
		checker.bodyCheckType = BodyCheckContains
	} else {
		checker.body = body
		checker.bodyCheckType = BodyCheckEqual
	}
	return checker
}
