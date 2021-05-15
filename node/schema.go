package node

import (
	"fmt"
	"io"
	"os"

	"github.com/auxten/go-sqldb/utils"
)

var (
	RowSize = (&Row{}).Size()
)

/*
	Id       uint32
	Sex      byte
	Age      uint8
	Username [32]byte
	Email    [128]byte
	Phone    [64]byte
*/
func PrintRow(row *Row) {
	_, _ = WriteRow(os.Stdout, row)
}

func WriteRow(w io.Writer, row *Row) (int, error) {
	return fmt.Fprintf(w, "%d\t%c\t%d\t%s\t%s\t%s\n",
		row.Id,
		row.Sex,
		row.Age,
		string(row.Username[:utils.Length(row.Username[:])]),
		string(row.Email[:utils.Length(row.Email[:])]),
		string(row.Phone[:utils.Length(row.Phone[:])]),
	)
}

func dumpConst() {
	if LeafNodeSize > PageSize {
		panic("LeafNode too big")
	}
	if RowSize > LeafNodeCellSize {
		panic("Row too big")
	}

	fmt.Printf("Row Size %d\n", RowSize)
	fmt.Printf("Common Header Size %d\n", CommonHeaderSize)
	fmt.Printf("InternalNode Header Size %d\n", InternalNodeHeaderSize)
	fmt.Printf("InternalNode Size %d\n", InternalNodeSize)
	fmt.Printf("InternalNode Cell Size %d\n", InternalNodeCellSize)
	fmt.Printf("InternalNode Max Cell %d\n", InternalNodeMaxCells)
	fmt.Printf("LeafNode Header Size %d\n", LeafNodeHeaderSize)
	fmt.Printf("LeafNode Size %d\n", LeafNodeSize)
	fmt.Printf("LeafNode Cell Size %d\n", LeafNodeCellSize)
	fmt.Printf("LeafNode Max Cell %d\n", LeafNodeMaxCells)
	fmt.Printf("LeftSplitCount %d\n", LeftSplitCount)
	fmt.Printf("RightSplitCount %d\n", RightSplitCount)
}
