package parser

import (
	"strconv"
	"time"

	"github.com/Zac-Garby/booleang/ast"
	"github.com/Zac-Garby/booleang/token"
)

func (p *Parser) next() {
	p.cur = p.peek
	p.peek = p.lex()

	if p.peek.Type == token.Illegal {
		p.err(
			"illegal token found: `%s`",
			p.peek.Range, p.peek.Literal,
		)
	}
}

func (p *Parser) curIs(ts ...token.Type) bool {
	for _, t := range ts {
		if p.cur.Type == t {
			return true
		}
	}

	return false
}

func (p *Parser) peekIs(ts ...token.Type) bool {
	for _, t := range ts {
		if p.peek.Type == t {
			return true
		}
	}

	return false
}

func (p *Parser) expect(t token.Type) bool {
	if p.peekIs(t) {
		p.next()
		return true
	}

	p.peekErr(t)
	return false
}

func (p *Parser) parseInt() (int64, error) {
	return strconv.ParseInt(p.cur.Literal, 10, 64)
}

func (p *Parser) parseDuration() *time.Duration {
	val, err := p.parseInt()
	if err != nil {
		p.Errors = append(p.Errors, err)
		return nil
	}

	if !p.expect(token.Ident) {
		return nil
	}

	var dur time.Duration

	switch p.cur.Literal {
	case "ns":
		dur = time.Nanosecond
	case "ms":
		dur = time.Millisecond
	case "s":
		dur = time.Second
	case "m":
		dur = time.Minute
	case "h":
		dur = time.Hour
	default:
		p.curErr("expexcted ns, ms, s, m, or h to complete the clock duration. got %s", p.cur.Literal)
		return nil
	}

	dur *= time.Duration(val)

	return &dur
}

func (p *Parser) parseExprs(end token.Type) []ast.Expression {
	var exprs []ast.Expression

	if p.peekIs(end) {
		p.next()
		return exprs
	}

	p.next()
	exprs = append(exprs, p.parseExpression())

	for p.peekIs(token.Comma) {
		p.next()

		if p.peekIs(end) {
			p.next()
			return exprs
		}

		p.next()
		exprs = append(exprs, p.parseExpression())
	}

	if !p.expect(end) {
		return nil
	}

	return exprs
}

func (p *Parser) parseIdents(end token.Type) []string {
	var idents []string

	if p.peekIs(end) {
		p.next()
		return exprs
	}

	if !p.expect(token.Ident) {
		return idents
	}
	idents = append(exprs, p.cur.Literal)

	for p.peekIs(token.Comma) {
		p.next()

		if p.peekIs(end) {
			p.next()
			return exprs
		}

		if !p.expect(token.Ident) {
			return idents
		}
		idents = append(exprs, p.cur.Literal)
	}

	if !p.expect(end) {
		return nil
	}

	return idents
}

func (p *Parser) parseStatements() []ast.Statement {
	var stmts []ast.Statement
	p.next()

	for !p.curIs(token.RightBrace) && !p.curIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			stmts = append(stmts, stmt)
		}
		p.next()
	}

	return stmts
}
