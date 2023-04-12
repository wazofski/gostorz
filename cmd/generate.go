/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wazofski/storz/mgen"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate golang model class meta",
	Long: `Scan the provided path for .yaml files and 
generate the corresponding class meta.

For example:
	storz generate model`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Missing argument: model path")
			fmt.Println()
			cmd.Help()
			return
		}

		err := mgen.Generate(args[0])
		if err != nil {
			fmt.Printf("Code-gen failed. %s", err)
			fmt.Println()
		} else {
			fmt.Println("Code-gen complete")
		}
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	// generateCmd.Flags().String("path", "", "Model directory")
}
