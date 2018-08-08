package query

import (
	"github.com/lfittl/pg_query_go"
	nodes "github.com/lfittl/pg_query_go/nodes"
	"github.com/moisespsena-go/iodata/api"
)

type QueryTable struct {
	Name   string
	Alias  *QueryTable
	Header api.DataHeader
	//Columns map of header index from header name
	Columns  map[string]int
	Fields   map[string]bool
	RealName string
}

type QueryFunc struct {
	Args []string
}

type Query struct {
	Context *api.DataSourceContext
	Source  string
	Args    []interface{}
	// tables tables used on query
	tables                                              map[string]*QueryTable
	hasSelect, hasInsert, HasUpdate, hasDelete, hasFunc bool
	funcs                                               map[string]*QueryFunc
	nodes                                               []nodes.Node
	tree                                                pg_query.ParsetreeList
}

func (q *Query) NewCompiler() *QueryString {
	return &QueryString{Tree: q.tree}
}
