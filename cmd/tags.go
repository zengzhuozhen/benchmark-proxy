package cmd

import (
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"os"
)

var tagsCmd = &cobra.Command{
	Use:   "tags",
	Short: "Print the `replaceTag` list",
	Run: func(cmd *cobra.Command, args []string) {
		data := [][]string{
			{"int8", "77"},
			{"int16", "19813"},
			{"int32", "1298498081"},
			{"int", "5577006791947779410"},
			{"float", "0.681078"},
			{"float64", "0.604660"},
			{"incr", "1"},
			{"uuid", "d035581b-53a3-48e5-9461-ba24709f06c9"},
			{"string", "a2c44582e086515705c3bc4181c61f9069cf99856dd98ca69c"},
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
