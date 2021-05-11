package planner

import (
	"strconv"
	"strings"

	"github.com/auxten/go-sqldb/node"
	"github.com/auxten/go-sqldb/parser"
)

func (plan *Plan) Insert(ast *parser.InsertTree) (count int, err error) {
	for _, r := range ast.Values {
		// 这里暂时都假定我们插入的 Schema 是固定的 node.Row 类型
		// 根据 InsertTree.Columns 的字段顺序，我们强制类型转换还原出 node.Row 结构
		row := &node.Row{}
		if len(ast.Columns) == 0 {
			ast.Columns = []string{"id", "username", "email"}
		}
		for i, col := range ast.Columns {
			switch strings.ToUpper(col) {
			case "ID":
				var parsed int64
				if parsed, err = strconv.ParseInt(r[i], 10, 64); err != nil {
					return
				}
				row.Id = uint32(parsed)
			case "USERNAME":
				copy(row.Username[:], r[i])
			case "EMAIL":
				copy(row.Email[:], r[i])
			}
		}

		if err = plan.table.Insert(row); err != nil {
			return
		} else {
			count++
		}
	}
	return
}
