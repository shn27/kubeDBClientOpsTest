package main

import (
	"fmt"
	"github.com/spf13/cobra"
	kmapi "kmodules.xyz/client-go/api/v1"
)

var runbook = &cobra.Command{
	Use:   "runbook",
	Short: "Greet the user",
	Long:  `This subcommand greets the user with a custom message.`,
	Run: func(cmd *cobra.Command, args []string) {
		isRunbookCRExit()
	},
}

var markdown1 = &cobra.Command{
	Use:   "markdown",
	Short: "Greet the user",
	Long:  `This subcommand greets the user with a custom message.`,
	Run: func(cmd *cobra.Command, args []string) {
		GetMarkdown()
	},
}

var version = &cobra.Command{
	Use:   "version",
	Short: "Greet the user",
	Long:  `This subcommand greets the user with a custom message.`,
	Run: func(cmd *cobra.Command, args []string) {
		ref := kmapi.TypedObjectReference{
			APIGroup: "kubedb.com", //k8s_group
			Kind:     "MongoDB",    //k8s_kind
		}
		version, err := getPreferredResourceVersion(ref)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("=====================", version)
	},
}
var resource = &cobra.Command{
	Use:   "resource",
	Short: "Greet the user",
	Long:  `This subcommand greets the user with a custom message.`,
	Run: func(cmd *cobra.Command, args []string) {
		getResource("mongodb-down")
	},
}

var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "This is a simple CLI application",
	Long:  `A simple CLI application built with Cobra in Go.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello from Cobra CLI!")
	},
}

func init() {
	rootCmd.AddCommand(markdown1)
	rootCmd.AddCommand(runbook)
	rootCmd.AddCommand(version)
	rootCmd.AddCommand(resource)
}

func main() {
	rootCmd.Execute()
}
