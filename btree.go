package main

import (
	"github.com/gridbase/sqldb/node"
)

/*

 */

func initLeafNode(node *node.LeafNode) {
	node.CommonHeader.IsInternal = false
	node.CommonHeader.IsRoot = false
	node.Header.Cells = 0
	node.Header.NextLeaf = 0
}
