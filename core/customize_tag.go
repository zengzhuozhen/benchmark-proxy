package core

import (
	"encoding/hex"
	"fmt"
	"github.com/google/uuid"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	TagUUID   = "uuid"
	TagInt    = "int"
	TagString = "string"
	TagFloat  = "float"
	TagIncr   = "incr"
	TagList   = "list:"
	TagRange  = "range:"
)

var deepTagReg = map[string]*regexp.Regexp{
	TagList:  regexp.MustCompile(`list:\[([^\]]*)\]`),
	TagRange: regexp.MustCompile(`range:\[([\d,]+)\]`),
}

var defaultParseReg = regexp.MustCompile("\\${(.+?)\\}")

func NewTagCompoundParser() *TagCompoundParser {
	return &TagCompoundParser{iteratorCount: 1}
}

type TagCompoundParser struct {
	// 全局共用字段,这里不能区分不同字段,考虑不同字段有不同的tag后缀
	iteratorCount int
	rangeCount    int
}

func (p *TagCompoundParser) ParseCustomizeTag(content string) string {
	return p.parseCustomizeTag(defaultParseReg, content)
}

func (p *TagCompoundParser) parseCustomizeTag(reg *regexp.Regexp, content string) string {
	match := reg.FindStringSubmatch(content)
	if len(match) == 0 || len(match) > 2 {
		return content
	}
	s := match[1]
	switch s {
	case TagInt:
		return fmt.Sprintf("%d", rand.Int31()>>24)
	case TagFloat:
		return fmt.Sprintf("%.2f", rand.Float32())
	case TagString:
		rand.Seed(time.Now().UnixNano())
		result := make([]byte, 10/2)
		rand.Read(result)
		return hex.EncodeToString(result)
	case TagUUID:
		return uuid.New().String()
	case TagIncr:
		defer func() { p.iteratorCount++ }()
		return fmt.Sprintf("%d", p.iteratorCount)
	default:
		if strings.Contains(s, TagList) {
			noTag := p.parseCustomizeTag(deepTagReg[TagList], s)
			list := strings.Split(noTag, ",")
			rand.Seed(time.Now().UnixNano())
			return list[rand.Intn(len(list))]
		}
		if strings.Contains(s, TagRange) {
			noTag := p.parseCustomizeTag(deepTagReg[TagRange], s)
			list := strings.Split(noTag, ",")
			begin, _ := strconv.Atoi(list[0])
			end, _ := strconv.Atoi(list[1])
			if p.rangeCount < begin || p.rangeCount > end {
				p.rangeCount = begin
			}
			defer func() { p.rangeCount++ }()
			return fmt.Sprintf("%d", p.rangeCount)
		}
		return s
	}
}
