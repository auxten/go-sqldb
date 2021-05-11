package page

import (
	"github.com/gridbase/sqldb/node"
)

type Cursor struct {
	Table      *Table
	PageIdx    uint32
	CellIdx    uint32
	EndOfTable bool
}

func (cursor *Cursor) LeafNodeInsert(key uint32, row *node.Row) (err error) {
	var (
		page *Page
	)

	if page, err = cursor.Table.Pager.GetPage(cursor.PageIdx); err != nil {
		return
	}
	cells := page.LeafNode.Header.Cells
	if cells >= node.LeafNodeMaxCells {
		// Split leaf node
		if err = cursor.LeafNodeSplitInsert(key, row); err != nil {
			return
		}
		return
	}

	if cursor.CellIdx < cells {
		// Need make room for new cell
		for i := cells; i > cursor.CellIdx; i-- {
			page.LeafNode.Cells[i] = page.LeafNode.Cells[i-1]
		}
	}
	page.LeafNode.Header.Cells += 1
	cell := &page.LeafNode.Cells[cursor.CellIdx]
	err = saveToCell(cell, key, row)
	return
}

func (cursor *Cursor) LeafNodeSplitInsert(key uint32, row *node.Row) (err error) {
	/*
	  Create a new node and move half the cells over.
	  Insert the new value in one of the two nodes.
	  Update parent or create a new parent.
	*/
	var (
		oldMaxKey, newPageNum uint32
		oldPage, newPage      *Page
		parentPage            *Page
		pager                 *Pager
	)
	pager = cursor.Table.Pager
	if oldPage, err = pager.GetPage(cursor.PageIdx); err != nil {
		return
	}
	oldMaxKey = oldPage.GetMaxKey()
	newPageNum = pager.PageNum
	// put new page in the end
	// TODO: Page recycle
	if newPage, err = pager.GetPage(newPageNum); err != nil {
		return
	}
	InitLeafNode(newPage.LeafNode)
	newPage.LeafNode.CommonHeader.Parent = oldPage.LeafNode.CommonHeader.Parent
	newPage.LeafNode.Header.NextLeaf = oldPage.LeafNode.Header.NextLeaf
	oldPage.LeafNode.Header.NextLeaf = newPageNum

	/*
	  All existing keys plus new key should should be divided
	  evenly between old (left) and new (right) nodes.
	  Starting from the right, move each key to correct position.
	*/
	for i := node.LeafNodeMaxCells; ; i-- {
		if i+1 == 0 {
			break
		}
		var destPage *Page
		if i > node.LeftSplitCount {
			destPage = newPage
		} else {
			destPage = oldPage
		}
		cellIdx := i % node.LeftSplitCount
		destCell := &destPage.LeafNode.Cells[cellIdx]

		if i == cursor.CellIdx {
			if err = saveToCell(destCell, key, row); err != nil {
				return
			}
		} else if i > cursor.CellIdx {
			*destCell = oldPage.LeafNode.Cells[i-1]
		} else {
			*destCell = oldPage.LeafNode.Cells[i]
		}
	}

	/* Update cell count on both leaf nodes */
	oldPage.LeafNode.Header.Cells = node.LeftSplitCount
	newPage.LeafNode.Header.Cells = node.RightSplitCount

	if oldPage.LeafNode.CommonHeader.IsRoot {
		return cursor.Table.CreateNewRoot(newPageNum)
	} else {
		parentPageIdx := oldPage.LeafNode.CommonHeader.Parent
		if parentPage, err = pager.GetPage(parentPageIdx); err != nil {
			return
		}
		// parent page is an internal node
		oldChildIdx := parentPage.InternalNode.FindChildByKey(oldMaxKey)
		if oldChildIdx >= node.InternalNodeMaxCells {
			panic("InternalNodeMaxCells exceeds")
		}
		parentPage.InternalNode.ICells[oldChildIdx].Key = oldPage.GetMaxKey()
		err = cursor.Table.InternalNodeInsert(parentPageIdx, newPageNum)
	}
	return
}

func saveToCell(cell *node.Cell, key uint32, row *node.Row) (err error) {
	rowBuf := make([]byte, row.Size())
	if _, err = row.Marshal(rowBuf); err != nil {
		return
	}
	cell.Key = key
	copy(cell.Value[:], rowBuf)
	return
}
