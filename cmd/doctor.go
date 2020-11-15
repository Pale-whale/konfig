package cmd

import (
  "fmt"
  "io"
  "os"
  "os/exec"
  "bufio"
  "sync"
  "strings"
  "context"
  "encoding/base64"

  "github.com/Pale-whale/konfig/internal"
  "github.com/spf13/cobra"

  "k8s.io/client-go/kubernetes"
  metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
  "k8s.io/client-go/rest"
)

var doctorCommand = &cobra.Command{
  Use:   "doctor",
  Short: "Test clusters and remove not working ones",
  Long: `Doctor will try to connect to each of your clusters
and will mark for deletions when it can't, you can edit which one
you want to delete if you have ` + "`fzf`" + ` installed.`,
  Run: _doctorCommand,

  Aliases: []string{
    "doc",
  },
}

func init() {
}

func _doctorCommand(cmd *cobra.Command, args []string) {
  var mainLock, fastLock sync.WaitGroup
  config, err := internal .GetKonfig(internal.Flags.ConfigPath)
  if err != nil {
    panic(err)
  }
  var toDelete []string
  usedClusters := new([]string)
  usedUsers := new([]string)
  mainLock.Add(len(config.Data.Contexts))
  fastLock.Add(len(config.Data.Contexts))
  for _, ctx := range config.Data.Contexts {
    go func(clusterName string, userName string, contextName string) {
      conf := rest.Config {}
      clusterOK := false
      userOK := false
      for _, cluster := range config.Data.Clusters {
        if cluster.Name != clusterName {
          continue
        }
        clusterOK = true
        *usedClusters = append(*usedClusters, clusterName)
        conf.Host = cluster.Cluster.Server
        cadata, err := base64.StdEncoding.DecodeString(cluster.Cluster.CAData)
        if err != nil {
          panic(err)
        }
        conf.TLSClientConfig.CAData = cadata
      }
      for _, user := range config.Data.Users {
        if user.Name != userName {
          continue
        }
        userOK = true
        *usedUsers = append(*usedUsers, userName)
        conf.Username = user.User.Username
        conf.Password = user.User.Password
        conf.BearerToken = user.User.Token
        ccdata, err := base64.StdEncoding.DecodeString(user.User.CCData)
        if err != nil {
          panic(err)
        }
        conf.TLSClientConfig.CertData = ccdata
        ckdata, err := base64.StdEncoding.DecodeString(user.User.CKData)
        if err != nil {
          panic(err)
        }
        conf.TLSClientConfig.KeyData = ckdata
      }
      fastLock.Done()
      if !clusterOK || !userOK {
        toDelete = append(toDelete, contextName)
        fmt.Printf("Context %s have bad user or cluster\n", contextName)
        mainLock.Done()
        return
      }
      clientset, err := kubernetes.NewForConfig(&conf)
      if err != nil {
        panic(err)
      }
      fmt.Printf("Checking %s @ %s ...\n", contextName, conf.Host)
      _, err = clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
      if err != nil {
        toDelete = append(toDelete, contextName)
      }
      mainLock.Done()
    }(ctx.Context.Cluster, ctx.Context.User, ctx.Name)
  }
  fastLock.Wait()
  var orphanedClusters []string
  var orphanedUsers []string
  for _, cluster := range config.Data.Clusters {
    if internal.Contains(*usedClusters, cluster.Name) {
      continue
    }
    orphanedClusters = append(orphanedClusters, cluster.Name)
  }
  for _, user := range config.Data.Users {
    if internal.Contains(*usedUsers, user.Name) {
      continue
    }
    orphanedUsers = append(orphanedClusters, user.Name)
  }
  mainLock.Wait()

  toDelete = promptDelete("Can't connect to %d contexts:", toDelete, "Select the contexts you DON'T want to delete")

  internal.Flags.DeepRm = true
  internal.RemoveContexts(config, toDelete)

  orphanedClusters = promptDelete("Found %d orphaned clusters:", orphanedClusters, "Select the clusters you DON'T want to delete")
  internal.RemoveClusters(config, orphanedClusters)

  orphanedUsers = promptDelete("Found %d orphaned users:", orphanedUsers, "Select the users you DON'T want to delete")
  internal.RemoveUsers(config, orphanedUsers)

  config.OutputReplace()

}

func promptDelete(baseStr string, toDelete []string, header string) []string {
  str := fmt.Sprintf(baseStr, len(toDelete))
  for _, c := range toDelete {
    str = fmt.Sprintf("%s\n - %s", str, c)
  }
  fmt.Println(str)
  pth, err := exec.LookPath("fzf")
  if err == nil {
    toDelete = editDelete(toDelete, pth, header)
  } else {
    confirm, err := promptUser("\nWould you like to continue ? [y/N]")
    if err != nil {
      toDelete = []string{}
    }
    if confirm != 'y' {
      toDelete = []string{}
    }
  }
  return toDelete
}

func editDelete(toDelete []string, fzfPath string, header string) []string {
  confirm, err := promptUser("\nWould you like to continue ? [y/N/e]")
  if err != nil {
    return []string{}
  }
  if confirm == 'y' {
    return toDelete
  } else if confirm == 'e' {
    command := exec.Command(fzfPath, "-m", "--header", header)
    command.Stderr = os.Stderr
    command.Env = os.Environ()
    cmdStdIn, err := command.StdinPipe()
    if err != nil {
      panic(err)
    }
    go func() {
      for _, context := range toDelete {
        io.WriteString(cmdStdIn, context)
        io.WriteString(cmdStdIn, "\n")
      }
      cmdStdIn.Close()
    }()

    selected, err := command.Output()
    if err != nil {
      panic(err)
    }
    edited := strings.Split(string(selected), "\n")
    if len(edited) > 0 {
      str := "Will keep: "
      for _, c := range edited[:len(edited) - 1] {
        str = fmt.Sprintf("%s\n - %s", str, c)
      }
      fmt.Println(str)
      internal.DeleteStrArray(&toDelete, edited)
    }
    return toDelete
  } else {
    return []string{}
  }
}

func promptUser(prompt string) (rune, error) {
  fmt.Println(prompt)
  reader := bufio.NewReader(os.Stdin)
  confirm, _, err := reader.ReadRune()
  return confirm, err
}


