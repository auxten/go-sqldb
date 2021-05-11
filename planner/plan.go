package planner

import (
	"github.com/gridbase/sqldb/node"
	"github.com/gridbase/sqldb/page"
)

type Plan struct {
	table          *page.Table
	cursor         *page.Cursor
	UnFilteredPipe chan *node.Row
	FilteredPipe   chan *node.Row
	LimitedPipe    chan *node.Row
	ErrorsPipe     chan error
	Stop           chan bool
}

func NewPlan(t *page.Table) (p *Plan) {
	return &Plan{
		table:          t,
		FilteredPipe:   make(chan *node.Row),
		UnFilteredPipe: make(chan *node.Row),
		LimitedPipe:    make(chan *node.Row),
		ErrorsPipe:     make(chan error, 1),
		Stop:           make(chan bool),
	}
}
