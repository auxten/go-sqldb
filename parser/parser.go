package parser

import (
	"fmt"
	"strconv"
	"strings"
	"text/scanner"
)

type StatementType string

// SQL type tokens
const (
	UNSUPPORTED = "N/A"
	SELECT      = "SELECT"
	FROM        = "FROM"
	WHERE       = "WHERE"
	LIMIT       = "LIMIT"
	INSERT      = "INSERT"
	INTO        = "INTO"
	VALUES      = "VALUES"
	ASTERISK    = "*"
)

type Parser struct {
	s scanner.Scanner
}

type SelectTree struct {
	Projects []string
	Table    string
	Where    []string
	Limit    int64
}

type InsertTree struct {
	Table   string
	Columns []string
	Values  [][]string
}

func (p *Parser) GetSQLType(sql string) StatementType {
	s := p.s
	s.Init(strings.NewReader(sql))
	s.Mode = scanner.ScanIdents | scanner.ScanFloats | scanner.ScanChars | scanner.ScanStrings | scanner.ScanRawStrings

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
ParseSelect is a simple select statement parser.
It's just a demo of SELECT statement parser skeleton.
Currently, the most complex SQL supported here is something like:

	SELECT * FROM foo WHERE id < 3 LIMIT 1;

Even SQL-92 standard is far more complex.
For a production ready SQL parser, see: https://github.com/auxten/postgresql-parser
*/
func (p *Parser) ParseSelect(sel string) (ast *SelectTree, err error) {
	ast = &SelectTree{}
	s := p.s
	s.Init(strings.NewReader(sel))
	s.Mode = scanner.ScanIdents | scanner.ScanFloats | scanner.ScanChars | scanner.ScanStrings | scanner.ScanRawStrings

	if tok := s.Scan(); tok == scanner.EOF || strings.ToUpper(s.TokenText()) != SELECT {
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
			if txt == ASTERISK {
				ast.Projects = append(ast.Projects, ASTERISK)
			} else {
				if txt == "," {
					continue
				} else if strings.ToUpper(txt) == FROM {
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
		ast.Table = s.TokenText()
	}

	// WHERE
	if tok := s.Scan(); tok == scanner.EOF {
		// WHERE is not necessary
		return
	}

	txt := s.TokenText()
	if strings.ToUpper(txt) == WHERE {
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
			if strings.ToUpper(txt) == LIMIT {
				break
			}
			ast.Where = append(ast.Where, txt)
		}
	} else if strings.ToUpper(txt) != LIMIT {
		err = fmt.Errorf("expect WHERE or LIMIT here")
		return
	}

	// token LIMIT is scanned, try to get the limit
	if tok := s.Scan(); tok == scanner.EOF {
		err = fmt.Errorf("expect LIMIT clause here")
		return
	}
	txt = s.TokenText()
	ast.Limit, err = strconv.ParseInt(txt, 10, 64)
	return
}

/*
ParseInsert can parse a simple INSERT statement, eg.
 	INSERT INTO table_name VALUES (value1, value2, …)
	or
	INSERT INTO table_name(column1, column2, …) VALUES (value1, value2, …)
*/
func (p *Parser) ParseInsert(insert string) (ast *InsertTree, err error) {
	ast = &InsertTree{}
	s := p.s
	s.Init(strings.NewReader(insert))
	s.Mode = scanner.ScanIdents | scanner.ScanFloats | scanner.ScanChars | scanner.ScanStrings | scanner.ScanRawStrings

	if tok := s.Scan(); tok == scanner.EOF || strings.ToUpper(s.TokenText()) != INSERT {
		err = fmt.Errorf("%s is not INSERT statement", insert)
		return
	}

	if tok := s.Scan(); tok == scanner.EOF || strings.ToUpper(s.TokenText()) != INTO {
		err = fmt.Errorf("%s expect INTO after INSERT", insert)
		return
	}

	// Table name
	if tok := s.Scan(); tok == scanner.EOF {
		err = fmt.Errorf("%s expect table after INSERT INTO", insert)
		return
	} else {
		ast.Table = s.TokenText()
	}

	var columnCnt int
	// try get colNames
	if tok := s.Scan(); tok == scanner.EOF {
		err = fmt.Errorf("%s expect VALUES or (colNames)", insert)
		return
	} else {
		txt := strings.ToUpper(s.TokenText())
		if txt == "(" {
			ast.Columns = make([]string, 0, 4)
			for {
				if tok := s.Scan(); tok == scanner.EOF {
					if len(ast.Columns) == 0 {
						err = fmt.Errorf("%s get Columns failed", insert)
					}
					return
				} else {
					txt := s.TokenText()
					//log.Print(txt)
					if txt == "," {
						continue
					} else if txt == ")" {
						continue
					} else if strings.ToUpper(txt) == VALUES {
						break
					} else {
						ast.Columns = append(ast.Columns, txt)
					}
				}
			}
		} else if txt != VALUES {
			err = fmt.Errorf("%s expect VALUES or '(' here", insert)
			return
		}
	}
	columnCnt = len(ast.Columns)

	// VALUES has been scanned try to get (value1, value2), (value3, value4)
	ast.Values = make([][]string, 0, 4)
rowLoop:
	for {
		if tok := s.Scan(); tok == scanner.EOF {
			break rowLoop
		} else {
			txt := s.TokenText()
			if txt == "," {
				// next row
				continue
			}
			if txt == "(" {
				var row []string
				if columnCnt != 0 {
					row = make([]string, 0, columnCnt)
				} else {
					row = make([]string, 0, 4)
				}
				for {
					if tok := s.Scan(); tok == scanner.EOF {
						break rowLoop
					} else {
						txt := s.TokenText()
						//log.Print(txt)
						if txt == "," {
							continue
						} else if txt == ")" {
							ast.Values = append(ast.Values, row)
							break
						} else {
							row = append(row, txt)
						}
					}
				}
			}
		}
	}

	// Check if column count identical
	for _, row := range ast.Values {
		if columnCnt == 0 {
			columnCnt = len(row)
		} else {
			if columnCnt != len(row) {
				err = fmt.Errorf(
					"%s expected column count is %d, got %d, %v",
					insert, columnCnt, len(row), row,
				)
				return
			}
		}
	}

	return
}
