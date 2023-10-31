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
					content: "{${int}:123}",
				},
				notWant: "{${int}:123}",
			},
			{
				name: "替换${int8}",
				args: args{
					content: "{${int8}:123}",
				},
				notWant: "{${int8}:123}",
			},
			{
				name: "替换${int16}",
				args: args{
					content: "{${int16}:123}",
				},
				notWant: "{${int16}:123}",
			},
			{
				name: "替换${int32}",
				args: args{
					content: "{${int32}:123}",
				},
				notWant: "{${int32}:123}",
			},
			{
				name: "替换${float}",
				args: args{
					content: "{${float}:123}",
				},
				notWant: "{${float}:123}",
			},
			{
				name: "替换${float64}",
				args: args{
					content: "{${float64}:123}",
				},
				notWant: "{${float64}:123}",
			},
			{
				name: "替换${string}",
				args: args{
					content: "{${string}:123}",
				},
				notWant: "{${string}:123}",
			},
			{
				name: "替换${incr}",
				args: args{
					content: "{${incr}:123}",
				},
				notWant: "{${incr}:123}",
			},
			{
				name: "替换${uuid}",
				args: args{
					content: "{${uuid}:123}",
				},
				notWant: "{${uuid}:123}",
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
