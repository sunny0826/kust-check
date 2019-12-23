package cmd

import (
  "fmt"
  "os"

  "github.com/spf13/cobra"
)

type Cli struct {
  rootCmd *cobra.Command
}

//NewCli returns the cli instance used to register and execute command
func NewCli() *Cli {
  cli := &Cli{
    rootCmd: &cobra.Command{
      Use:   "kust-check",
      Short: "chek kustomize docs",
      Long:  "chek kustomize docs",
    },
  }
  cli.rootCmd.SetOutput(os.Stdout)
  cli.setFlags()
  return cli
}

// setFlags defines flags for root command
func (cli *Cli) setFlags() {
  flags := cli.rootCmd.PersistentFlags()
  fmt.Println(flags)
}

//Run command
func (cli *Cli) Run() error {
  return cli.rootCmd.Execute()
}
