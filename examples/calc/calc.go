//
//	calc.go implements a simple calculator using the iNamik lexer and parser api.
//
//	Input is read from STDIN
//
//	The input expression is matched against the following pattern:
//
//	input_exp:
//	( id '=' )? general_exp
//	general_exp:
//		operand ( operator operand )?
//	operand:
//		number | id | '(' general_exp ')'
//	operator:
//		'+' | '-' | '*' | '/'
//	number:
//		digit+ ( '.' digit+ )?
//	digit:
//		['0'..'9']
//	id:
//		alpha ( alpa | digit )*
//	alpha:
//		['a'..'z'] | ['A'..'Z']
//
//	Precedence is as expected, with '*' and '/' have higher precedence
//	than '+' and '-', as follows:
//
//	1 + 2 * 3 - 4 / 5  ==  1 + (2 * 3) - (4 / 5)
//

package main

import (
	"os"
	"bytes"
	"strings"
	"strconv"
	"bufio"
	"fmt"
	"github.com/iNamik/go_lexer"
	"github.com/iNamik/go_parser"
)

// We define our lexer tokens starting from the pre-defined EOF token
const (
	T_EOF   lexer.TokenType = lexer.TokenTypeEOF
	T_NIL                   = lexer.TokenTypeEOF + iota
	T_ID
	T_NUMBER
	T_PLUS
	T_MINUS
	T_MULTIPLY
	T_DIVIDE
	T_EQUALS
	T_OPEN_PAREN
	T_CLOSE_PAREN
)

// To store variables
var vars = map[string] float64 {}

// Single-character tokens
var singleChars  = []byte            { '+'   , '-'    , '*'       , '/'     , '='     , '('         , ')'           }
var singleTokens = []lexer.TokenType { T_PLUS, T_MINUS, T_MULTIPLY, T_DIVIDE, T_EQUALS, T_OPEN_PAREN, T_CLOSE_PAREN }

// Multi-character tokens
var rangeWhitespace = []byte { ' ', '\t' }
var rangeDigits     = lexer.RangeToBytes("0-9")
var rangeAlpha      = lexer.RangeToBytes("a-zA-Z")
var rangeAlphaNum   = lexer.RangeToBytes("0-9a-zA-Z")

// main
func main() {

	// Create a buffered reader from STDIN
	stdin := bufio.NewReader(os.Stdin)

	for {
		// Read a line of input
		input, _, err := stdin.ReadLine()

		// Error? we're done
		if nil != err { break }

		// Anything to process?
		if len(input) > 0 {
			// Create a new lexer to turn the input text into tokens
			l := lexer.NewLexer(lex, strings.NewReader(string(input)), len(input), 2)

			// Create a new parser that feeds off the lexer and generates expression values
			p := parser.NewParser(parse, l, 2)

			// Loop over parser emits
			for i := p.Next() ; nil != i ; i = p.Next() {
				fmt.Printf("%v\n", i)
			}
		}
	}
}

// lex is the starting (and only) StateFn for lexing the input into tokens
func lex(l lexer.Lexer) lexer.StateFn {

	// EOF
	if l.MatchEOF() {
		l.EmitEOF()
		return nil // We're done here
	}

	// Single-char token?
	if i := bytes.IndexRune(singleChars, l.PeekRune(0)) ; i >= 0 {
		l.NextRune()
		l.EmitToken(singleTokens[i])
		return lex
	}

	switch {

		// Skip whitespace
		case l.MatchOneOrMore(rangeWhitespace) :
			l.IgnoreToken()

		// Number
		case l.MatchOneOrMore(rangeDigits) :
			if l.PeekRune(0) == '.' {
				l.NextRune() // skip '.'
				if ! l.MatchOneOrMore(rangeDigits) {
					printError(l.Column(), "Illegal number format - Missing digits after '.'")
					l.IgnoreToken()
					break
				}
			}
			l.EmitTokenWithBytes(T_NUMBER)

		// ID
		case l.MatchOne(rangeAlpha) && l.MatchNoneOrMore(rangeAlphaNum):
			l.EmitTokenWithBytes(T_ID)

		// Unknown
		default :
			l.NextRune()
			printError(l.Column(), "Unknown Character")
			l.IgnoreToken()
	}

	// See you again soon!
	return lex
}

