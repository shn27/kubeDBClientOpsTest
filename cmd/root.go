package cmd

import (
	"fmt"
	"github.com/shn27/Test/work"
	"github.com/spf13/cobra"
	kmapi "kmodules.xyz/client-go/api/v1"
)

var runbook = &cobra.Command{
	Use:   "runbook",
	Short: "Greet the user",
	Long:  `This subcommand greets the user with a custom message.`,
	Run: func(cmd *cobra.Command, args []string) {
		work.IsRunbookCRExit()
	},
}

var markdown1 = &cobra.Command{
	Use:   "markdown",
	Short: "Greet the user",
	Long:  `This subcommand greets the user with a custom message.`,
	Run: func(cmd *cobra.Command, args []string) {
		work.GetMarkdown()
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
		version, err := work.GetPreferredResourceVersion(ref)
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
		work.GetResource("mongodb-down")
	},
}

var RootCmd = &cobra.Command{
	Use:   "app",
	Short: "This is a simple CLI application",
	Long:  `A simple CLI application built with Cobra in Go.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello from Cobra CLI!")
	},
}

func init() {
	RootCmd.AddCommand(markdown1)
	RootCmd.AddCommand(runbook)
	RootCmd.AddCommand(version)
	RootCmd.AddCommand(resource)
}
