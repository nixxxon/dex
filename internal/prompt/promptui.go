package prompt

import "github.com/manifoldco/promptui"

type PromptuiService struct{}

func (p *PromptuiService) DisplaySelect(
	label string,
	options []Option,
) (Option, error) {
	labels := []string{}
	for _, option := range options {
		labels = append(labels, option.Label)
	}

	prompt := promptui.Select{
		Label: label,
		Items: labels,
	}
	i, _, err := prompt.Run()
	if err != nil {
		return Option{}, err
	}

	return options[i], nil
}

func (p *PromptuiService) DisplayPrompt(label string) (string, error) {
	prompt := promptui.Prompt{
		Label: label,
	}
	result, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return result, nil
}

func NewPromptuiService() Service {
	return &PromptuiService{}
}