// parse tries to execute a general expression from the lexed tokens.
// Returns nil - We only take one pass at the input string
func parse(p parser.Parser) parser.StateFn {

	if p.PeekTokenType(0) != T_EOF {
		// Assignment ( id = general_expression )
		if p.PeekTokenType(0) == T_ID && p.PeekTokenType(1) == T_EQUALS {
			tId := p.NextToken()

			p.SkipToken() // skip '='

			val, ok := pGeneralExpression(p)

			if ok {
				id := string(tId.Bytes())
				vars[id] = val

			}
		// General expression
		} else {
			val, ok := pGeneralExpression(p)

			if ok {
				p.Emit(val)
			}
		}
	}

	p.ClearTokens()
	p.Emit(nil) // We're done - One pass only

	return nil
}

// pGeneralExpression is the starting point for parsing a General Expression.
// It is basically a pass-through to pAdditiveExpression, but it feels cleaner
func pGeneralExpression(p parser.Parser) (f float64, ok bool) { return pAdditiveExpression(p) }

// pAdditiveExpression parses [ expression ( ( '+' | '-' ) expression )? ]
func pAdditiveExpression(p parser.Parser) (f float64, ok bool) {

	f, ok = pMultiplicitiveExpression(p)

	if ok {
		t := p.NextToken()
		switch t.Type() {

			// Add (+)
			case T_PLUS :
				r, ok := pAdditiveExpression(p)
				if ok {
					f += r
				}

			// Subtract (-)
			case T_MINUS :
				r, ok := pAdditiveExpression(p)
				if ok {
					f -= r
				}

			// EOF
			case T_EOF :
				p.BackupToken()
				ok = true

			// Unknown
			default :
				printError(t.Column(), "Expecting operator")
				ok = false
		}
	}

	return
}

// pMultiplicitiveExpression parses [ expression ( ( '*' | '/' ) expression )? ]
func pMultiplicitiveExpression(p parser.Parser) (f float64, ok bool) {

	f, ok = pOperand(p)

	if ok {
		t := p.NextToken()
		switch t.Type() {

			// Multiply (*)
			case T_MULTIPLY :
				r, ok := pMultiplicitiveExpression(p)
				if ok {
					f *= r
				}

			// Divide (/)
			case T_DIVIDE :
				r, ok := pMultiplicitiveExpression(p)
				if ok {
					f /= r
				}

			// Unknown - Send it back upstream
			default :
				p.BackupToken()
				ok = true
		}
	}

	return
}

// pOperand parses [ id | number | '(' expression ')' ]
func pOperand (p parser.Parser) (f float64, ok bool) {

	var err os.Error

	m := p.Marker()
	t := p.NextToken()

	switch t.Type() {

		// ID
		case T_ID :
			var id = string( t.Bytes() )
			f, ok = vars[ id ]
			if !ok {
				printError(t.Column(), fmt.Sprint("id '",id,"' not defined"))
				f = 0.0
			}

		// Number
		case T_NUMBER :
			f, err = strconv.Atof64( string( t.Bytes() ) )
			ok = nil == err
			if !ok {
				printError(t.Column(), fmt.Sprint("Error reading number: ",err.String()))
				f = 0.0
			}

		// '(' Expresson ')'
		case T_OPEN_PAREN :
			f, ok = pGeneralExpression(p)
			if ok {
				t2 := p.NextToken()
				if t2.Type() != T_CLOSE_PAREN {
					printError(t.Column(), "Unbalanced Paren")
					ok = false
					f = 0.0
				}
			}

		// EOF
		case T_EOF:
			printError(t.Column(), "Unexpected EOF - Expecting operand")
			ok = false
			f = 0.0

		// Unknown
		default:
			printError(t.Column(), "Expecting operand")
			ok = false
			f = 0.0
	}

	if !ok {
		p.Reset(m)
	}

	return
}

// printError prints an error msg pointing to the specified column of the input.
func printError(col int, msg string) {
	fmt.Print(strings.Repeat(" ", col-1), "^ ", msg, "\n")
}

