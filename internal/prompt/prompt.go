package prompt

type Option struct {
	Label string
	Value string
}

//go:generate mockery --name Service
type Service interface {
	DisplaySelect(label string, options []Option) (Option, error)
	DisplayPrompt(label string) (string, error)
}
