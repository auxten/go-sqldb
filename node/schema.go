package node

import (
	"fmt"
)

var (
	RowSize = (&Row{}).Size()
)

func printRow(row *Row) {
	fmt.Printf("%d %s %s", row.Id, row.Username, row.Email)
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
