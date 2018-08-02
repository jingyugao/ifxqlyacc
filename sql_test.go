package ifxqlyacc

import (
	"testing"
	"strings"
)

func TestYyParser(t *testing.T) {
	//s := NewScanner(strings.NewReader("SELECT value as a from myseries WHERE a = 'b"))
	tokenizer := &Tokenizer{
		query:Query{},
		//scanner:NewScanner(strings.NewReader("select *  From b where a = 1 order by time")),
	}
	for i:=0;i<100000;i++{
		tokenizer.scanner = NewScanner(strings.NewReader("select *  From b where a = 1 order by time"))
		yyParse(tokenizer)
	}
	println("lex cost " ,tokenizer.dur)
}

