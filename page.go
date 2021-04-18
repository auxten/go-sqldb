package main

import (
	"fmt"
	"io"
	"os"

	"github.com/gridbase/sqldb/node"
)

type Page struct {
	// Either InternalNode or LeafNode
	InternalNode *node.InternalNode
	LeafNode     *node.LeafNode
}

type Pager struct {
	file    *os.File
	fileLen int64
	pageNum uint32  // pageNum is the boundary of db memory page.
	pages   []*Page // Page pointer slice, nil member indicates cache missing.
}

func pagerOpen(fileName string) (pager *Pager, err error) {
	var (
		dbFile  *os.File
		fileLen int64
	)
	if dbFile, err = os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0600); err != nil {
		return
	}

	// get file length
	if fileLen, err = dbFile.Seek(0, io.SeekEnd); err != nil {
		return
	}

	// dbFile length must be n * node.PageSize, node.PageSize is usually 4096
	if fileLen%node.PageSize != 0 {
		return
	}

	pageNum := uint32(fileLen / node.PageSize)
	if pageNum >= node.MaxPages {
		panic("file length exceeds max pages limit")
	}
	pager = &Pager{
		file:    dbFile,
		fileLen: fileLen,
		pageNum: pageNum,
		pages:   make([]*Page, node.MaxPages),
	}

	return
}

func (p *Pager) getPage(pageIdx uint32) (page *Page, err error) {
	if pageIdx >= node.MaxPages {
		return nil, fmt.Errorf("page index %d out of node.MaxPages %d", pageIdx, node.MaxPages)
	}

	if p.pages[pageIdx] == nil {
		// Cache miss
		// If pageIdx within data file, just read,
		// else just return blank page which will be flushed to db file later.
		if pageIdx <= p.pageNum {
			// Load page from file
			buf := make([]byte, node.PageSize)
			if _, err = p.file.ReadAt(buf, int64(pageIdx*node.PageSize)); err != nil {
				if err != io.EOF {
					return
				}
			}
			if buf[0] == 0 {
				// Leaf node
				leaf := &node.LeafNode{}
				if _, err = leaf.Unmarshal(buf); err != nil {
					return
				}
				p.pages[pageIdx] = &Page{LeafNode: leaf}
			} else {
				// Internal node
				internal := &node.InternalNode{}
				if _, err = internal.Unmarshal(buf); err != nil {
					return
				}
				p.pages[pageIdx] = &Page{InternalNode: internal}
			}
			if pageIdx >= p.pageNum {
				p.pageNum = pageIdx + 1
			}
		}
	}

	return p.pages[pageIdx], nil
}

func (p *Pager) Flush(pageIdx uint32) (err error) {
	page := p.pages[pageIdx]
	if page == nil {
		return fmt.Errorf("flushing nil page")
	}

	buf := make([]byte, node.PageSize)
	if page.LeafNode != nil {
		if _, err = page.LeafNode.Marshal(buf); err != nil {
			return
		}
	} else if page.InternalNode != nil {
		if _, err = page.InternalNode.Marshal(buf); err != nil {
			return
		}
	} else {
		panic("neither Leaf nor Internal node")
	}
	_, err = p.file.WriteAt(buf, int64(pageIdx*node.PageSize))

	return
}
