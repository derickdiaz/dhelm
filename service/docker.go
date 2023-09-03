package service

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type DockerService struct {
	Client *client.Client
}

func NewDockerService() (*DockerService, error) {
	dockerClient, err := client.NewClientWithOpts()
	if err != nil {
		return nil, err
	}
	return &DockerService{Client: dockerClient}, nil
}

func (d *DockerService) DownloadImages(images []string, platform string) error {
	ctx := context.Background()
	reciever := make(chan struct{})
	for _, image := range images {
		go func(image string, reciever chan<- struct{}) {
			reader, _ := d.Client.ImagePull(ctx, image, types.ImagePullOptions{
				Platform: platform,
			})
			io.Copy(os.Stdout, reader)
			reciever <- struct{}{}
		}(image, reciever)
	}

	for i := 0; i < len(images); i++ {
		<-reciever
	}
	return nil
}

func (d *DockerService) TagImages(images []string, repo string) ([]string, error) {
	ctx := context.Background()
	tags := []string{}
	for _, image := range images {
		arr := strings.Split(image, "/")
		if len(arr) <= 2 {
			newArr := strings.Split(repo, "/")
			newArr = append(newArr, arr...)
			newTag := strings.Join(newArr, "/")
			tags = append(tags, newTag)
			d.Client.ImageTag(ctx, image, newTag)
		}
		ending := strings.Join(arr[1:], "/")
		newTag := fmt.Sprintf("%v/%v", repo, ending)
		tags = append(tags, newTag)
		err := d.Client.ImageTag(ctx, image, newTag)
		if err != nil {
			return nil, err
		}
	}
	return tags, nil
}

func (d *DockerService) PushImages(images []string) error {
	ctx := context.Background()
	reciever := make(chan struct{})
	for _, image := range images {
		go func(image string, reciever chan<- struct{}) {
			reader, err := d.Client.ImagePush(ctx, image, types.ImagePushOptions{})
			if err != nil {
				panic(err)
			}
			io.Copy(os.Stdout, reader)
			reciever<- struct{}{}
		}(image, reciever) 
	}
	for i := 0; i < len(images); i++ {
		<-reciever
	}
	return nil
}