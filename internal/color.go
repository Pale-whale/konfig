package internal

import (
  "fmt"

  "github.com/fatih/color"
  "github.com/goccy/go-yaml/lexer"
  "github.com/goccy/go-yaml/printer"
  "github.com/mattn/go-colorable"
)

const escape = "\x1b"

func format(attr color.Attribute) string {
  return fmt.Sprintf("%s[%dm", escape, attr)
}

func PrintColorized(yaml string) {
  tokens := lexer.Tokenize(yaml)
  var p printer.Printer
  p.LineNumber = false
  p.Bool = func() *printer.Property {
    return &printer.Property{
      Prefix: format(color.FgMagenta),
      Suffix: format(color.Reset),
    }
  }
  p.Number = func() *printer.Property {
    return &printer.Property{
      Prefix: format(color.FgMagenta),
      Suffix: format(color.Reset),
    }
  }
  p.MapKey = func() *printer.Property {
    return &printer.Property{
      Prefix: format(color.FgBlue),
      Suffix: format(color.Reset),
    }
  }
  p.Anchor = func() *printer.Property {
    return &printer.Property{
      Prefix: format(color.FgYellow),
      Suffix: format(color.Reset),
    }
  }
  p.Alias = func() *printer.Property {
    return &printer.Property{
      Prefix: format(color.FgYellow),
      Suffix: format(color.Reset),
    }
  }
  p.String = func() *printer.Property {
    return &printer.Property{
      Prefix: format(color.FgGreen),
      Suffix: format(color.Reset),
    }
  }
  writer := colorable.NewColorableStdout()
  writer.Write([]byte(p.PrintTokens(tokens) + "\n"))
}
