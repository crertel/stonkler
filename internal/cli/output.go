package cli

type outputFormat string

const (
	outputTable outputFormat = "table"
	outputJSON  outputFormat = "json"
	outputCSV   outputFormat = "csv"
)
