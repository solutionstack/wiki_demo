package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "Wikipedia Demo",
}

func init() {
	rootCmd.AddCommand(serverCmd)
}

// Execute : runs registered cobra commands
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		Log.Fatal().Err(err).Msg("")
	}
}

var serverCmd = &cobra.Command{
	Use:   "demo",
	Short: "Starts wiki demo server",
	Run:   startServer,
}

func startServer(cmd *cobra.Command, args []string) {

	err := StartNew()
	if err != nil {
		Log.Fatal().Err(err).Msg("fatal startup error")
	}

}
