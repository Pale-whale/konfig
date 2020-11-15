package cmd

import (
  "fmt"
  "io/ioutil"
  "os"
  "path"
  "strconv"
  "strings"
  "sort"

  "github.com/Pale-whale/konfig/internal"
  "github.com/spf13/cobra"
)

var backupCommand = &cobra.Command{
  Use:   "backup",
  Short: "Backup kubeconfig",
  Long: `Backup will create a backup of your kubeconfig, a path can be specified thanks
to the ` + "`--path`" + ` option.
The backup will be created by appending .backup and an increment to the original path`,
  Run: _backupCommand,
}

func init() {
  backupCommand.Flags().StringVarP(&internal.Flags.BackupPath, "backup-path", "b", "", "Custom path for backups location")
}

// TODO: Just call cp binary
func _backupCommand(cmd *cobra.Command, args []string) {
  config, err := internal.GetKonfig(internal.Flags.ConfigPath)
  if err != nil {
    panic(err)
  }

  var newPath string
  if internal.Flags.BackupPath != "" {
    newPath = internal.Flags.BackupPath // TODO: Roll is buggy
  } else {
    newPath = config.Path + ".backup.0"
  }
  if _, err := os.Stat(newPath); !os.IsNotExist(err) {
    files, err := ioutil.ReadDir(path.Dir(newPath))
    if err != nil {
      panic(err)
    }

    var takens []string
    for _, f := range files {
      if strings.HasPrefix(path.Base(f.Name()), path.Base(newPath)) {
        takens = append(takens, f.Name())
      }
    }
    sort.Strings(takens)
    np := strings.Split(takens[len(takens) - 1], ".")
    id, err := strconv.Atoi(np[len(np) - 1])
    np[len(np) - 1] = strconv.Itoa(id + 1)
    newPath = path.Dir(newPath) + "/" + strings.Join(np, ".")
  }
  config.OutputNew(newPath)
  fmt.Println(newPath, "wrote successfully")
}
