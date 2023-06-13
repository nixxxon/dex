package prompt

type PromptOption struct {
	Label string
	Value string
}

//go:generate mockery --name PromptService
type PromptService interface {
	DisplaySelect(label string, options []PromptOption) (PromptOption, error)
	DisplayPrompt(label string) (string, error)
}
