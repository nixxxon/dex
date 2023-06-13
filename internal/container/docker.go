package container

import (
	"context"
	"encoding/csv"
	"io"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type DockerService struct {
	client *client.Client
}

func (s *DockerService) GetAll() ([]Container, error) {
	var containers []Container

	ctx := context.Background()
	dockerContainers, err := s.client.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return containers, err
	}

	for _, dockerContainer := range dockerContainers {
		containers = append(containers, Container{
			ID:    dockerContainer.ID,
			Names: dockerContainer.Names,
			Image: dockerContainer.Image,
		})
	}
	return containers, nil
}

func (s *DockerService) RunCommand(containerID string, command string) error {
	arguments, err := splitCommand(command)
	if err != nil {
		return err
	}

	ctx := context.Background()
	execConfig := types.ExecConfig{
		AttachStdin:  true,
		AttachStderr: true,
		AttachStdout: true,
		Tty:          true,
		Detach:       false,
		Cmd:          arguments,
	}

	exec, err := s.client.ContainerExecCreate(ctx, containerID, execConfig)
	if err != nil {
		return err
	}

	execAttachConfig := types.ExecStartCheck{
		Tty:    true,
		Detach: false,
	}

	attach, err := s.client.ContainerExecAttach(ctx, exec.ID, execAttachConfig)
	if err != nil {
		return err
	}
	defer attach.Close()

	errChan := make(chan error, 1)

	go func() {
		_, err = io.Copy(os.Stdout, attach.Reader)
		errChan <- err
	}()

	go func() {
		_, err := io.Copy(attach.Conn, os.Stdin)
		errChan <- err
	}()

	err = <-errChan
	return err
}

func (s *DockerService) Close() {
	s.client.Close()
}

func splitCommand(command string) ([]string, error) {
	r := csv.NewReader(strings.NewReader(command))
	r.Comma = ' '
	arguments, err := r.Read()
	if err != nil {
		return []string{}, nil
	}
	return arguments, nil
}

func NewDockerService() (ContainerService, error) {
	dockerClient, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		return nil, err
	}
	return &DockerService{client: dockerClient}, nil
}
