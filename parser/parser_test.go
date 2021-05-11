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

		p = &Parser{}
		typ = p.GetSQLType("UPSERT INTO foo VALUES (1,2,3)")
		So(typ, ShouldEqual, UNSUPPORTED)
	})
}

func TestParserSelect(t *testing.T) {
	var p *Parser
	Convey("SELECT SQL", t, func() {
		p = &Parser{}
		ast, err := p.ParseSelect("SELECT ab,b, c FROM foo WHERE id < 3")
		So(err, ShouldBeNil)
		So(ast.Projects, ShouldResemble, []string{"ab", "b", "c"})
		So(ast.Table, ShouldEqual, "foo")

		p = &Parser{}
		ast, err = p.ParseSelect("SELECT ab,b, c FROM foo LIMIT 3")
		So(err, ShouldBeNil)
		So(ast.Projects, ShouldResemble, []string{"ab", "b", "c"})
		So(ast.Table, ShouldEqual, "foo")
		So(ast.Limit, ShouldEqual, 3)

		p = &Parser{}
		ast, err = p.ParseSelect("SELECT ab,b,c FROM foo WHERE id < 3 AND ab > 10 LIMIT 11")
		So(err, ShouldBeNil)
		So(ast.Projects, ShouldResemble, []string{"ab", "b", "c"})
		So(ast.Where, ShouldResemble, []string{"id", "<", "3", "AND", "ab", ">", "10"})
		So(ast.Table, ShouldEqual, "foo")
		So(ast.Limit, ShouldEqual, 11)

		p = &Parser{}
		ast, err = p.ParseSelect("SELECT 1")
		So(err, ShouldBeNil)
		So(ast.Projects, ShouldResemble, []string{"1"})
		So(ast.Table, ShouldEqual, "")
	})
}

func TestParserInsert(t *testing.T) {
	var p *Parser
	Convey("INSERT SQL with Column names", t, func() {
		p = &Parser{}
		ast, err := p.ParseInsert("INSERT INTO table_name(column1, column2) VALUES (value1, value2)")
		So(err, ShouldBeNil)
		So(ast.Table, ShouldEqual, "table_name")
		So(ast.Columns, ShouldResemble, []string{"column1", "column2"})
		So(ast.Values, ShouldResemble, [][]string{{"value1", "value2"}})
	})

	Convey("column count miss match", t, func() {
		p = &Parser{}
		_, err := p.ParseInsert("INSERT INTO table_name(column1, column2, column3) VALUES (value1, value2)")
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldContainSubstring, "expected column count")
	})

	Convey("INSERT SQL", t, func() {
		p = &Parser{}
		ast, err := p.ParseInsert("INSERT INTO table_name VALUES (value1, value2)")
		So(err, ShouldBeNil)
		So(ast.Table, ShouldEqual, "table_name")
		So(ast.Columns, ShouldBeNil)
		So(ast.Values, ShouldResemble, [][]string{{"value1", "value2"}})
	})

	Convey("INSERT multiple rows", t, func() {
		p = &Parser{}
		ast, err := p.ParseInsert("INSERT INTO table_name VALUES (\"value1\", value2), (\"value3\", value4)")
		So(err, ShouldBeNil)
		So(ast.Table, ShouldEqual, "table_name")
		So(ast.Columns, ShouldBeNil)
		So(ast.Values, ShouldResemble, [][]string{{"\"value1\"", "value2"}, {"\"value3\"", "value4"}})
	})

	Convey("INSERT multiple rows 2", t, func() {
		p = &Parser{}
		ast, err := p.ParseInsert("INSERT INTO table (id, username, email) VALUES " +
			"(0, auxten, \"auxtenwpc@gmail.com\")," +
			"(1, hahaha, \"hahaha@gmail.com\")," +
			"(2, jijiji, \"jijiji@gmail.com\")")
		So(err, ShouldBeNil)
		So(ast.Table, ShouldEqual, "table")
		So(ast.Columns, ShouldResemble, []string{"id", "username", "email"})
		So(ast.Values, ShouldResemble, [][]string{{"0", "auxten", "\"auxtenwpc@gmail.com\""},
			{"1", "hahaha", "\"hahaha@gmail.com\""}, {"2", "jijiji", "\"jijiji@gmail.com\""}})
	})

	Convey("column count miss match 2", t, func() {
		p = &Parser{}
		_, err := p.ParseInsert("INSERT INTO table_name VALUES (value1, value2), (value3, value4, value5)")
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldContainSubstring, "expected column count")
	})
}
