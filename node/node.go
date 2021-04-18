package node

const (
	PageSize = 4096
	MaxPages = 1024
)

var (
	leaf         LeafNode
	internalNode InternalNode

	CommonHeaderSize       = leaf.CommonHeader.Size()
	InternalNodeHeaderSize = CommonHeaderSize + internalNode.Header.Size()
	InternalNodeSize       = internalNode.Size()
	InternalNodeCellSize   = internalNode.ICells[0].Size()
	InternalNodeMaxCells   = len(internalNode.ICells)

	LeafNodeHeaderSize = CommonHeaderSize + leaf.Header.Size()
	LeafNodeSize       = leaf.Size()
	LeafNodeCellSize   = leaf.Cells[0].Size()
	LeafNodeMaxCells   = len(leaf.Cells)

	RightSplitCount = (LeafNodeMaxCells + 1) / 2
	LeftSplitCount  = LeafNodeMaxCells + 1 - RightSplitCount
)
