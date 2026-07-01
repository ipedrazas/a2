package cmd

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// CommandHelp is the JSON schema for a CLI command, intended for AI agent consumption.
type CommandHelp struct {
	Name        string        `json:"name"`
	Use         string        `json:"use"`
	Short       string        `json:"short"`
	Long        string        `json:"long,omitempty"`
	Aliases     []string      `json:"aliases,omitempty"`
	Flags       []FlagHelp    `json:"flags,omitempty"`
	Subcommands []CommandHelp `json:"subcommands,omitempty"`
}

// FlagHelp is the JSON schema for a CLI flag.
type FlagHelp struct {
	Name      string `json:"name"`
	Shorthand string `json:"shorthand,omitempty"`
	Usage     string `json:"usage"`
	Default   string `json:"default,omitempty"`
	Type      string `json:"type"`
}

var skipHelpFlags = map[string]bool{
	"help":      true,
	"help-json": true,
	"version":   true,
}

func buildCommandHelp(cmd *cobra.Command) CommandHelp {
	h := CommandHelp{
		Name:  cmd.Name(),
		Use:   cmd.Use,
		Short: cmd.Short,
		Long:  strings.TrimSpace(cmd.Long),
	}
	if len(cmd.Aliases) > 0 {
		h.Aliases = cmd.Aliases
	}

	seen := map[string]bool{}
	addFlag := func(f *pflag.Flag) {
		if skipHelpFlags[f.Name] || seen[f.Name] {
			return
		}
		seen[f.Name] = true
		h.Flags = append(h.Flags, FlagHelp{
			Name:      f.Name,
			Shorthand: f.Shorthand,
			Usage:     f.Usage,
			Default:   f.DefValue,
			Type:      f.Value.Type(),
		})
	}

	cmd.Flags().VisitAll(addFlag)
	cmd.PersistentFlags().VisitAll(addFlag)

	for _, sub := range cmd.Commands() {
		if sub.Hidden || sub.Name() == "help" {
			continue
		}
		h.Subcommands = append(h.Subcommands, buildCommandHelp(sub))
	}

	return h
}

func outputHelpJSON(cmd *cobra.Command) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(buildCommandHelp(cmd))
	os.Exit(0)
}

// checkHelpJSON inspects os.Args for --help-json before cobra parses flags.
// If found, it resolves the target command from the non-flag positional args,
// prints its JSON help, and exits.
func checkHelpJSON() {
	args := os.Args[1:]
	found := false
	for _, a := range args {
		if a == "--help-json" {
			found = true
			break
		}
	}
	if !found {
		return
	}

	var cmdArgs []string
	for _, a := range args {
		if a == "--help-json" || strings.HasPrefix(a, "-") {
			continue
		}
		cmdArgs = append(cmdArgs, a)
	}

	targetCmd, _, _ := rootCmd.Find(cmdArgs)
	if targetCmd == nil {
		targetCmd = rootCmd
	}
	outputHelpJSON(targetCmd)
}
