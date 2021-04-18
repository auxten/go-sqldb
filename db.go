package main

type Table struct {
	pager       *Pager
	rootPageNum int
}

type Cursor struct {
	table      *Table
	pageNum    int
	cellNum    int
	endOfTable bool
}

func dbOpen(fileName string) (table *Table, err error) {
	var (
		pager *Pager
		page  *Page
	)
	if pager, err = pagerOpen(fileName); err != nil {
		return
	}

	table = &Table{
		pager:       pager,
		rootPageNum: 0,
	}

	if pager.pageNum == 0 {
		// New database file, initialize page 0 as leaf node.
		if page, err = pager.getPage(0); err != nil {
			return
		}
		initLeafNode(page.LeafNode)
		page.LeafNode.CommonHeader.IsRoot = true
	}

	return
}

func dbClose(table *Table) (err error) {
	pager := table.pager
	for i := uint32(0); i < pager.pageNum; i++ {
		if pager.pages[i] != nil {
			if err = pager.Flush(i); err != nil {
				return
			}
			pager.pages[i] = nil
		}
	}

	err = pager.file.Close()

	for i := range pager.pages {
		pager.pages[i] = nil
	}
	return
}
