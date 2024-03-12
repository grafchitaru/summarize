package ai

type AI interface {
	Send(text string, prompt string) (string, error)
	GetCountTokens(text string) int
}
