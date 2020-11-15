package cmd

import (
  "fmt"

  "github.com/Pale-whale/konfig/internal"
  "github.com/spf13/cobra"
)

var listCommand = &cobra.Command{
  Use:   "list",
  Short: "Display kubeconfig",
  Long: `Lists all the entries in the kubeconfig.
This includes the contexts, clusters and users.`,
  Run: _listCommand,

  Aliases: []string{
    "ls",
  },
}

var listClustersCommand = &cobra.Command{
  Use:   "clusters",
  Short: "Display clusters",
  Long:  `Lists all the clusters present in the kubeconfig.`,
  Run:   _listClustersCommand,
}

var listContextsCommand = &cobra.Command{
  Use:   "contexts",
  Short: "Display contexts",
  Long:  `Lists all the contexts present in the kubeconfig.`,
  Run:   _listContextsCommand,
}

var listUsersCommand = &cobra.Command{
  Use:   "users",
  Short: "Display users",
  Long:  `Lists all the users present in the kubeconfig.`,
  Run:   _listUsersCommand,
}

func init() {
  listCommand.AddCommand(listClustersCommand)
  listCommand.AddCommand(listContextsCommand)
  listCommand.AddCommand(listUsersCommand)
}

func _listCommand(cmd *cobra.Command, args []string) {
  config, err := internal.GetKonfig(internal.Flags.ConfigPath)
  if err != nil {
    panic(err)
  }
  str := "Clusters:"
  for _, cluster := range config.Data.Clusters {
    str = fmt.Sprintf("%s\n│ %s", str, cluster.Name)
  }
  str = fmt.Sprintf("%s\n\nContexts:", str)
  for _, context := range config.Data.Contexts {
    str = fmt.Sprintf("%s\n│ %s", str, context.Name)
  }
  str = fmt.Sprintf("%s\n\nUsers:", str)
  for _, user := range config.Data.Users {
    str = fmt.Sprintf("%s\n│ %s", str, user.Name)
  }
  fmt.Println(str)
}

func _listClustersCommand(cmd *cobra.Command, args []string) {
  config, err := internal.GetKonfig(internal.Flags.ConfigPath)
  if err != nil {
    panic(err)
  }
  str := "Clusters:"
  for _, cluster := range config.Data.Clusters {
    str = fmt.Sprintf("%s\n│ %s", str, cluster.Name)
  }
  fmt.Println(str)
}

func _listContextsCommand(cmd *cobra.Command, args []string) {
  config, err := internal.GetKonfig(internal.Flags.ConfigPath)
  if err != nil {
    panic(err)
  }
  str := "Contexts:"
  for _, context := range config.Data.Contexts {
    str = fmt.Sprintf("%s\n│ %s", str, context.Name)
  }
  fmt.Println(str)
}

func _listUsersCommand(cmd *cobra.Command, args []string) {
  config, err := internal.GetKonfig(internal.Flags.ConfigPath)
  if err != nil {
    panic(err)
  }
  str := "Users:"
  for _, user := range config.Data.Users {
    str = fmt.Sprintf("%s\n│ %s", str, user.Name)
  }
  fmt.Println(str)
}
