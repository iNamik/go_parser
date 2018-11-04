package parser

import "github.com/iNamik/go_lexer"
import "github.com/iNamik/go_container/queue"

// StateFn represents the state of the parser as a function that returns the next state.
type StateFn func(Parser) StateFn

// Marker stores the state of the parser to allow rewinding
type Marker struct {
	sequence int
	pos      int
}

// Parser helps you process lexer tokens
type Parser interface {

	// Line returns the line number of the next token (aka PeekTotenType(0).Line() )
	Line() int

	// Column returns the column number of the next token (aka PeekTokenType(0).column() )
	Column() int

	// PeekTokenType allows you to look ahead at tokens without consuming them
	PeekTokenType(int) lexer.TokenType

	// PeekToken allows you to look ahead at tokens without consuming them
	PeekToken(int) *lexer.Token

	// NextToken consumes and returns the next token
	NextToken() *lexer.Token

	// SkipToken consumes the next token without returning it
	SkipToken()

	// SkipTokens consumes the next n tokens without returning them
	SkipTokens(int)

	// BackupToken un-consumes the last token
	BackupToken()

	// BackupTokens un-consumes the last n tokens
	BackupTokens(int)

	// ClearTokens clears all consumed tokens
	ClearTokens()

	// Emit emits an object, consuming matched tokens
	Emit(interface{})

	EOF() bool

	// Next retrieves the next emitted item
	Next() interface{}

	// Marker returns a marker that you can use to reset the parser state later
	Marker() *Marker

	// Reset resets the lexer state to the specified marker
	Reset(*Marker)
}

// New returns a new Parser object
func New(startState StateFn, lex lexer.Lexer, channelCap int) Parser {
	p := &parser{
		lex:      lex,
		tokens:   queue.New(4), // 4 is just a nice number that seems appropriate
		pos:      0,
		sequence: 0,
		eofToken: nil,
		eof:      false,
		state:    startState,
		chn:      make(chan interface{}, channelCap),
	}
	return p
}
