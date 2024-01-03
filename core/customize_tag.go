package core

import (
	"encoding/hex"
	"fmt"
	"github.com/google/uuid"
	"math/rand"
	"regexp"
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
)

var defaultParseReg = regexp.MustCompile("\\${(.+?)\\}")

func NewTagCompoundParser() *TagCompoundParser {
	return &TagCompoundParser{iteratorCount: 1}
}

type TagCompoundParser struct {
	iteratorCount int
}

func (p *TagCompoundParser) ParseCustomizeTag(content string) string {
	return p.parseCustomizeTag(defaultParseReg, content)
}

func (p *TagCompoundParser) parseCustomizeTag(reg *regexp.Regexp, content string) string {
	return reg.ReplaceAllStringFunc(content, func(s string) string {
		switch s {
		case fmt.Sprintf("${%s}", TagInt):
			return fmt.Sprintf("%d", rand.Int31()>>24)
		case fmt.Sprintf("${%s}", TagFloat):
			return fmt.Sprintf("%.2f", rand.Float32())
		case fmt.Sprintf("${%s}", TagString):
			rand.Seed(time.Now().UnixNano())
			result := make([]byte, 10/2)
			rand.Read(result)
			return hex.EncodeToString(result)
		case fmt.Sprintf("${%s}", TagUUID):
			return uuid.New().String()
		case fmt.Sprintf("${%s}", TagIncr):
			defer func() { p.iteratorCount++ }()
			return fmt.Sprintf("%d", p.iteratorCount)
		default:
			s = strings.ReplaceAll(strings.ReplaceAll(s, "${", ""), "}", "")
			if strings.HasPrefix(s, TagList) {
				s = strings.TrimRight(strings.ReplaceAll(s, TagList+"[", ""), "]")
				list := strings.Split(s, ",")
				rand.Seed(time.Now().UnixNano())
				return list[rand.Intn(len(list))]
			}
			return ""
		}
	})
}
