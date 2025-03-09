/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Supkaa/release/internal/commit"
	"github.com/Supkaa/release/internal/conventional/types"
	"github.com/Supkaa/release/internal/git"
	"github.com/Supkaa/release/internal/linter"
	"github.com/spf13/cobra"
)

// commitgenCmd represents the commitgen command
var commitgenCmd = &cobra.Command{
	Use:   "commitgen",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !git.IsGitAddExecuted() {
			log.Fatal("git add not executed")
		}

		reader := bufio.NewReader(os.Stdin)
		commit := commit.Commit{
			IsMerge: false,
		}
		commit.Type = selectCommitTypeommit(reader)
		commit.Scope = enterScope(reader)
		commit.Message = enterShortDescription(reader)

		if err := linter.LintCommit(commit); err != nil {
			log.Fatal(err)
		}

		git.Commit(commit)
	},
}

func init() {
	rootCmd.AddCommand(commitgenCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// commitgenCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// commitgenCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func selectCommitTypeommit(reader *bufio.Reader) types.CommitType {
	fmt.Println("Select commit type:")
	for i, t := range types.ValidCommitTypes {
		fmt.Printf("%d: %s\n", i+1, t)
	}

	typeChoice := readChoice(reader, len(types.ValidCommitTypes))

	return types.ValidCommitTypes[typeChoice-1]
}

func enterScope(reader *bufio.Reader) string {
	fmt.Println("Enter a scope:")
	scope, _ := reader.ReadString('\n')

	return strings.TrimSpace(scope)
}

func enterShortDescription(reader *bufio.Reader) string {
	fmt.Print("Enter a short description: ")
	message, _ := reader.ReadString('\n')

	return strings.TrimSpace(message)
}

func readChoice(reader *bufio.Reader, max int) int {
	for {
		fmt.Print("Enter a number: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		for i := 1; i <= max; i++ {
			if fmt.Sprintf("%d", i) == input {
				return i
			}
		}

		fmt.Println("❌ Invalid input! Please try again.")
	}
}
