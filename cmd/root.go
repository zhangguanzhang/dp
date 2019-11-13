package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)


var rootCmd = &cobra.Command{
	Use:   "dp",
	Short: "a docker images tool which without docker daemon",
	Long: `
could pull image without docker daemon ath ths all os.`,
	//Run: func(cmd *cobra.Command, args []string) {
	//
	//},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}


