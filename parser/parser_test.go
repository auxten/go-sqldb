package parser

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestParser(t *testing.T) {
	var p *Parser
	Convey("SQL type", t, func() {
		p = &Parser{}
		typ := p.GetSQLType("SELECT * FROM foo WHERE id < 3")
		So(typ, ShouldEqual, SELECT)

		p = &Parser{}
		typ = p.GetSQLType("INSERT INTO foo VALUES (1,2,3)")
		So(typ, ShouldEqual, INSERT)
	})
}

func TestParserSelect(t *testing.T) {
	var p *Parser

	Convey("SELECT SQL", t, func() {
		p = &Parser{}
		ast, err := p.ParserSelect("SELECT ab,b, c FROM foo WHERE id < 3")
		So(err, ShouldBeNil)
		So(ast.Projects, ShouldResemble, []string{"ab", "b", "c"})
		So(ast.From, ShouldEqual, "foo")

		p = &Parser{}
		ast, err = p.ParserSelect("SELECT ab,b,c FROM foo WHERE id < 3 AND ab > 10 LIMIT 11")
		So(err, ShouldBeNil)
		So(ast.Projects, ShouldResemble, []string{"ab", "b", "c"})
		So(ast.Where, ShouldResemble, []string{"id", "<", "3", "AND", "ab", ">", "10"})
		So(ast.From, ShouldEqual, "foo")
		So(ast.Limit, ShouldEqual, 11)

		p = &Parser{}
		ast, err = p.ParserSelect("SELECT 1")
		So(err, ShouldBeNil)
		So(ast.Projects, ShouldResemble, []string{"1"})
		So(ast.From, ShouldEqual, "")
	})
}
