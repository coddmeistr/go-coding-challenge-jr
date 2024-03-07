package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Root command for all cli application
var rootCmd = &cobra.Command{
	Use:   "",
	Short: "",
	Long:  ``,
	Run:   func(cmd *cobra.Command, args []string) {},
}

// Persistent flag, every command inherits this flag
var address string

func Execute() {
	rootCmd.PersistentFlags().StringVarP(&address, "address", "a", "localhost:6000", "gRPC server address")
	if err := rootCmd.Execute(); err != nil {
		_, err := fmt.Fprintf(os.Stderr, "Error, when executing cli: %s", err)
		if err != nil {
			return
		}
		os.Exit(1)
	}
}
