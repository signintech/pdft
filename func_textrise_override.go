package pdft

//FuncTextriseOverride override text rise
type FuncTextriseOverride func(
	leftRune rune,
	rightRune rune,
	fontsize int,
	allText string,
	currTextIndex int,
) float32

//FuncKernOverride  return your custome pair value
type FuncKernOverride func(
	leftRune rune,
	rightRune rune,
	leftPair uint,
	rightPair uint,
	pairVal int16,
) int16
