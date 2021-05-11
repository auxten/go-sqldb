package planner

import (
	"fmt"
	"go/token"
	"go/types"
	"strings"

	"github.com/auxten/go-sqldb/node"
	"github.com/auxten/go-sqldb/page"
	"github.com/auxten/go-sqldb/parser"
)

func (plan *Plan) SelectPrepare(ast *parser.SelectTree) (filteredPipe chan *node.Row, err error) {
	var (
		table = plan.table
		p     *page.Page
	)
	// Get the start of table
	if plan.cursor, err = table.Seek(0); err != nil {
		return
	}
	if p, err = table.Pager.GetPage(plan.cursor.PageIdx); err != nil {
		return
	}
	// If the key not found in the table and the page to insert that key
	// has 0 cells. We got the end of the table.
	plan.cursor.EndOfTable = p.LeafNode.Header.Cells == 0

	/*
		The code below demonstrates a simple "Volcano Model" query plan.
		 For more please refer to https://doi.org/10.1109/69.273032
	*/

	// Fetch rows from storage pages
	go func(out chan<- *node.Row) {
		var (
			err error
			row *node.Row
		)
		defer close(out)
		if row, err = plan.fetchRow(table); err != nil {
			plan.ErrorsPipe <- err
			return
		}
		out <- row
		for !plan.cursor.EndOfTable {
			select {
			case <-plan.Stop:
				return
			case out <- row:
				if row, err = plan.fetchRow(table); err != nil {
					plan.ErrorsPipe <- err
					return
				}
			}
		}
	}(plan.UnFilteredPipe)

	// Filter rows according the ast.Where
	go func(in <-chan *node.Row, out chan<- *node.Row, where []string) {
		defer close(out)
		for row := range in {
			if len(where) == 0 {
				out <- row
				continue
			}
			filtered, err := isRowFiltered(where, row)
			if err != nil {
				plan.ErrorsPipe <- err
				return
			}
			if !filtered {
				out <- row
			}
		}
	}(plan.UnFilteredPipe, plan.FilteredPipe, ast.Where)

	// Count row count for LIMIT clause.
	go func(in <-chan *node.Row, out chan<- *node.Row, limit int64) {
		defer close(out)
		defer close(plan.Stop)
		i := int64(0)
		for row := range in {
			i++
			if i > limit && limit > 0 {
				return
			}
			out <- row
		}
	}(plan.FilteredPipe, plan.LimitedPipe, ast.Limit)

	return plan.LimitedPipe, nil
}

func (plan *Plan) fetchRow(table *page.Table) (row *node.Row, err error) {
	var p *page.Page
	if p, err = table.Pager.GetPage(plan.cursor.PageIdx); err != nil {
		return
	}
	row = new(node.Row)
	if _, err = row.Unmarshal(p.LeafNode.Cells[plan.cursor.CellIdx].Value[:]); err != nil {
		row = nil
		return
	}
	plan.cursor.CellIdx++
	if plan.cursor.CellIdx >= p.LeafNode.Header.Cells {
		// Move cursor to next leaf
		if p.LeafNode.Header.NextLeaf == 0 {
			// 已经移动到了最右的的叶子节点
			plan.cursor.EndOfTable = true
		} else {
			plan.cursor.PageIdx = p.LeafNode.Header.NextLeaf
			plan.cursor.CellIdx = 0
		}
	}
	return
}

func isRowFiltered(where []string, row *node.Row) (filtered bool, err error) {
	// This is a very dirty hack to use Eval to evaluate the Where statement.
	var (
		normalized = make([]string, len(where))
		tv         types.TypeAndValue
	)

	/*
		struct Row {
		  Id uint32
		  Username [32]byte
		  Email [256]byte
		}
	*/
	for i, w := range where {
		upper := strings.ToUpper(w)
		switch upper {
		case "AND":
			normalized[i] = "&&"
		case "OR":
			normalized[i] = "||"
		case "ID":
			normalized[i] = fmt.Sprintf("%d", row.Id)
		case "USERNAME":
			normalized[i] = fmt.Sprintf("%s", string(row.Username[:]))
		case "EMAIL":
			normalized[i] = fmt.Sprintf("%s", string(row.Email[:]))
		default:
			normalized[i] = w
		}
	}
	expr := strings.Join(normalized, " ")
	fSet := token.NewFileSet()
	if tv, err = types.Eval(fSet, nil, token.NoPos, expr); err != nil {
		return
	}
	if tv.Type == nil {
		err = fmt.Errorf("eval(%q) got nil type but no error", expr)
		return
	}
	if !strings.Contains(tv.Type.String(), "bool") {
		err = fmt.Errorf("eval(%q) got non bool type", expr)
		return
	}
	if tv.Value.ExactString() == "true" {
		filtered = false
	} else {
		filtered = true
	}
	return
}
