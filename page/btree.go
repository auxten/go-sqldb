package page

import (
	"github.com/auxten/go-sqldb/node"
)

func InitLeafNode(node *node.LeafNode) {
	node.CommonHeader.IsInternal = false
	node.CommonHeader.IsRoot = false
	node.Header.Cells = 0
	node.Header.NextLeaf = 0
}

func InitInternalNode(node *node.InternalNode) {
	node.CommonHeader.IsInternal = true
	node.CommonHeader.IsRoot = false
	node.Header.KeysNum = 0
}
