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
		}{
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
		}
	)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTagCompoundParser().ParseCustomizeTag(tt.args.content); got == tt.notWant {
				t.Errorf("ParseCustomizeTag() = %v, notWant %v", got, tt.notWant)
			} else {
				fmt.Println("got:", got)
			}
		})
	}
}
