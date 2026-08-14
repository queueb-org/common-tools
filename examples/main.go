package main

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"common.queueb.org/tools"
)

var (
	// flags
	password = &tools.Sensitive{}

	// command
	rootCmd = &cobra.Command{
		Use:  "example",
		RunE: appRun,
	}
)

// CLI Arguments settings.
const (
	PasswordFlag         = "password"
	PasswordShortFlag    = "p"
	PasswordDefaultValue = ""
)

func addFlags(flags *pflag.FlagSet) {
	env := tools.MakeEnv("APP")

	// default value exists, however, it can be used only in-app,
	// and it will be hidden from standard --help usage.
	p := flags.VarPF(password, PasswordFlag, PasswordShortFlag,
		tools.EnvUsageP(env, PasswordFlag, "password to set"),
	)
	_ = p.Value.Set(tools.EnvP(env, PasswordFlag, PasswordDefaultValue))
}

func main() {
	addFlags(rootCmd.PersistentFlags())

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("got issue: %v", err)
	}
}

func appRun(cmd *cobra.Command, args []string) (err error) {
	fmt.Printf("your password: %v\n", password.Reveal())
	fmt.Printf("also try --help\n")
	return
}
