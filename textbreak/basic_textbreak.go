package textbreak

type BasicTextbreak struct {
}

func (BasicTextbreak) BreakTextToToken(text string) ([]string, error) {
	var tokens []string
	var tmp = ""
	for _, ru := range text {
		if ru == ' ' {
			tokens = append(tokens, tmp)
			tokens = append(tokens, string(ru))
			tmp = ""
			continue
		}
		tmp += string(ru)
	}
	if tmp != "" {
		tokens = append(tokens, tmp)
	}

	return tokens, nil
}
