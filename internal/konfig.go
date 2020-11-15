package internal

import (
  "io"
  "io/ioutil"
  "os"

  "gopkg.in/yaml.v3"
)

type konfig struct {
  ApiVersion     string `yaml:"apiVersion,omitempty"`
  CurrentContext string `yaml:"current-context,omitempty"`
  Kind           string `yaml:",omitempty"`
  Preferences    map[interface{}]interface{} // ?
  Clusters       []struct {
    Name    string `yaml:",omitempty"`
    Cluster struct {
      CAData string `yaml:"certificate-authority-data,omitempty"`
      Server string `yaml:",omitempty"`
    }
  } `yaml:",omitempty"`
  Contexts        []struct {
    Name          string `yaml:",omitempty"`
    Context       struct {
      Cluster     string `yaml:",omitempty"`
      User        string `yaml:",omitempty"`
      Namespace   string `yaml:",omitempty"`
    }
  } `yaml:",omitempty"`
  Users []struct {
    Name string `yaml:",omitempty"`
    User struct {
      CCData   string `yaml:"client-certificate-data,omitempty"`
      CKData   string `yaml:"client-key-data,omitempty"`
      Token    string `yaml:"token,omitempty"`
      Username string `yaml:"username,omitempty"`
      Password string `yaml:"password,omitempty"`
    } `yaml:"user,omitempty"`
  } `yaml:"users,omitempty"`
}

type Konfig struct {
  Path string
  Data konfig
}

func GetKonfig(path string) (*Konfig, error) {
  conf := Konfig{}
  conf.Path = path

  data, err := ioutil.ReadFile(path)
  if err != nil {
    return nil, err
  }

  err = yaml.Unmarshal(data, &conf.Data)
  if err != nil {
    return nil, err
  }
  return &conf, nil
}

func (konfig *Konfig) OutputReplace() {
  if os.Remove(konfig.Path) != nil {
    konfig.OutputNew(".failed." + konfig.Path)
    return
  }
  konfig.OutputNew(konfig.Path)
}

func (konfig *Konfig) OutputNew(filePath string) {
  f, err := os.Create(filePath)
  if err != nil {
    panic(err)
  }
  konfig.Output(f)
}

func (konfig *Konfig) Output(filePath io.Writer) {
  pretty, err := yaml.Marshal(&konfig.Data)
  if err != nil {
    panic(err)
  }
  io.WriteString(filePath, string(pretty))
}
