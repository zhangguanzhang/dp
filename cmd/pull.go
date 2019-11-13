package cmd

import (
	"docker-pull/registry"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strings"
	"time"
)

var (
	//strict bool
	saveName string
)
var cpCmd = &cobra.Command{
	Use:     "pull",
	Aliases: []string{"p"},
	Short:   "pull images",
	Long: `
pull all images and write to a tar.gz file without docker daemon.`,
	Example: `
# pull a image or set the name to save
dp pull nginx:alpine
dp pull -o nginx.tar.gz nginx:alpine

# pull image use sha256
dp pull mcr.microsoft.com/windows/nanoserver@sha256:ae443bd9609b9ef06d21d6caab59505cb78f24a725cc24716d4427e36aedabf2

# pull images and set the name to save
dp pull -o project.tar.gz nginx:alpine nginx:1.17.5-alpine-perl

# pull from different registry 
dp pull -o project.tar.gz nginx:alpine gcr.azk8s.cn/google_containers/pause-amd64:3.1
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
			return
		}
		if len(args) == 1 && saveName == "" {
			saveName = strings.ReplaceAll(args[0], "/", "_")
			saveName = fmt.Sprintf("%s.tar.gz", strings.Replace(args[0], ":", "@", 1))
		}
		// todo regex check
		//for _, name := range args {
		//https://github.com/docker/distribution/blob/master/reference/regexp.go
		//}
		if saveName == "" {
			saveName = fmt.Sprintf("%s.tar.gz", time.Now().Format("2006-1-2-15:04:05"))
		}
		if err := registry.Save(args, saveName);err != nil {
			_ = os.Remove(saveName)
			log.Fatal("Save failed: ", err)
		}
		log.Printf("Successfully written to file %s", saveName)

	},
}

func init() {
	rootCmd.AddCommand(cpCmd)
	//cpCmd.Flags().BoolVarP(&strict, "strict-mode", "s", false,
	//	"The image name of the pull is strictly checked. If it is wrong, it will not be pulled.")
	cpCmd.Flags().StringVarP(&saveName, "out-file", "o", "", "the name will write to,default use timeformat")
}
