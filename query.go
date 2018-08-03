package iodata

import (
	"github.com/moisespsena-go/iodata/api"
	"github.com/moisespsena/go-error-wrap"
	"github.com/xwb1989/sqlparser"
)

type QueryTable struct {
	Name   string
	Alias  *QueryTable
	Header *api.DataHeader
	//Columns map of header index from header name
	Columns map[string]int
	Fields  map[string]bool
}

type QueryFunc struct {
	Args []string
}

type QueryContext struct {
	Tables map[string]*api.DataHeader
}

type Query struct {
	Context *QueryContext
	Source  string
	Args    []interface{}
	// tables tables used on query
	tables                                              map[string]*QueryTable
	hasSelect, hasInsert, HasUpdate, hasDelete, hasFunc bool
	funcs                                               map[string]*QueryFunc
}

func (q *Query) validateTableExpr(e sqlparser.TableExpr) (err error) {
	switch e.(type) {
	case *sqlparser.AliasedTableExpr:
	case *sqlparser.JoinTableExpr:
	case *sqlparser.ParenTableExpr:
	}
	return
}

func (q *Query) validateTableExprs(s sqlparser.TableExprs) (err error) {
	for i, expr := range s {
		if err = q.validateTableExpr(expr); err != nil {
			return errwrap.Wrap(err, "Table Expression %d", i)
		}
	}
	return
}

func (q *Query) validateSelect(s *sqlparser.Select) error {

	err := q.validateTableExprs(s.From)
	return err
}

func (q *Query) Build() error {
	return nil
}
