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
		}
	)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewTagCompoundParser().ParseCustomizeTag(tt.args.content)
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
