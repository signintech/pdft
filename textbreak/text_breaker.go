package textbreak

type TextBreaker interface {
	BreakTextToToken(text string) ([]string, error)
}
