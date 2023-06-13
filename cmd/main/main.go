package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/containerd/console"

	"dex/internal/container"
	"dex/internal/prompt"
)

func main() {
	containerService, err := container.NewDockerService()
	if err != nil {
		panic(err)
	}
	defer containerService.Close()

	promptService := prompt.NewPromptuiService()

	current := console.Current()
	defer current.Reset()

	if err := current.SetRaw(); err != nil {
		panic(err)
	}

	err = run(containerService, promptService)
	if err != nil {
		fmt.Printf("Something went wrong, got error: %v\n", err)
		os.Exit(1)
	}
}

func run(
	containerService container.ContainerService,
	promptService prompt.PromptService,
) error {
	containers, err := containerService.GetAll()
	if err != nil {
		return err
	}

	promptOptions := containersToPromptOptions(containers)
	selectedOption, err := promptService.DisplaySelect("Select Docker Container", promptOptions)
	if err != nil {
		return nil
	}

	containerID := selectedOption.Value

	command, err := promptService.DisplayPrompt("Command (default is shell - bash or sh)")
	if err != nil {
		return nil
	}

	if command != "" {
		err = containerService.RunCommand(containerID, command)
		return err
	}

	err = containerService.RunCommand(containerID, "bash")
	if err != nil {
		err = containerService.RunCommand(containerID, "sh")
		if err != nil {
			return err
		}
	}

	return nil
}

func containersToPromptOptions(
	dockerContainers []container.Container,
) []prompt.PromptOption {
	var options []prompt.PromptOption

	for _, container := range dockerContainers {
		var names []string
		for _, name := range container.Names {
			names = append(names, strings.Trim(name, "/"))
		}

		label := fmt.Sprintf(
			"[%s] (%s) %s",
			strings.Join(names, ", "),
			container.Image,
			container.ID[0:12],
		)

		option := prompt.PromptOption{
			Label: label,
			Value: container.ID,
		}

		options = append(options, option)
	}

	return options
}
