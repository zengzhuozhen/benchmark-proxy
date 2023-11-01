package cmd

import (
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"os"
)

var tagsCmd = &cobra.Command{
	Use:   "tags",
	Short: "Print the `replaceTag` list for http params",
	Run: func(cmd *cobra.Command, args []string) {
		data := [][]string{
			{"int", "77"},
			{"float", "0.94"},
			{"incr", "1"},
			{"string", "88fa7ac2bf"},
			{"uuid", "d035581b-53a3-48e5-9461-ba24709f06c9"},
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Tag", "Example"})
		table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
		for _, v := range data {
			table.Append(v)
		}
		table.Render() // Send output
	},
}
