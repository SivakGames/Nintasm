package macroTable

import "regexp"

type Replacer struct {
	ReplaceRegex  *regexp.Regexp
	ReplaceString string
}

type ReplacementList []Replacer

var ReplacementStack []ReplacementList

//Create a replacer.
//substitutionID is the desired ID
//replaceString is what will actually be replaced when instances of substitutionID are found
func NewReplacer(substitutionID string, replaceString string) Replacer {
	return Replacer{
		ReplaceRegex:  regexp.MustCompile(`\` + substitutionID + `\b`),
		ReplaceString: replaceString,
	}
}

func AppendToReplacementStack() {
	ReplacementStack = append(ReplacementStack, []Replacer{})
}
func PopFromReplacementStack() {
	ReplacementStack = ReplacementStack[:len(ReplacementStack)-1]
}
func GetReplacementListOnTopOfStack() *ReplacementList {
	return &ReplacementStack[len(ReplacementStack)-1]
}
func AddToReplacementListOnTopOfStack(replaceName string, replaceValue string) {
	current := GetReplacementListOnTopOfStack()
	*current = append(*current, NewReplacer(replaceName, replaceValue))
}
func AddNumToReplacementListOnTopOfStack(replaceValue string) {
	current := GetReplacementListOnTopOfStack()
	*current = append(*current, Replacer{
		ReplaceRegex:  regexp.MustCompile(`\\#`),
		ReplaceString: replaceValue,
	})
}
