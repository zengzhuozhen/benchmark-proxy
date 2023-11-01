package cmd

import (
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"os"
)

var headersCmd = &cobra.Command{
	Use:   "headers",
	Short: "Print the `CustomHeaders` list for http Header",
	Run: func(cmd *cobra.Command, args []string) {
		data := [][]string{
			{"Benchmark-Proxy-Times", "indicate how many times exec in each http request"},
			{"Benchmark-Proxy-Duration", "indicate how much second exec in each http requests"},
			{"Benchmark-Proxy-Concurrency", "concurrency in running"},
			{"Benchmark-Proxy-Check-Result-Status", "indicate the response status to determine whether request is success"},
			{"Benchmark-Proxy-Check-Result-Body", "indicate the response body to determine whether request is success"},
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Header", "Meaning"})
		table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
		table.SetColWidth(200)
		for _, v := range data {
			table.Append(v)
		}
		table.Render() // Send output
	},
}
