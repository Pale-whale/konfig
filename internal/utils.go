package internal

import (
  "io"
  "os"
  "os/exec"
  "strings"
)

const (
  FzfNotPresent string = `exec: "fzf": executable file not found in $PATH`
)

func DeleteStrArray(ary *[]string, toDelete []string) *[]string {
  for _, td := range toDelete {
    for idx, a := range *ary {
      if a == td {
        (*ary)[idx] = (*ary)[len(*ary) - 1]
        *ary = (*ary)[:len(*ary) - 1]
        break ;
      }
    }
  }
  return ary
}

func Contains(arr []string, str string) bool {
   for _, a := range arr {
      if a == str {
         return true
      }
   }
   return false
}

func SelectCluster(header string, preview string) (*[]string, error) {
  return _select(0, header, preview)
}

func SelectContext(header string, preview string) (*[]string, error) {
  return _select(1, header, preview)
}

func SelectUser(header string, preview string) (*[]string, error) {
  return _select(2, header, preview)
}

func _select(what int, header string, preview string) (*[]string, error) {
  fzfPath, err := exec.LookPath("fzf")
  if err != nil {
    return nil, err
  }
  command := exec.Command(fzfPath, "-m", "--header", header, "--preview", preview, "--preview-window=right:70%:wrap")
  command.Stderr = os.Stderr
  command.Env = os.Environ()
  cmdStdIn, err := command.StdinPipe()
  if err != nil {
    return nil, err
  }

  go func() {
    config, err := GetKonfig(Flags.ConfigPath)
    if err != nil {
      panic(err)
    }
    switch what {
      case 0:
        for _, cluster := range config.Data.Clusters {
          io.WriteString(cmdStdIn, cluster.Name)
          io.WriteString(cmdStdIn, "\n")
        }
      case 1:
        for _, context := range config.Data.Contexts {
          io.WriteString(cmdStdIn, context.Name)
          io.WriteString(cmdStdIn, "\n")
        }
      case 2:
        for _, user := range config.Data.Users {
          io.WriteString(cmdStdIn, user.Name)
          io.WriteString(cmdStdIn, "\n")
        }
    }
    cmdStdIn.Close()
  }()

  selected, err := command.Output()
  if err != nil {
    return nil, err
  }
  splitted := strings.Split(string(selected), "\n")
  return &splitted, nil
}
