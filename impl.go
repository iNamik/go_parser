package parser

import "github.com/iNamik/lexer.go"

/**
 * Parser::Nextn - Returns the next emit from the parser.
 */
func (p *parser) Next() interface{} {
	for {
		select {
			case i := <- p.chn :
				return i
			default:
				p.state = p.state(p)
		}
	}
	panic("not reached")
}

/**
 * Parser::PeekTokenType
 */
func (p *parser) PeekTokenType(n int) lexer.TokenType {
	return p.PeekToken(n).Type();
}

/**
 * Parser::PeekToken
 */
func (p *parser) PeekToken(n int) *lexer.Token {
	ok := p.ensureTokenLen( p.pos + n + 1 ) // Correct for 0-based 'n'

	if !ok {
		if nil == p.eofToken { panic("illegal state: eofToken is nil") }
		return p.eofToken
	}

	i := p.tokens.Peek( p.pos + n )

	return i.(*lexer.Token)
}

/**
 * Parser::NextToken
 */
func (p *parser) NextToken() *lexer.Token {
	ok := p.ensureTokenLen( p.pos + 1 )

	if !ok {
		if nil == p.eofToken { panic("illegal state: eofToken is nil") }
		return p.eofToken
	}

	i := p.tokens.Peek( p.pos ) // 0-based

	p.pos++

	return i.(*lexer.Token)
}

/**
 * Parser::SkipToken
 */
func (p *parser) SkipToken() {
	ok := p.ensureTokenLen( p.pos + 1 )

	if ok {
		p.pos++
	}
}

/**
 * Parser::SkipTokens
 */
func (p *parser) SkipTokens(n int) {
	ok := p.ensureTokenLen( p.pos + n + 1 )

	if ok {
		p.pos += n
	}
}

/**
 * Parser::BackupToken
 */
func (p *parser) BackupToken() {
	p.BackupTokens(1)
}

/**
 * Parser::BackupTokens
 */
func (p *parser) BackupTokens(n int) {
	if n > p.pos {
		panic("Underflow Exception")
	}
	p.pos -= n
}

/**
 * Parser::ClearTokens
 */
func (p *parser) ClearTokens() {
	for ; p.pos > 0 ; p.pos-- {
		p.tokens.Remove()
	}
}

/**
 * Parser::Emit
 */
func (p *parser) Emit(i interface{}) {
	p.chn <- i
}

/**
 * Parser::EOF
 */
func (p *parser) EOF() bool {
	return p.eof
}

/**
 * Parser::Marker
 */
func (p * parser) Marker() *Marker {
	return &Marker{sequence: p.sequence, pos: p.pos}
}

/**
 * Parser::Reset
 */
func (p *parser) Reset(m *Marker) {
	if (m.sequence != p.sequence || m.pos < 0 || m.pos >= p.tokens.Len()) {
		panic("Invalid marker")
	}
	p.pos = m.pos
}
