package cmd

import (
  "os"

  "github.com/Pale-whale/konfig/internal"
  "github.com/spf13/cobra"
)

var rmCommand = &cobra.Command{
  Use:   "rm",
  Short: "Remove stuff",
  Long: `Rm can remove fields from your kubeconfig.
If rm is executed without arguments, then if ` + "`fzf`" + ` is present you'll be prompted to choose what to delete.
Be careful, any edit made with ` + "`rm`" + `is irreversible and can't be undone unless you have backups.
See the help of the backup command for more informations.`,
  Run:    _rmCommand,
  Aliases: []string{
    "remove",
    "delete",
  },
}

var rmClusterCommand = &cobra.Command{
  Use:   "cluster",
  Short: "Remove cluster",
  Long:  `Removes a cluster entry from your kubeconfig`,
  Run:   _rmClusterCommand,
}

var rmContextCommand = &cobra.Command{
  Use:   "context",
  Short: "Remove contexts",
  Long:  `Removes a context entry from your kubeconfig`,
  Run:   _rmContextCommand,
}

var rmUserCommand = &cobra.Command{
  Use:   "user",
  Short: "Remove users",
  Long:  `Removes a user entry from your kubeconfig`,
  Run:   _rmContextCommand,
}

func init() {
  rmContextCommand.Flags().BoolVarP(&internal.Flags.DeepRm, "deep", "d", false, "Mark dependencies for removal (i.e. cluster and user when removing a context)")

  rmCommand.AddCommand(rmClusterCommand)
  rmCommand.AddCommand(rmContextCommand)
  rmCommand.AddCommand(rmUserCommand)
}

func _rmCommand(cmd *cobra.Command, args []string) {
  selected, err := internal.SelectContext("Select the contexts you want to delete", os.Args[0] + " get context --deep --color {}")
  if err != nil {
    if err.Error() == internal.FzfNotPresent {
      cmd.Help()
      return
    }
    panic(err)
  }

  internal.Flags.DeepRm = true
  _rmContextCommand(cmd, *selected)
}

func _rmClusterCommand(cmd *cobra.Command, args []string) {
  config, err := internal.GetKonfig(internal.Flags.ConfigPath)
  if err != nil {
    panic(err)
  }
  if len(args) == 0 {
    newArgs, err := internal.SelectUser("Select the cluster(s) to remove", os.Args[0] + " get cluster --color {}")
    if err != nil {
      if err.Error() == internal.FzfNotPresent {
        cmd.Help()
        return
      }
      panic(err)
    }
    args = *newArgs
  }

  internal.RemoveClusters(config, args)
  config.OutputReplace()

  config.OutputReplace()
}

func _rmContextCommand(cmd *cobra.Command, args []string) {
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
    newArgs, err := internal.SelectUser("Select the context(s) to remove", os.Args[0] + preview)
    if err != nil {
      if err.Error() == internal.FzfNotPresent {
        cmd.Help()
        return
      }
      panic(err)
    }
    args = *newArgs
  }

  internal.RemoveContexts(config, args)
  config.OutputReplace()
}

func _rmUserCommand(cmd *cobra.Command, args []string) {
  config, err := internal.GetKonfig(internal.Flags.ConfigPath)
  if err != nil {
    panic(err)
  }
  if len(args) == 0 {
    newArgs, err := internal.SelectUser("Select the user(s) to remove", os.Args[0] + " get user --color {}")
    if err != nil {
      if err.Error() == internal.FzfNotPresent {
        cmd.Help()
        return
      }
      panic(err)
    }
    args = *newArgs
  }

  internal.RemoveUsers(config, args)
  config.OutputReplace()
}
