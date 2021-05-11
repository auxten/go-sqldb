package page

import (
	"fmt"
	"io"
	"os"

	"github.com/auxten/go-sqldb/node"
)

type Page struct {
	// Either InternalNode or LeafNode
	InternalNode *node.InternalNode
	LeafNode     *node.LeafNode
}

func (p *Page) GetMaxKey() uint32 {
	if p.InternalNode != nil {
		return p.InternalNode.ICells[p.InternalNode.Header.KeysNum-1].Key
	} else if p.LeafNode != nil {
		return p.LeafNode.Cells[p.LeafNode.Header.Cells-1].Key
	} else {
		panic("neither Leaf nor Internal node")
	}
}

type Pager struct {
	File    *os.File
	fileLen int64
	PageNum uint32  // PageNum is the boundary of db memory page.
	Pages   []*Page // Page pointer slice, nil member indicates cache missing.
}

func PagerOpen(fileName string) (pager *Pager, err error) {
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
		File:    dbFile,
		fileLen: fileLen,
		PageNum: pageNum,
		Pages:   make([]*Page, node.MaxPages),
	}

	return
}

func (p *Pager) GetPage(pageIdx uint32) (page *Page, err error) {
	if pageIdx >= node.MaxPages {
		return nil, fmt.Errorf("page index %d out of node.MaxPages %d", pageIdx, node.MaxPages)
	}

	if p.Pages[pageIdx] == nil {
		// Cache miss
		// If pageIdx within data file, just read,
		// else just return blank page which will be flushed to db file later.
		if pageIdx <= p.PageNum {
			// Load page from file
			buf := make([]byte, node.PageSize)
			if _, err = p.File.ReadAt(buf, int64(pageIdx*node.PageSize)); err != nil {
				if err != io.EOF {
					return
				}
			}
			// Empty new page will be leaf node
			if buf[0] == 0 {
				// Leaf node
				leaf := &node.LeafNode{}
				if _, err = leaf.Unmarshal(buf); err != nil {
					return
				}
				p.Pages[pageIdx] = &Page{LeafNode: leaf}
			} else {
				// Internal node
				internal := &node.InternalNode{}
				if _, err = internal.Unmarshal(buf); err != nil {
					return
				}
				p.Pages[pageIdx] = &Page{InternalNode: internal}
			}
			if pageIdx >= p.PageNum {
				p.PageNum = pageIdx + 1
			}
		}
	}

	return p.Pages[pageIdx], nil
}

func (p *Pager) Flush(pageIdx uint32) (err error) {
	page := p.Pages[pageIdx]
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
		panic("neither leaf nor internal node")
	}
	_, err = p.File.WriteAt(buf, int64(pageIdx*node.PageSize))

	return
}
