package internal

func RemoveClusters(config *Konfig, toDelete []string) {
  for _, td := range toDelete {
    for idx, cluster := range config.Data.Clusters {
      if td == cluster.Name {
        config.Data.Clusters[idx] = config.Data.Clusters[len(config.Data.Clusters) - 1]
        config.Data.Clusters = config.Data.Clusters[:len(config.Data.Clusters) - 1]
        break
      }
    }
  }
}

func RemoveContexts(config *Konfig, toDelete []string) {
  var clusters []string
  var users []string

  for _, td := range toDelete {
    for idx, a := range config.Data.Contexts {
      if td == a.Name {
        config.Data.Contexts[idx] = config.Data.Contexts[len(config.Data.Contexts) - 1]
        config.Data.Contexts = config.Data.Contexts[:len(config.Data.Contexts) - 1]
        break
      }
    }
  }
  if Flags.DeepRm {
    RemoveClusters(config, clusters)
    RemoveUsers(config, users)
  }
}

func RemoveUsers(config *Konfig, toDelete []string) {
  for _, td := range toDelete {
    for idx, user := range config.Data.Users {
      if td == user.Name {
        config.Data.Users[idx] = config.Data.Users[len(config.Data.Users) - 1]
        config.Data.Users = config.Data.Users[:len(config.Data.Users) - 1]
      }
    }
  }
}

