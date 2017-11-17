package cmd

import (
	"os"

	"github.com/gesquive/cli"
	"github.com/gesquive/shield/engine"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// saveCmd represents the save command
var saveCmd = &cobra.Command{
	Use:     "save",
	Aliases: []string{"out", "list", "scribe"},
	Short:   "Output the generated firewall rules",
	Long:    `Generate firewall rules and output them.`,
	Run:     runSave,
}

func init() {
	RootCmd.AddCommand(saveCmd)

	// #viperbug
	// saveCmd.Flags().StringP("rules", "r", "",
	// 	"The templated firewall rules")
	saveCmd.Flags().StringP("output", "o", "-",
		"Output location for generated iptable rules, use '-' for stdout")

	// viper.BindPFlag("rules", saveCmd.Flags().Lookup("rules"))
	viper.BindPFlag("output", saveCmd.Flags().Lookup("output"))

	viper.SetDefault("output", "-")
}

func runSave(cmd *cobra.Command, args []string) {
	rulePath := viper.GetString("rules")

	rules, err := engine.NewRuleset(rulePath)
	if err != nil {
		cli.Error("%v", err)
		os.Exit(2)
	}

	b, err := rules.GenerateRules(displayVersion, []string{})
	if err != nil {
		cli.Error("%v", err)
		os.Exit(2)
	}

	output := viper.GetString("output")
	var pipe *os.File
	if output == "-" {
		pipe = os.Stdout
	} else {
		var err error
		pipe, err = os.OpenFile(output, os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			cli.Error("%v", err)
			os.Exit(2)
		}
		defer pipe.Close()
	}
	pipe.Write(b)
}