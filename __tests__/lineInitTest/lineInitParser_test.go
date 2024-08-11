package line_init_parser_test

import (
	"misc/nintasm/parser"
	"testing"
)

func TestInitialLineParser(t *testing.T) {
	parser := parser.NewInitialLineParser()

	t.Run("Positive", func(t *testing.T) {
		result, _ := parser.Process("\tlda #1 ;Something")
		if result != " lda #1" {
			t.Error()
		}
	})
}
