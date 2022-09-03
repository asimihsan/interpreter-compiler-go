package repl

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"monkey-book/lexer"
	"monkey-book/token"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	for {
		_, _ = fmt.Fprintf(out, PROMPT)

		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()

		l := lexer.New(strings.NewReader(line))
		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			_, _ = fmt.Fprintf(out, "%+v\n", tok)
		}
	}
}