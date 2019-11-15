package cmd

import (
	"dp/registry"
	"fmt"
	"github.com/spf13/cobra"
)

var (
	//strict bool
	only bool
)
var checkCmd = &cobra.Command{
	Use:     "check",
	Aliases: []string{"c"},
	Short:   "Check if the images belongs to scheme2.Manifest",

	Example: `
dp c gcr.io/google_containers/bustbox

dp c --only nginx:alpine

dp check nginx:alpine gcr.io/google_containers/pause-amd64:3.1

`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
			return
		}
		var (
			v2 = make([]string, 0)
			v1 = make([]string, 0)

		)
		for _, name := range args {
			p := registry.NewPull(name)
			manifest, _ := p.Manifests()
			if manifest == nil || manifest.SchemaVersion != 2 {
				v1 = append(v1, name)
			} else {
				v2 = append(v2, name)
			}
		}
		fmt.Printf("scheme2.Manifest: %v\n", v2)
		if !only {
			fmt.Printf("scheme1.Manifest: %v\n", v1)
		}
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
	checkCmd.Flags().BoolVarP(&only, "only", "o", false, "only print which is scheme2.Manifest")
}
