package node

import (
	"fmt"
)

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
	InternalNodeMaxCells   = uint32(len(internalNode.ICells))

	LeafNodeHeaderSize = CommonHeaderSize + leaf.Header.Size()
	LeafNodeSize       = leaf.Size()
	LeafNodeCellSize   = leaf.Cells[0].Size()
	LeafNodeMaxCells   = uint32(len(leaf.Cells))

	RightSplitCount = (LeafNodeMaxCells + 1) / 2
	LeftSplitCount  = LeafNodeMaxCells + 1 - RightSplitCount
)

// FindChildByKey returns the index of the child which should contain
//  the given key.
func (d *InternalNode) FindChildByKey(key uint32) uint32 {
	var (
		minIdx = uint32(0)
		maxIdx = d.Header.KeysNum
	)
	for minIdx != maxIdx {
		idx := (minIdx + maxIdx) / 2
		rightKey := d.ICells[idx].Key
		if rightKey >= key {
			maxIdx = idx
		} else {
			minIdx = idx + 1
		}
	}

	return minIdx
}

func (d *InternalNode) Child(childIdx uint32) (ptr *uint32) {
	keysNum := d.Header.KeysNum
	if childIdx > keysNum {
		panic(fmt.Sprintf("childIdx %d out of keysNum %d", childIdx, keysNum))
	} else if childIdx == keysNum {
		return &d.Header.RightChild
	} else {
		return &d.ICells[childIdx].Key
	}
}
