package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/containerd/console"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/manifoldco/promptui"
)

func main() {
	dockerClient := getDockerClient()
	defer dockerClient.Close()

	dockerContainers := getDockerContainers(dockerClient)
	items, keys := getPromptSelectItems(dockerContainers)

	selectedKey := displayPromptSelect(keys)
	if selectedKey == "" {
		return
	}

	selectedContainer := items[selectedKey]

	command := displayPromptCommand()

	current := console.Current()
	defer current.Reset()

	if err := current.SetRaw(); err != nil {
		panic(err)
	}

	if command != "" {
		dockerExec(dockerClient, selectedContainer.ID, []string{command})
		return
	}

	err := dockerExec(dockerClient, selectedContainer.ID, []string{"bash"})
	if err != nil {
		err = dockerExec(dockerClient, selectedContainer.ID, []string{"sh"})
		if err != nil {
			panic(err)
		}
	}
}

func getDockerClient() *client.Client {
	client, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	return client
}

func getDockerContainers(client *client.Client) []types.Container {
	ctx := context.Background()
	containers, err := client.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	return containers
}

func getPromptSelectItems(
	dockerContainers []types.Container,
) (map[string]types.Container, []string) {
	items := make(map[string]types.Container)
	keys := []string{}

	for _, container := range dockerContainers {
		var names []string
		for _, name := range container.Names {
			names = append(names, strings.Trim(name, "/"))
		}

		key := fmt.Sprintf(
			"[%s] (%s) %s",
			strings.Join(names, ", "),
			container.Image,
			container.ID[0:12],
		)
		keys = append(keys, key)
		items[key] = container
	}

	return items, keys
}

func displayPromptSelect(items []string) string {
	prompt := promptui.Select{
		Label: "Select Docker Container",
		Items: items,
	}
	_, result, err := prompt.Run()
	if err != nil {
		return ""
	}

	return result
}

func displayPromptCommand() string {
	prompt := promptui.Prompt{
		Label: "Command (default is shell - bash / sh)",
	}
	result, err := prompt.Run()
	if err != nil {
		return ""
	}

	return result
}

func dockerExec(
	dockerClient *client.Client,
	containerID string,
	command []string,
) error {
	ctx := context.Background()
	config := types.ExecConfig{
		AttachStdin:  true,
		AttachStderr: true,
		AttachStdout: true,
		Tty:          true,
		Detach:       false,
		Cmd:          command,
	}

	exec, err := dockerClient.ContainerExecCreate(ctx, containerID, config)
	if err != nil {
		return err
	}

	execAttachConfig := types.ExecStartCheck{
		Tty:    true,
		Detach: false,
	}

	attach, err := dockerClient.ContainerExecAttach(ctx, exec.ID, execAttachConfig)
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
