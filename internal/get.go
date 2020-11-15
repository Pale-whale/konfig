package internal

import (
  "fmt"
  "gopkg.in/yaml.v3"
)

func GetCluster(config *Konfig, toGet []string) {
  for _, cluster := range config.Data.Clusters {
    if Contains(toGet, cluster.Name) || len(toGet) == 0 {
      d, err := yaml.Marshal(&cluster)
      if err != nil {
        panic(err)
      }
      if Flags.ColoredGet {
        PrintColorized(string(d))
      } else {
        fmt.Println(string(d))
      }
    }
  }
}

func GetContext(config *Konfig, toGet []string) {
  for _, context := range config.Data.Contexts {
    if Contains(toGet, context.Name) || len(toGet) == 0 {
      d, err := yaml.Marshal(&context)
      if err != nil {
        panic(err)
      }
      if Flags.ColoredGet {
        PrintColorized(string(d))
      } else {
        fmt.Println(string(d))
      }
      if Flags.DeepGet == false {
        continue
      }
      GetCluster(config, []string{context.Context.Cluster})
      GetUser(config, []string{context.Context.User})
    }
  }
}

func GetUser(config *Konfig, toGet []string) {
  for _, user := range config.Data.Users {
    if Contains(toGet, user.Name) || len(toGet) == 0 {
      d, err := yaml.Marshal(&user)
      if err != nil {
        panic(err)
      }
      if Flags.ColoredGet {
        PrintColorized(string(d))
      } else {
        fmt.Println(string(d))
      }
    }
  }
}
