/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"dhelm/service"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List helm file docker images",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}
		helmChart := args[0]	
		helmService := service.NewHelmService()
		dockerService, err := service.NewDockerService()
		if err != nil {
			panic(err)
		}
		images, err := helmService.ListDockerImages(helmChart)
		if err != nil {
			fmt.Printf("Unable to find helm chart: %v\n", helmChart)
			return
		}

		downloadDockerImages, _ := cmd.Flags().GetBool("download")
		platform, _ := cmd.Flags().GetString("platform")
		if downloadDockerImages {
			dockerService.DownloadImages(images, platform)
		}

		repo, _ := cmd.Flags().GetString("retag")
		if repo != "" {
			images, err = dockerService.TagImages(images, repo)
			if err != nil {
				panic(err)
			}
		}

		push, _ := cmd.Flags().GetBool("push")
		if push {
			dockerService.PushImages(images)
		}

		for _, image := range images {
			fmt.Println(image)
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringP("platform", "p", "amd64", "Platform to download docker images")
	listCmd.Flags().BoolP("download", "d", false, "Download Docker Images")
	listCmd.Flags().StringP("retag", "t", "", "Appends string to downloaded Docker images")
	listCmd.Flags().BoolP("push", "a", false, "Pushes Docker Images")
}
