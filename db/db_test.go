package db

import (
	"fmt"
	"os"
	"testing"

	"github.com/gridbase/sqldb/node"
	. "github.com/smartystreets/goconvey/convey"
)

func TestDB(t *testing.T) {
	Convey("Open and Close", t, func() {
		const testFile = "test.db"
		defer func() {
			_ = os.Remove(testFile)
		}()
		table, err := Open(testFile)
		So(err, ShouldBeNil)
		So(table.Pager.Pages[0].LeafNode.CommonHeader.IsInternal, ShouldBeFalse)
		So(err, ShouldBeNil)
		So(table, ShouldNotBeNil)

		err = table.Insert(&node.Row{
			Id:       1,
			Username: [32]byte{'a', 'u', 'x', 't', 'e', 'n'},
			Email:    [256]byte{'a', 'u', 'x', 't', 'e', 'n', '@'},
		})
		So(err, ShouldBeNil)

		err = table.Insert(&node.Row{
			Id: 1,
		})
		So(err.Error(), ShouldContainSubstring, "duplicate key 1")

		for i := uint32(2); i < 35; i++ {
			err = table.Insert(&node.Row{
				Id: i,
			})
			fmt.Println(i)
			So(err, ShouldBeNil)
		}
		err = Close(table)
		So(err, ShouldBeNil)

		table, err = Open(testFile)
		So(err, ShouldBeNil)
	})
}
