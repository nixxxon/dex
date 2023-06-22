package prompt

import "github.com/manifoldco/promptui"

// PromptuiService for displaying prompts
type PromptuiService struct{}

// DisplaySelect displays a select prompt
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

// DisplayPrompt displays a normal prompt
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

// NewPromptuiService returns a new PromptuiService pointer
func NewPromptuiService() Service {
	return &PromptuiService{}
}
