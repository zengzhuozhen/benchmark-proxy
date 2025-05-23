package core

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
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

// TagGenerator 生成器接口，每种tag类型实现该接口
// arg: 解析后的参数字符串，ctx: 线程安全上下文
// 返回生成的字符串
type TagGenerator interface {
	Generate(arg string, ctx *TagContext) string
}

// TagContext 用于存储有状态tag（如incr/range）的上下文，线程安全
// 每个请求/虚拟用户应有独立TagContext
type TagContext struct {
	mu       sync.Mutex
	IncrMap  map[string]int
	RangeMap map[string]int
}

func NewTagContext() *TagContext {
	return &TagContext{
		IncrMap:  make(map[string]int),
		RangeMap: make(map[string]int),
	}
}

// TagRegistry 注册表，支持动态注册/获取tag生成器
type TagRegistry struct {
	mu         sync.RWMutex
	generators map[string]TagGenerator
}

func NewTagRegistry() *TagRegistry {
	return &TagRegistry{generators: make(map[string]TagGenerator)}
}

func (r *TagRegistry) Register(tag string, gen TagGenerator) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.generators[tag] = gen
}

func (r *TagRegistry) Get(tag string) (TagGenerator, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	gen, ok := r.generators[tag]
	return gen, ok
}

// 默认tag注册表，包含常用tag类型
func DefaultTagRegistry() *TagRegistry {
	r := NewTagRegistry()
	r.Register("int", IntGenerator{})
	r.Register("float", FloatGenerator{})
	r.Register("string", StringGenerator{})
	r.Register("bool", BoolGenerator{})
	r.Register("date", DateGenerator{})
	r.Register("timestamp", TimestampGenerator{})
	r.Register("uuid", UUIDGenerator{})
	r.Register("incr", IncrGenerator{})
	r.Register("range", RangeGenerator{})
	r.Register("list", ListGenerator{})
	r.Register("const", ConstGenerator{})
	// 可扩展: email/phone/ip/url/name/ref ...
	return r
}

// TagCompoundParser 负责解析和替换所有${tag}，支持参数化
type TagCompoundParser struct {
	ctx      *TagContext
	registry *TagRegistry
}

func NewTagCompoundParser(registry *TagRegistry) *TagCompoundParser {
	return &TagCompoundParser{
		ctx:      NewTagContext(),
		registry: registry,
	}
}

// tag语法说明：
// ${int}                随机整数
// ${int:100,200}        100~200随机整数
// ${float:1.5,3.5}      1.5~3.5随机浮点
// ${string:8}           8位随机字符串
// ${bool}               随机布尔值
// ${date:2006-01-02}    当前日期，指定格式
// ${timestamp}          当前时间戳
// ${uuid}               UUID
// ${incr:100,2}         从100开始，步长2递增
// ${range:1,100}        1~100递增
// ${list:[a,b,c]}       随机选取a/b/c
// ${const:hello}        常量hello
// ${ref:userId}         引用变量userId（需配合变量池实现）
var defaultParseReg = regexp.MustCompile(`\${(.+?)}`)

// ParseCustomizeTag 替换所有${tag}为生成值
func (p *TagCompoundParser) ParseCustomizeTag(content string) string {
	return defaultParseReg.ReplaceAllStringFunc(content, func(match string) string {
		sub := defaultParseReg.FindStringSubmatch(match)
		if len(sub) != 2 {
			return match
		}
		tagExpr := sub[1]
		var tagType, tagArg string
		if idx := strings.Index(tagExpr, ":"); idx != -1 {
			tagType = tagExpr[:idx]
			tagArg = strings.Trim(tagExpr[idx+1:], "[]")
		} else {
			tagType = tagExpr
			tagArg = ""
		}
		if gen, ok := p.registry.Get(tagType); ok {
			return gen.Generate(tagArg, p.ctx)
		}
		return match
	})
}

// ------------------- 生成器实现 -------------------

// IntGenerator: ${int} 或 ${int:min,max}
type IntGenerator struct{}

func (g IntGenerator) Generate(arg string, _ *TagContext) string {
	if arg == "" {
		return fmt.Sprintf("%d", rand.Int31()>>24)
	}
	parts := strings.Split(arg, ",")
	if len(parts) == 2 {
		min, _ := strconv.Atoi(strings.TrimSpace(parts[0]))
		max, _ := strconv.Atoi(strings.TrimSpace(parts[1]))
		if max > min {
			return fmt.Sprintf("%d", rand.Intn(max-min+1)+min)
		}
	}
	return "0"
}

