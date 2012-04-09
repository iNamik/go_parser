package parser

import "github.com/iNamik/go_lexer"
import "github.com/iNamik/go_container/queue"

type parser struct {
	lex       lexer.Lexer
	tokens    queue.Interface
	pos       int
	sequence  int
	eof       bool
	eofToken *lexer.Token
	state     StateFn
	chn       chan interface{} // channel of objects
}

func (p *parser) ensureTokenLen(n int) bool{
	for !p.eof && p.tokens.Len() < n {
		token := p.lex.NextToken()
		if token.EOF() {
			p.eofToken = token
			p.eof = true
		}
		p.tokens.Add(token)
	}
	return p.tokens.Len() >= n
}
