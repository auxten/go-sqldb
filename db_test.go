package main

import (
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDB(t *testing.T) {
	Convey("Open and Close", t, func() {
		const testFile = "test.db"
		defer func() {
			_ = os.Remove(testFile)
		}()
		table, err := dbOpen(testFile)
		So(table.pager.pages[0].LeafNode.CommonHeader.IsInternal, ShouldBeFalse)
		So(err, ShouldBeNil)
		So(table, ShouldNotBeNil)

		err = dbClose(table)
		So(err, ShouldBeNil)

		table, err = dbOpen(testFile)
		So(err, ShouldBeNil)
	})
}