// FloatGenerator: ${float} 或 ${float:min,max}
type FloatGenerator struct{}

func (g FloatGenerator) Generate(arg string, _ *TagContext) string {
	if arg == "" {
		return fmt.Sprintf("%.2f", rand.Float32())
	}
	parts := strings.Split(arg, ",")
	if len(parts) == 2 {
		min, _ := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
		max, _ := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
		if max > min {
			return fmt.Sprintf("%.2f", min+rand.Float64()*(max-min))
		}
	}
	return "0.0"
}

// StringGenerator: ${string} 或 ${string:length}
type StringGenerator struct{}

func (g StringGenerator) Generate(arg string, _ *TagContext) string {
	length := 10
	if arg != "" {
		if l, err := strconv.Atoi(arg); err == nil && l > 0 {
			length = l
		}
	}
	rand.Seed(time.Now().UnixNano())
	result := make([]byte, length/2)
	rand.Read(result)
	return hex.EncodeToString(result)[:length]
}

// BoolGenerator: ${bool}
type BoolGenerator struct{}

func (g BoolGenerator) Generate(_ string, _ *TagContext) string {
	if rand.Intn(2) == 0 {
		return "false"
	}
	return "true"
}

// DateGenerator: ${date:format}，如${date:2006-01-02}
type DateGenerator struct{}

func (g DateGenerator) Generate(arg string, _ *TagContext) string {
	format := "2006-01-02"
	if arg != "" {
		format = arg
	}
	return time.Now().Format(format)
}

// TimestampGenerator: ${timestamp}
type TimestampGenerator struct{}

func (g TimestampGenerator) Generate(_ string, _ *TagContext) string {
	return fmt.Sprintf("%d", time.Now().Unix())
}

// UUIDGenerator: ${uuid}
type UUIDGenerator struct{}

func (g UUIDGenerator) Generate(_ string, _ *TagContext) string {
	return uuid.New().String()
}

// IncrGenerator: ${incr} 或 ${incr:start,step}
type IncrGenerator struct{}

func (g IncrGenerator) Generate(arg string, ctx *TagContext) string {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	key := arg
	start, step := 1, 1
	if arg != "" {
		parts := strings.Split(arg, ",")
		if len(parts) >= 1 {
			if v, err := strconv.Atoi(strings.TrimSpace(parts[0])); err == nil {
				start = v
			}
		}
		if len(parts) == 2 {
			if v, err := strconv.Atoi(strings.TrimSpace(parts[1])); err == nil {
				step = v
			}
		}
	}
	if ctx.IncrMap[key] == 0 {
		ctx.IncrMap[key] = start
	}
	val := ctx.IncrMap[key]
	ctx.IncrMap[key] += step
	return fmt.Sprintf("%d", val)
}

// RangeGenerator: ${range:start,end}
type RangeGenerator struct{}

func (g RangeGenerator) Generate(arg string, ctx *TagContext) string {
	parts := strings.Split(arg, ",")
	if len(parts) != 2 {
		return ""
	}
	begin, _ := strconv.Atoi(strings.TrimSpace(parts[0]))
	end, _ := strconv.Atoi(strings.TrimSpace(parts[1]))
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	key := arg
	if ctx.RangeMap[key] < begin || ctx.RangeMap[key] > end {
		ctx.RangeMap[key] = begin
	}
	val := ctx.RangeMap[key]
	ctx.RangeMap[key]++
	return fmt.Sprintf("%d", val)
}

// ListGenerator: ${list:[a,b,c]}
type ListGenerator struct{}

func (g ListGenerator) Generate(arg string, _ *TagContext) string {
	items := strings.Split(arg, ",")
	if len(items) == 0 {
		return ""
	}
	rand.Seed(time.Now().UnixNano())
	return strings.TrimSpace(items[rand.Intn(len(items))])
}

// ConstGenerator: ${const:value}
type ConstGenerator struct{}

func (g ConstGenerator) Generate(arg string, _ *TagContext) string {
	return arg
}

// 你可以继续扩展如 EmailGenerator、PhoneGenerator、IPGenerator、NameGenerator、RefGenerator 等
