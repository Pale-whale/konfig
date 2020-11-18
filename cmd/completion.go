package cmd

import (
  "os"
  "github.com/spf13/cobra"
)

var completionCommand = &cobra.Command{
  Use: "completion [bash|zsh|fish|powershell]",
  Short: "Generate auto-completion",
  Long: `To load completions:

Bash:

$ source <(` + os.Args[0] + ` completion bash)

# To load completions for each session, execute once:
Linux:
  $ ` + os.Args[0] + ` completion bash > /etc/bash_completion.d/` + os.Args[0] + `
MacOS:
  $ ` + os.Args[0] + ` completion bash > /usr/local/etc/bash_completion.d/` + os.Args[0] + `

Zsh:

# If shell completion is not already enabled in your environment you will need
# to enable it.  You can execute the following once:

$ echo "autoload -U compinit; compinit" >> ~/.zshrc

# To load completions for each session, execute once:
$ ` + os.Args[0] + ` completion zsh > "${fpath[1]}/_` + os.Args[0] + `"

# You will need to start a new shell for this setup to take effect.

Fish:

$ ` + os.Args[0] + ` completion fish | source

# To load completions for each session, execute once:
$ ` + os.Args[0] + ` completion fish > ~/.config/fish/completions/` + os.Args[0] + `.fish
`,
  DisableFlagsInUseLine: true,
  ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
  Args:                  cobra.ExactValidArgs(1),
  Run: func(cmd *cobra.Command, args []string) {
    switch args[0] {
    case "bash":
      cmd.Root().GenBashCompletion(os.Stdout)
    case "zsh":
      cmd.Root().GenZshCompletion(os.Stdout)
    case "fish":
      cmd.Root().GenFishCompletion(os.Stdout, true)
    case "powershell":
      cmd.Root().GenPowerShellCompletion(os.Stdout)
    }
  },
}
