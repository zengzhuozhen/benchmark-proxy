package core

import (
	"fmt"
	"testing"
)

func TestParseCustomizeTag(t *testing.T) {
	type args struct {
		content string
	}
	var (
		tests = []struct {
			name    string
			args    args
			notWant string
			want    string
		}{
			{
				name: "不替换",
				args: args{
					content: "不替换",
				},
				want: "不替换",
			},
			{
				name: "替换${int}",
				args: args{
					content: "{${int}}",
				},
				notWant: "{${int}}",
			},
			{
				name: "替换${float}",
				args: args{
					content: "{${float}}",
				},
				notWant: "{${float}}",
			},
			{
				name: "替换${string}",
				args: args{
					content: "{${string}}",
				},
				notWant: "{${string}}",
			},
			{
				name: "替换${incr}",
				args: args{
					content: "{${incr}}",
				},
				notWant: "{${incr}}",
			},
			{
				name: "替换${uuid}",
				args: args{
					content: "{${uuid}}",
				},
				notWant: "{${uuid}}",
			},
			{
				name: "替换${list:[30,60]}",
				args: args{
					content: "{${list:[30,60]}}",
				},
				notWant: "{${list:[30,60]}}",
			},
			{
				name: "替换${range:[30,60]}",
				args: args{
					content: "{${range:[30,60]}}",
				},
				notWant: "{${range:[30,60]}}",
			},
			{
				name: "int范围",
				args: args{
					content: "${int:10,20}",
				},
				notWant: "${int:10,20}",
			},
			{
				name: "float范围",
				args: args{
					content: "${float:1.5,2.5}",
				},
				notWant: "${float:1.5,2.5}",
			},
			{
				name: "string长度",
				args: args{
					content: "${string:12}",
				},
				notWant: "${string:12}",
			},
			{
				name: "bool类型",
				args: args{
					content: "${bool}",
				},
				notWant: "${bool}",
			},
			{
				name: "date格式",
				args: args{
					content: "${date:2006-01-02}",
				},
				notWant: "${date:2006-01-02}",
			},
			{
				name: "timestamp",
				args: args{
					content: "${timestamp}",
				},
				notWant: "${timestamp}",
			},
			{
				name: "const常量",
				args: args{
					content: "${const:hello}",
				},
				want: "hello",
			},
			{
				name: "list多选",
				args: args{
					content: "${list:[a,b,c]}",
				},
				notWant: "${list:[a,b,c]}",
			},
			{
				name: "incr递增",
				args: args{
					content: "${incr:100,2}",
				},
				notWant: "${incr:100,2}",
			},
			{
				name: "range递增",
				args: args{
					content: "${range:5,7}",
				},
				notWant: "${range:5,7}",
			},
			{
				name: "uuid唯一",
				args: args{
					content: "${uuid}",
				},
				notWant: "${uuid}",
			},
		}
	)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewTagCompoundParser(DefaultTagRegistry()).ParseCustomizeTag(tt.args.content)
			if tt.want != "" && got != tt.want {
				t.Errorf("ParseCustomizeTag() = %s, want %s,", got, tt.notWant)
			} else if tt.notWant != "" && got == tt.notWant {
				t.Errorf("ParseCustomizeTag() = %s, notWant %s,", got, tt.notWant)
			} else {
				fmt.Println("got:", got)
			}
		})
	}
}
