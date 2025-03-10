/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"

	"github.com/Supkaa/release/internal/git"
	"github.com/Supkaa/release/internal/linter"
	"github.com/spf13/cobra"
)

// lintCmd represents the lint command
var lintCmd = &cobra.Command{
	Use:   "lint",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		lastCommit, err := git.GetLatestCommit()
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("%#v", lastCommit)
		if err := linter.LintCommit(lastCommit); err != nil {
			log.Fatal(err)
		}

		log.Print(lastCommit.String())
	},
}

func init() {
	rootCmd.AddCommand(lintCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// lintCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// lintCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
