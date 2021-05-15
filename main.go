package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/auxten/go-sqldb/db"
	"github.com/auxten/go-sqldb/node"
	"github.com/auxten/go-sqldb/page"
	"github.com/auxten/go-sqldb/parser"
	"github.com/auxten/go-sqldb/planner"
	"github.com/auxten/go-sqldb/utils"
)

func main() {
	if len(os.Args) < 2 {
		panic("Database file name must provide")
	}

	var (
		err        error
		t          *page.Table
		resultPipe chan *node.Row
	)

	fileName := os.Args[1]
	if t, err = db.Open(fileName); err != nil {
		panic(err)
	}
	defer db.Close(t)

	http.HandleFunc("/query", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		q := request.URL.Query()
		query := q.Get("q")
		if query != "" {
			var (
				ast *parser.SelectTree
			)

			p := &parser.Parser{}
			if p.GetSQLType(query) != parser.SELECT {
				_, _ = fmt.Fprintf(writer, "not a SELECT statement")
				return
			}
			if ast, err = p.ParseSelect(query); err != nil {
				_, _ = fmt.Fprintf(writer, "parse %s, error: %v", query, err)
				return
			}

			plan := planner.NewPlan(t)
			if resultPipe, err = plan.SelectPrepare(ast); err != nil {
				_, _ = fmt.Fprintf(writer, "%s", err)
				return
			}

			for row := range resultPipe {
				if len(ast.Projects) == 1 && ast.Projects[0] == parser.ASTERISK {
					node.WriteRow(writer, row)
					continue
				}
				var outRow = make([]string, 0, 3)
				for _, proj := range ast.Projects {
					switch strings.ToUpper(proj) {
					case "ID":
						outRow = append(outRow, fmt.Sprintf("%d", row.Id))
					case "SEX":
						outRow = append(outRow, fmt.Sprintf("%c", row.Sex))
					case "AGE":
						outRow = append(outRow, fmt.Sprintf("%d", row.Age))
					case "USERNAME":
						outRow = append(outRow, string(row.Username[:utils.Length(row.Username[:])]))
					case "EMAIL":
						outRow = append(outRow, string(row.Email[:utils.Length(row.Email[:])]))
					case "PHONE":
						outRow = append(outRow, string(row.Phone[:utils.Length(row.Phone[:])]))
					}
				}
				_, _ = fmt.Fprint(writer, strings.Join(outRow, "\t"), "\n")
			}
			return
		}
		_, _ = fmt.Fprintf(writer, "need /query?q=SELECT")
		return
	})
	http.HandleFunc("/exec", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "application/json")
		q := request.URL.Query()
		query := q.Get("q")
		if query != "" {
			var (
				ast *parser.InsertTree
				cnt int
			)

			p := &parser.Parser{}
			if p.GetSQLType(query) != parser.INSERT {
				_, _ = fmt.Fprintf(writer, "not a INSERT statement")
				return
			}
			ast, err = p.ParseInsert(query)

			plan := planner.NewPlan(t)
			if cnt, err = plan.Insert(ast); err != nil {
				_, _ = fmt.Fprintf(writer, "%s", err)
				return
			}
			_, _ = fmt.Fprintf(writer, "%d", cnt)
			return
		}
		_, _ = fmt.Fprintf(writer, "need /exec?q=INSERT")
		return
	})
	go http.ListenAndServe(":8080", nil)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
}
