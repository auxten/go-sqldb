package parser

import (
	"fmt"
	"strconv"
	"strings"
	"text/scanner"
)

type Token int

// Common tokens
const (
	ILLEGAL = iota
	SPACE   // space, tab, \r, \n
	COMMA
)

type StatementType int

// SQL type tokens
// Start from 10 to reserve common tokens
const (
	UNSUPPORTED = iota + 10
	SELECT
	INSERT
)

// Other tokens
// Start from 20 to reserve SQL type tokens
const (
	FROM = iota + 20
	INTO
)

type Parser struct {
	s scanner.Scanner
}

type SelectTree struct {
	Projects []string
	From     string
	Where    []string
	Limit    int64
}

type InsertTree struct {
}

func (p *Parser) GetSQLType(sql string) StatementType {
	s := p.s
	s.Init(strings.NewReader(sql))
	if tok := s.Scan(); tok != scanner.EOF {
		txt := strings.ToUpper(s.TokenText())
		switch txt {
		case "SELECT":
			return SELECT
		case "INSERT":
			return INSERT
		default:
			return UNSUPPORTED
		}
	}

	return UNSUPPORTED
}

/*
ParserSelect is a simple select statement parser.
It's just a demo of SELECT statement parser skeleton.
Currently, the most complex SQL supported here is something like:

	SELECT * FROM foo WHERE id < 3 LIMIT 1;

Even SQL-92 standard is far more complex.
For a production ready SQL parser, see: https://github.com/auxten/postgresql-parser
*/
func (p *Parser) ParserSelect(sel string) (ast *SelectTree, err error) {
	ast = &SelectTree{}
	s := p.s
	s.Init(strings.NewReader(sel))
	if tok := s.Scan(); tok == scanner.EOF || strings.ToUpper(s.TokenText()) != "SELECT" {
		err = fmt.Errorf("%s is not SELECT statement", sel)
		return
	}

	ast.Projects = make([]string, 0, 4)
	for {
		if tok := s.Scan(); tok == scanner.EOF {
			if len(ast.Projects) == 0 {
				err = fmt.Errorf("%s get select projects failed", sel)
			}
			return
		} else {
			txt := s.TokenText()
			//log.Print(txt)
			if txt == "*" {
				ast.Projects = append(ast.Projects, "*")
			} else {
				if txt == "," {
					continue
				} else if strings.ToUpper(txt) == "FROM" {
					break
				} else {
					ast.Projects = append(ast.Projects, txt)
				}
			}
		}
	}

	// token FROM is scanned, try to get the table name here
	// FROM ?
	if tok := s.Scan(); tok == scanner.EOF {
		// if projects are all constant value, source table is not necessary.
		// eg.  SELECT 1;
		return
	} else {
		ast.From = s.TokenText()
	}

	// WHERE
	if tok := s.Scan(); tok == scanner.EOF {
		// WHERE is not necessary
		return
	}

	txt := s.TokenText()
	if strings.ToUpper(txt) != "WHERE" {
		err = fmt.Errorf("expect WHERE here")
		return
	}

	// token WHERE is scanned, try to get the WHERE clause.
	ast.Where = make([]string, 0, 4)
	for {
		if tok := s.Scan(); tok == scanner.EOF {
			if len(ast.Where) == 0 {
				err = fmt.Errorf("missing WHERE clause")
			}
			return
		}
		txt := s.TokenText()
		if strings.ToUpper(txt) == "LIMIT" {
			break
		}
		ast.Where = append(ast.Where, txt)
	}

	// token LIMIT is scanned, try to get the limit
	if tok := s.Scan(); tok == scanner.EOF {
		err = fmt.Errorf("expect LIMIT clause here")
		return
	}
	txt = s.TokenText()
	ast.Limit, err = strconv.ParseInt(txt, 10, 32)
	return
}

func (p *Parser) ParserInsert(insert string) (ast *InsertTree, err error) {

	return
}
