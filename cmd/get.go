package cmd

import (
  "os"

  "github.com/Pale-whale/konfig/internal"
  "github.com/spf13/cobra"
)

var getCommand = &cobra.Command{
  Use:   "get",
  Short: "Get will display informations about a context, cluster or user",
  Long: `Get will display informations about a context cluster or user,
If nothing specified and ` + "`fzf`" + ` is present, get will display informations
about the selected contexts.

You can use the subcommands to get more specific informations.`,
  Run: _getCommand,
}

var getClusterCommand = &cobra.Command{
  Use:   "cluster",
  Short: "",
  Long:  ``,
  Run:   _getClusterCommand,
}

var getContextCommand = &cobra.Command{
  Use:   "context",
  Short: "",
  Long:  ``,
  Run:   _getContextCommand,
}

var getUserCommand = &cobra.Command{
  Use:   "user",
  Short: "",
  Long:  ``,
  Run:   _getUserCommand,
}

func init() {
  getCommand.PersistentFlags().BoolVarP(&internal.Flags.ColoredGet, "color", "c", false, "Use colored output.")
  getContextCommand.Flags().BoolVarP(&internal.Flags.DeepGet, "deep", "d", false, "Also get the info about the user and cluster.")

  getCommand.AddCommand(getClusterCommand)
  getCommand.AddCommand(getContextCommand)
  getCommand.AddCommand(getUserCommand)
}

func _getCommand(cmd *cobra.Command, args []string) {
  selected, err := internal.SelectContext("Select the contexts you want to get", os.Args[0] + " get context --color --deep {}")
  if err != nil {
    if err.Error() == `exec: "fzf": executable file not found in $PATH` {
      cmd.Help()
      return
    }
    panic(err)
  }

  internal.Flags.DeepGet = true
  _getContextCommand(cmd, *selected)
}

func _getClusterCommand(cmd *cobra.Command, args []string) {
  config, err := internal.GetKonfig(internal.Flags.ConfigPath)
  if err != nil {
    panic(err)
  }

  if len(args) == 0 {
    newArgs, err := internal.SelectCluster("Select the cluster to get", os.Args[0] + " get cluster --color {}")
    if err != nil {
      if err.Error() == internal.FzfNotPresent {
        internal.GetCluster(config, args)
        return
      }
      panic(err)
    }
    args = *newArgs
  }
  internal.GetCluster(config, args)
}

func _getContextCommand(cmd *cobra.Command, args []string) {
  config, err := internal.GetKonfig(internal.Flags.ConfigPath)
  if err != nil {
    panic(err)
  }

  if len(args) == 0 {
    preview := " get context --color "
    if internal.Flags.DeepRm {
      preview = preview + "--deep {}"
    } else {
      preview = preview + "{}"
    }
    newArgs, err := internal.SelectContext("Select the cluster to get", os.Args[0] + preview)
    if err != nil {
      if err.Error() == internal.FzfNotPresent {
        internal.GetContext(config, args)
        return
      }
      panic(err)
    }
    args = *newArgs
  }
  internal.GetContext(config, args)
}

func _getUserCommand(cmd *cobra.Command, args []string) {
  config, err := internal.GetKonfig(internal.Flags.ConfigPath)
  if err != nil {
    panic(err)
  }

  if len(args) == 0 {
    newArgs, err := internal.SelectUser("Select the user(s) to get", os.Args[0] + " get user --color {}")
    if err != nil {
      if err.Error() == internal.FzfNotPresent {
        internal.GetUser(config, args)
        return
      }
      panic(err)
    }
    args = *newArgs
  }
  internal.GetUser(config, args)
}
