/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"

	"github.com/Supkaa/release/internal/git"
	"github.com/Supkaa/release/internal/taggen"
	"github.com/spf13/cobra"
)

// taggenCmd represents the taggen command
var taggenCmd = &cobra.Command{
	Use:   "taggen",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		newTag, err := taggen.GenerateTag()
		if err != nil {
			log.Fatal(err)
		}

		log.Print(newTag.String())
		git.Tag(newTag)
	},
}

func init() {
	rootCmd.AddCommand(taggenCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// taggenCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// taggenCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
