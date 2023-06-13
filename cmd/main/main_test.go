package main

import (
	"dex/internal/container"
	"dex/internal/prompt"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	containermocks "dex/internal/container/mocks"

	promptmocks "dex/internal/prompt/mocks"
)

func TestRun(t *testing.T) {
	t.Run("FailsIfErrorWhenGettingContainers", func(t *testing.T) {
		expectedError := errors.New("Mock Error")
		mockContainerService := containermocks.NewContainerService(t)
		mockContainerService.On("GetAll").Return([]container.Container{}, expectedError)

		mockPromptService := promptmocks.NewPromptService(t)

		actualError := run(mockContainerService, mockPromptService)

		assert.Equal(t, expectedError, actualError)
	})

	t.Run("FailsSilentlyIfErrorWhenDisplayingContainerSelectPrompt", func(t *testing.T) {
		mockContainerService := containermocks.NewContainerService(t)
		mockContainerService.On("GetAll").Return(containersData(), nil)

		mockPromptService := promptmocks.NewPromptService(t)
		mockPromptService.On("DisplaySelect", mock.Anything, promptOptionsData()).
			Return(prompt.PromptOption{}, errors.New("Mock Error"))
		mockPromptService.AssertNotCalled(t, "DisplayPrompt", mock.Anything)

		actualError := run(mockContainerService, mockPromptService)

		assert.Equal(t, nil, actualError)
	})

	t.Run("FailsSilentlyIfErrorWhenDisplayingCommandPrompt", func(t *testing.T) {
		mockContainerService := containermocks.NewContainerService(t)
		mockContainerService.On("GetAll").Return(containersData(), nil)
		mockContainerService.AssertNotCalled(t, "RunCommand", mock.Anything, mock.Anything)

		mockPromptService := promptmocks.NewPromptService(t)
		mockPromptService.On("DisplaySelect", mock.Anything, promptOptionsData()).
			Return(prompt.PromptOption{}, nil)
		mockPromptService.On("DisplayPrompt", mock.Anything).
			Return("", errors.New("Mock Error"))

		actualError := run(mockContainerService, mockPromptService)

		assert.Equal(t, nil, actualError)
	})

	t.Run("NonEmptyCommandRunsSuccessfully", func(t *testing.T) {
		expectedCommand := "ls -la"
		promptOptions := promptOptionsData()
		selectedPromptOption := promptOptions[0]

		mockContainerService := containermocks.NewContainerService(t)
		mockContainerService.On("GetAll").Return(containersData(), nil)
		mockContainerService.On("RunCommand", selectedPromptOption.Value, expectedCommand).
			Return(nil)

		mockPromptService := promptmocks.NewPromptService(t)
		mockPromptService.On("DisplaySelect", mock.Anything, promptOptions).
			Return(selectedPromptOption, nil)
		mockPromptService.On("DisplayPrompt", mock.Anything).
			Return(expectedCommand, nil)

		actualError := run(mockContainerService, mockPromptService)

		assert.Equal(t, nil, actualError)
	})

	t.Run("EmptyCommandDefaultsToBash", func(t *testing.T) {
		expectedCommand := "bash"
		promptOptions := promptOptionsData()
		selectedPromptOption := promptOptions[0]

		mockContainerService := containermocks.NewContainerService(t)
		mockContainerService.On("GetAll").Return(containersData(), nil)
		mockContainerService.On("RunCommand", selectedPromptOption.Value, expectedCommand).
			Return(nil)

		mockPromptService := promptmocks.NewPromptService(t)
		mockPromptService.On("DisplaySelect", mock.Anything, promptOptions).
			Return(selectedPromptOption, nil)
		mockPromptService.On("DisplayPrompt", mock.Anything).
			Return("", nil)

		actualError := run(mockContainerService, mockPromptService)

		assert.Equal(t, nil, actualError)
	})

	t.Run("EmptyCommandAndFailedBashDefaultsToSh", func(t *testing.T) {
		expectedCommand := "sh"
		promptOptions := promptOptionsData()
		selectedPromptOption := promptOptions[0]

		mockContainerService := containermocks.NewContainerService(t)
		mockContainerService.On("GetAll").Return(containersData(), nil)
		mockContainerService.On("RunCommand", selectedPromptOption.Value, "bash").
			Return(errors.New("Mock Error")).Once()
		mockContainerService.On("RunCommand", selectedPromptOption.Value, expectedCommand).
			Return(nil).Once()

		mockPromptService := promptmocks.NewPromptService(t)
		mockPromptService.On("DisplaySelect", mock.Anything, promptOptions).
			Return(selectedPromptOption, nil)
		mockPromptService.On("DisplayPrompt", mock.Anything).
			Return("", nil)

		actualError := run(mockContainerService, mockPromptService)

		assert.Equal(t, nil, actualError)
	})
}

func promptOptionsData() []prompt.PromptOption {
	return []prompt.PromptOption{
		{
			Label: "[containerName1] (containerImage1) aa64bf12226b",
			Value: "aa64bf12226bc4dff14310a8fec7d5c5f0439ed2e69b3b590c413220650c0174",
		},
		{
			Label: "[containerName2, containerName3] (containerImage2) ab64bf12226b",
			Value: "ab64bf12226bc4dff14310a8fec7d5c5f0439ed2e69b3b590c413220650c0174",
		},
	}
}

func containersData() []container.Container {
	return []container.Container{
		{
			ID:    "aa64bf12226bc4dff14310a8fec7d5c5f0439ed2e69b3b590c413220650c0174",
			Names: []string{"containerName1"},
			Image: "containerImage1",
		},
		{
			ID:    "ab64bf12226bc4dff14310a8fec7d5c5f0439ed2e69b3b590c413220650c0174",
			Names: []string{"containerName2", "containerName3"},
			Image: "containerImage2",
		},
	}
}
