package gen

import "github.com/spf13/cobra"

var ServiceCmd = &cobra.Command{
	Use:               "service",
	Short:             "gen service",
	Long:              "gen service",
	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	Run:               GeneratorModel,
}

func init() {
	ModelCmd.Flags().StringVarP(&table, "table", "t", "", "[Required] The name of the db table name")
	ModelCmd.Flags().StringVarP(&structName, "struct", "s", "", "[Required] The name of the db table model")
}
