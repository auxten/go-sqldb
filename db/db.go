package db

import (
	"github.com/auxten/go-sqldb/page"
)

func Open(fileName string) (t *page.Table, err error) {
	var (
		pager *page.Pager
		p     *page.Page
	)
	if pager, err = page.PagerOpen(fileName); err != nil {
		return
	}

	t = &page.Table{
		Pager:       pager,
		RootPageIdx: 0,
	}

	if pager.PageNum == 0 {
		// New database file, initialize page 0 as leaf node.
		if p, err = pager.GetPage(0); err != nil {
			return
		}
		page.InitLeafNode(p.LeafNode)
		p.LeafNode.CommonHeader.IsRoot = true
	}

	return
}

func Close(t *page.Table) (err error) {
	pager := t.Pager
	for i := uint32(0); i < pager.PageNum; i++ {
		if pager.Pages[i] != nil {
			if err = pager.Flush(i); err != nil {
				return
			}
			pager.Pages[i] = nil
		}
	}

	err = pager.File.Close()

	for i := range pager.Pages {
		pager.Pages[i] = nil
	}
	return
}
