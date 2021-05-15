package planner

import (
	"os"
	"testing"

	"github.com/auxten/go-sqldb/db"
	"github.com/auxten/go-sqldb/node"
	"github.com/auxten/go-sqldb/parser"
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
				Id: i,
				Sex: func(i int) byte {
					if uint8(i%2) == 1 {
						return 'm'
					} else {
						return 'f'
					}
				}(int(i)),
				Age:      uint8(i % 120),
				Username: [32]byte{'a', 'u', 'x', 't', 'e', 'n', byte('a' + i)},
				Email:    [128]byte{'a', 'u', 'x', 't', 'e', 'n', '@', byte('a' + i)},
				Phone:    [64]byte{'1', '2', '3', '4', '5', '6', '0' + uint8((i/100)%10), '0' + uint8((i/10)%10), '0' + uint8(i%10)},
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
			node.PrintRow(row)
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
			node.PrintRow(row)
		}
		So(i, ShouldEqual, InsertCnt)
	})
}
