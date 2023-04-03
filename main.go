package main

import (
	"flag"
	"fmt"
	"github.com/DragFAQ/uuid-generator/cmd"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"os"
)

var (
	rootCmd = &cobra.Command{
		Use:   "uuid-generator",
		Short: "uuid-generator for reference",
		Long:  "A uuid-generator for reference in repository.",
	}
)

func main() {
	rootCmd.AddCommand(cmd.Run())

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
