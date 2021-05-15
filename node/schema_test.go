package node

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSchema(t *testing.T) {
	Convey("", t, func() {
		dumpConst()
	})
}

type node struct {
	name  string
	idx   uint32
	cells [2]cell
}

type cell struct {
	Key   uint32
	Value [230]byte
}

func TestStructAssignment(t *testing.T) {
	Convey("cell assign", t, func() {
		c1 := Cell{
			Key:   1,
			Value: [230]byte{'1', '1'},
		}
		c2 := Cell{
			Key:   2,
			Value: [230]byte{'2', '2'},
		}
		So(c2.Key, ShouldEqual, 2)
		So(c2.Value[0], ShouldEqual, '2')
		So(c2.Value[1], ShouldEqual, '2')

		c2 = c1
		So(c2.Key, ShouldEqual, 1)
		So(c2.Value[0], ShouldEqual, '1')
		So(c2.Value[1], ShouldEqual, '1')
	})

	Convey("node assign", t, func() {
		n1 := node{
			name: "n1",
			idx:  1,
			cells: [2]cell{
				{
					Key:   1,
					Value: [230]byte{'1', '1'},
				},
				{
					Key:   2,
					Value: [230]byte{'2', '2'},
				},
			},
		}
		n2 := node{
			name: "n2",
			idx:  2,
			cells: [2]cell{
				{
					Key:   3,
					Value: [230]byte{'3', '3'},
				},
				{
					Key:   4,
					Value: [230]byte{'4', '4'},
				},
			},
		}

		n2 = n1
		So(n2.name, ShouldEqual, "n1")
		So(n2.idx, ShouldEqual, 1)
		So(n2.cells[0].Key, ShouldEqual, 1)
		So(n2.cells[0].Value[0], ShouldEqual, '1')
		So(n2.cells[0].Value[1], ShouldEqual, '1')
		So(n2.cells[1].Key, ShouldEqual, 2)
		So(n2.cells[1].Value[0], ShouldEqual, '2')
		So(n2.cells[1].Value[1], ShouldEqual, '2')
	})
}
