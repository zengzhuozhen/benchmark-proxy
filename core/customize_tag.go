package core

import (
	"encoding/hex"
	"fmt"
	"github.com/google/uuid"
	"math/rand"
	"regexp"
	"time"
)

const (
	TagUUID   = "uuid"
	TagInt8   = "int8"
	TagInt16  = "int16"
	TagInt32  = "int32"
	TagInt    = "int"
	TagString = "string"
	TagFloat  = "float"
	TagIncr   = "incr"
)

func NewTagCompoundParser() *TagCompoundParser {
	return &TagCompoundParser{iteratorCount: 1}
}

type TagCompoundParser struct {
	iteratorCount int
}

func (p *TagCompoundParser) ParseCustomizeTag(content string) string {
	reg := regexp.MustCompile("\\${(.+?)\\}")
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
			return s
		}
	})
}
