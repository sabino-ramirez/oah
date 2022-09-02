/*
Copyright © 2022 Sabino Ramirez <sabinoramirez017@gmail.com>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Test the endpoints with parameters entered in setup",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("test called")
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
}
