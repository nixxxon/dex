package prompt

// Option for the select prompt
type Option struct {
	Label string
	Value string
}

// Service interface for displaying prompts
//
//go:generate mockery --name Service
type Service interface {
	DisplaySelect(label string, options []Option) (Option, error)
	DisplayPrompt(label string) (string, error)
}
