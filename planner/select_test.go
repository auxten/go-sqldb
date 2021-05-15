package planner

import (
	"fmt"
	"os"
	"testing"

	"github.com/auxten/go-sqldb/db"
	"github.com/auxten/go-sqldb/node"
	"github.com/auxten/go-sqldb/parser"
	"github.com/auxten/go-sqldb/utils"
	. "github.com/smartystreets/goconvey/convey"
)

func TestPlannerSelect(t *testing.T) {
	Convey("Volcano model select implementation", t, func() {
		const testFile = "test_plan_insert.db"
		const InsertCnt = 512
		defer func() {
			_ = os.Remove(testFile)
		}()
		table, err := db.Open(testFile)
		So(err, ShouldBeNil)
		So(table.Pager.Pages[0].LeafNode.CommonHeader.IsInternal, ShouldBeFalse)
		So(err, ShouldBeNil)
		So(table, ShouldNotBeNil)

		for i := uint32(0); i < InsertCnt; i++ {
			row := &node.Row{
				Id:       i,
				Username: [32]byte{'a', 'u', 'x', 't', 'e', 'n', byte('a' + i)},
				Email:    [256]byte{'a', 'u', 'x', 't', 'e', 'n', '@', byte('a' + i)},
			}
			err = table.Insert(row)
			//fmt.Println(i, string(row.Username[:]), string(row.Email[:]))
			So(err, ShouldBeNil)
		}

		p := &parser.Parser{}
		ast, err := p.ParseSelect("SELECT id, username, email FROM table WHERE id > 5 AND id < 7 LIMIT 3")
		So(err, ShouldBeNil)

		plan := NewPlan(table)
		resultPipe, err := plan.SelectPrepare(ast)
		So(err, ShouldBeNil)
		var i int
		for row := range resultPipe {
			i++
			fmt.Println(row.Id,
				string(row.Username[:utils.Length(row.Username[:])]),
				string(row.Email[:utils.Length(row.Email[:])]))
		}
		So(i, ShouldEqual, 1)

		p = &parser.Parser{}
		ast, err = p.ParseSelect("SELECT id, username, email FROM table")
		So(err, ShouldBeNil)

		plan = NewPlan(table)
		resultPipe, err = plan.SelectPrepare(ast)
		So(err, ShouldBeNil)
		i = 0
		for row := range resultPipe {
			i++
			fmt.Println(row.Id,
				string(row.Username[:utils.Length(row.Username[:])]),
				string(row.Email[:utils.Length(row.Email[:])]))
		}
		So(i, ShouldEqual, InsertCnt)
	})
}
