package page

import (
	"fmt"

	"github.com/auxten/go-sqldb/node"
)

type Table struct {
	Pager       *Pager
	RootPageIdx uint32
}

// Seek the page of key, if not exist then return the place key should be
// for the later INSERT.
func (table *Table) Seek(key uint32) (cursor *Cursor, err error) {
	var (
		rootPage *Page
	)

	if rootPage, err = table.Pager.GetPage(table.RootPageIdx); err != nil {
		return
	}
	if rootPage.LeafNode != nil {
		return table.leafNodeSeek(table.RootPageIdx, key)
	} else if rootPage.InternalNode != nil {
		return table.internalNodeSeek(table.RootPageIdx, key)
	} else {
		panic("root page type")
	}
	return
}

func (table *Table) Insert(row *node.Row) (err error) {
	var (
		p   *Page
		cur *Cursor
	)

	if cur, err = table.Seek(row.Id); err != nil {
		return
	}
	if p, err = table.Pager.GetPage(cur.PageIdx); err != nil {
		return
	}
	// Must be leaf node
	if p.LeafNode == nil {
		panic("should be leaf node")
	}
	if cur.CellIdx < p.LeafNode.Header.Cells {
		if p.LeafNode.Cells[cur.CellIdx].Key == row.Id {
			return fmt.Errorf("duplicate key %d", row.Id)
		}
	}

	return cur.LeafNodeInsert(row.Id, row)
}

func (table *Table) leafNodeSeek(pageIdx uint32, key uint32) (cursor *Cursor, err error) {
	var (
		p                 *Page
		minIdx, maxIdx, i uint32
	)

	if p, err = table.Pager.GetPage(pageIdx); err != nil {
		return
	}
	maxIdx = p.LeafNode.Header.Cells

	cursor = &Cursor{
		Table:      table,
		PageIdx:    pageIdx,
		EndOfTable: false,
	}

	// Walk the btree
	for i = maxIdx; i != minIdx; {
		index := (minIdx + i) / 2
		keyIdx := p.LeafNode.Cells[index].Key
		if key == keyIdx {
			cursor.CellIdx = index
			return
		}
		if key < keyIdx {
			i = index
		} else {
			minIdx = index + 1
		}
	}

	cursor.CellIdx = minIdx
	return
}

func (table *Table) internalNodeSeek(pageIdx uint32, key uint32) (cursor *Cursor, err error) {
	var (
		p, childPage *Page
	)

	if p, err = table.Pager.GetPage(pageIdx); err != nil {
		return
	}

	nodeIdx := p.InternalNode.FindChildByKey(key)
	childIdx := *p.InternalNode.Child(nodeIdx)

	if childPage, err = table.Pager.GetPage(childIdx); err != nil {
		return
	}
	if childPage.InternalNode != nil {
		return table.internalNodeSeek(childIdx, key)
	} else if childPage.LeafNode != nil {
		return table.leafNodeSeek(childIdx, key)
	}
	return
}

func (table *Table) CreateNewRoot(rightChildPageIdx uint32) (err error) {
	/*
	  Handle splitting the root.
	  Old root copied to new page, becomes left child.
	  Address of right child passed in.
	  Re-initialize root page to contain the new root node.
	  New root node points to two children.
	*/
	var (
		rootPage, rightChildPage, leftChildPage *Page
	)
	if rootPage, err = table.Pager.GetPage(table.RootPageIdx); err != nil {
		return
	}
	if rightChildPage, err = table.Pager.GetPage(rightChildPageIdx); err != nil {
		return
	}
	leftChildPageIdx := table.Pager.PageNum
	if leftChildPage, err = table.Pager.GetPage(leftChildPageIdx); err != nil {
		return
	}

	// copy whatever kind of node to leftChildPage, and set nonRoot
	if rootPage.LeafNode != nil {
		*leftChildPage.LeafNode = *rootPage.LeafNode
		leftChildPage.LeafNode.CommonHeader.IsRoot = false
	} else if rootPage.InternalNode != nil {
		*leftChildPage.InternalNode = *rootPage.InternalNode
		leftChildPage.InternalNode.CommonHeader.IsRoot = false
	}

	// 重新初始化 root page，root page 将会有一个 key，两个子节点
	rootPage.LeafNode = nil
	rootPage.InternalNode = new(node.InternalNode)
	rootNode := rootPage.InternalNode
	InitInternalNode(rootNode)
	rootNode.CommonHeader.IsRoot = true
	rootNode.Header.KeysNum = 1
	childPageIdxPtr := rootNode.Child(0)
	*(childPageIdxPtr) = leftChildPageIdx
	leftChildMaxKey := leftChildPage.GetMaxKey()
	rootNode.ICells[0].Key = leftChildMaxKey
	rootNode.Header.RightChild = rightChildPageIdx
	if leftChildPage.LeafNode != nil {
		leftChildPage.LeafNode.CommonHeader.Parent = table.RootPageIdx
	} else if leftChildPage.InternalNode != nil {
		leftChildPage.InternalNode.CommonHeader.Parent = table.RootPageIdx
	}
	if rightChildPage.LeafNode != nil {
		rightChildPage.LeafNode.CommonHeader.Parent = table.RootPageIdx
	} else if rightChildPage.InternalNode != nil {
		rightChildPage.InternalNode.CommonHeader.Parent = table.RootPageIdx
	}

	return
}

func (table *Table) InternalNodeInsert(parentPageIdx uint32, childPageIdx uint32) (err error) {
	/*
	  Add a new child/key pair to parent that corresponds to child
	*/
	var (
		parentPage, childPage, rightChildPage *Page
	)

	if parentPage, err = table.Pager.GetPage(parentPageIdx); err != nil {
		return
	}
	if childPage, err = table.Pager.GetPage(childPageIdx); err != nil {
		return
	}
	childMaxKey := childPage.GetMaxKey()
	index := parentPage.InternalNode.FindChildByKey(childMaxKey)
	originalKeyCnt := parentPage.InternalNode.Header.KeysNum
	parentPage.InternalNode.Header.KeysNum += 1

	if parentPage.InternalNode.Header.KeysNum > node.InternalNodeMaxCells {
		panic("InternalNodeMaxCells exceeds")
	}

	rightChildPageIdx := parentPage.InternalNode.Header.RightChild
	if rightChildPage, err = table.Pager.GetPage(rightChildPageIdx); err != nil {
		return
	}

	if childMaxKey > rightChildPage.GetMaxKey() {
		/* Replace right child */
		*parentPage.InternalNode.Child(originalKeyCnt) = rightChildPageIdx
		parentPage.InternalNode.ICells[originalKeyCnt].Key = rightChildPage.GetMaxKey()
		parentPage.InternalNode.Header.RightChild = childPageIdx
	} else {
		/* Make room for the new cell */
		for i := originalKeyCnt; i > index; i-- {
			parentPage.InternalNode.ICells[i] = parentPage.InternalNode.ICells[i-1]
		}
		*parentPage.InternalNode.Child(index) = childPageIdx
		parentPage.InternalNode.ICells[index].Key = childMaxKey
	}
	return
}

func (table *Table) Select() {

}

func (table *Table) Prepare() {

}
