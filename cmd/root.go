package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "drygopher [flags] [<exclusion pattern>...]",
	Short: "Keep your coverage high, and your gopher dry.",
	Long: "For detailed information, visit http://github.com/eltorocorp/drygopher"
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		log.Println("Woot!")
		return nil
	},
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	setupFlags()
}

func setupFlags() {
	rootCmd.DisableFlagsInUseLine = true
}
