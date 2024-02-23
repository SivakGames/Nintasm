package macroTable

import "regexp"

type Replacer struct {
	ReplaceRegex  *regexp.Regexp
	ReplaceString string
}

//Create a replacer.
//substitutionID is the desired ID sans the preceding backslash
//replaceString is what will actually be replaced when instances of substitutionID are found
func NewReplacer(substitutionID string, replaceString string) Replacer {
	return Replacer{
		ReplaceRegex:  regexp.MustCompile(`\\` + substitutionID + `\b`),
		ReplaceString: replaceString,
	}
}
