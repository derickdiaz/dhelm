package service

import (
	"os/exec"
	"regexp"
	"strings"

	"golang.org/x/exp/slices"
)

type HelmService struct {}

func NewHelmService() *HelmService {
	return &HelmService{}
}

func (h *HelmService) ListDockerImages(helmChart string) ([]string, error) {
	cmd := exec.Command("helm", "template", helmChart)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	template := string(output[:])
	expression := regexp.MustCompile("image:.+")
	images := expression.FindAllString(template, -1)
	results := []string{}

	for _, image := range images {
		arr := strings.Split(image, " ")
		newImageName := arr[len(arr)-1]
		if slices.Contains(results, newImageName) {
			continue
		}
		results = append(results, newImageName)
	}
	return results, nil 
}