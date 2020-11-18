package cmd

import (
  "fmt"
  "os"
  "path/filepath"

  "github.com/Pale-whale/konfig/internal"
  home "github.com/mitchellh/go-homedir"
  "github.com/spf13/cobra"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
  Use:   "konfig",
  Short: "Kubeconfig manager",
  Long: `konfig is a kubeconfig manager, it's main purpose is to simplify kubernetes context creation/deletion.
With it you can simply add, remove or list contexts and maybe more
`,
}

func Execute() {
  if err := rootCmd.Execute(); err != nil {
    fmt.Println(err)
    os.Exit(1)
  }
}

func init() {
  cobra.OnInitialize()

  // rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

  // Root Flags
  rootCmd.PersistentFlags().StringVarP(&internal.Flags.ConfigPath, "kubeconfig", "k", findDefaultConfig(), "The path of the kubeconfig file to use")

  // Commands
  rootCmd.AddCommand(listCommand)
  rootCmd.AddCommand(rmCommand)
  rootCmd.AddCommand(getCommand)
  rootCmd.AddCommand(backupCommand)
  rootCmd.AddCommand(doctorCommand)
  rootCmd.AddCommand(completionCommand)
}

func findDefaultConfig() string {
  if p, ok := os.LookupEnv("KUBECONFIG"); ok {
    return p
  }

  if home, err := home.Dir(); err == nil {
    p := filepath.Join(home, ".kube", "config")
    _, err := os.Stat(p)
    if err != nil {
      return ""
    }
    return p
  }
  return ""
}
