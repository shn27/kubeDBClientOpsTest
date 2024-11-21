package work_postgres

import (
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
)

var PgCmdTest2 = &cobra.Command{
	Use:   "pgTestAll",
	Short: "This is a simple CLI application",
	Long:  `Test postgres cmd2`,
	Run: func(cmd *cobra.Command, args []string) {
		klog.Info("Testing 'TestPostgresServerStatus()':\n")
		TestPostgresServerStatus()
		klog.Info("Testing 'TestClientFuncs()':\n")
		TestClientFuncs()
		klog.Info("Testing 'TestSharedBuffers()':\n")
		TestSharedBuffers()
	},
}

var PgCmdTest = &cobra.Command{
	Use:   "pgCmdTest",
	Short: "This is a simple CLI application",
	Long:  `Test postgres cmd`,
	Run: func(cmd *cobra.Command, args []string) {
		TestPostgresServerStatus()
	},
}

var PgCmdTestSharedBuffers = &cobra.Command{
	Use:   "pgTestSharedBuffers",
	Short: "This is a simple CLI Application",
	Long:  `Test postgres cmd`,
	Run: func(cmd *cobra.Command, args []string) {
		klog.Info("Testing `TestSharedBuffers()`:\n")
		TestSharedBuffers()
	},
}
