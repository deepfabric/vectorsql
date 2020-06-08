package parser

import "github.com/deepfabric/vectorsql/pkg/sql/tree"

type Parser struct {
	lexer      lexer
	scanner    scanner
	parserImpl sqlParserImpl
}

func Parse(sql string) (*tree.Select, error) {
	var p Parser

	p.scanner.init(sql)
	_, tokens, _ := p.scanOneStmt()
	return p.parse(sql, tokens)
}

// parse parses a statement from the given scanned tokens.
func (p *Parser) parse(sql string, tokens []sqlSymType) (*tree.Select, error) {
	p.lexer.init(sql, tokens)
	defer p.lexer.cleanup()
	if p.parserImpl.Parse(&p.lexer) != 0 {
		if p.lexer.lastError == nil {
			// This should never happen -- there should be an error object
			// every time Parse() returns nonzero. We're just playing safe
			// here.
			p.lexer.Error("syntax error")
		}
		return nil, p.lexer.lastError
	}
	return p.lexer.stmt, nil
}

func (p *Parser) scanOneStmt() (string, []sqlSymType, bool) {
	var lval sqlSymType
	var tokens []sqlSymType

	// Scan the first token.
	for {
		p.scanner.scan(&lval)
		if lval.id == 0 {
			return "", nil, true
		}
		if lval.id != ';' {
			break
		}
	}

	startPos := lval.pos
	// We make the resulting token positions match the returned string.
	lval.pos = 0
	tokens = append(tokens, lval)
	for {
		if lval.id == ERROR {
			return p.scanner.in[startPos:], tokens, true
		}
		posBeforeScan := p.scanner.pos
		p.scanner.scan(&lval)
		if lval.id == 0 || lval.id == ';' {
			return p.scanner.in[startPos:posBeforeScan], tokens, (lval.id == 0)
		}
		lval.pos -= startPos
		tokens = append(tokens, lval)
	}
}
