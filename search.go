package iodata

import (
	"errors"
	"fmt"

	"github.com/moisespsena-go/iodata/api"
	"github.com/moisespsena-go/error-wrap"
)

var (
	ErrInvalidSQL = errors.New("invalid SQL")
)

// SQL expression
type expr struct {
	expr string
	args []interface{}
}

// Expr generate raw SQL expression, for example:
//     DB.Model(&product).Update("price", gorm.Expr("price * ? + ?", 2, 100))
func Expr(expression string, args ...interface{}) *expr {
	return &expr{expr: expression, args: args}
}

type Search struct {
	limit            interface{}
	offset           interface{}
	group            string
	query            string
	whereConditions  []map[string]interface{}
	orConditions     []map[string]interface{}
	notConditions    []map[string]interface{}
	havingConditions []map[string]interface{}
	joinConditions   []map[string]interface{}
	selects          map[string]interface{}
	orders           []interface{}
	omits            []string
	raw              bool
	ignoreOrderQuery bool
}

func (s *Search) Select(query interface{}, args ...interface{}) api.Searcher {
	s.selects = map[string]interface{}{"query": query, "args": args}
	return s
}

func (s *Search) Where(query interface{}, values ...interface{}) api.Searcher {
	s.whereConditions = append(s.whereConditions, map[string]interface{}{"query": query, "args": values})
	return s
}

func (s *Search) Not(query interface{}, values ...interface{}) api.Searcher {
	s.notConditions = append(s.notConditions, map[string]interface{}{"query": query, "args": values})
	return s
}

func (s *Search) Or(query interface{}, values ...interface{}) api.Searcher {
	s.orConditions = append(s.orConditions, map[string]interface{}{"query": query, "args": values})
	return s
}

func (s *Search) Limit(limit interface{}) api.Searcher {
	s.limit = limit
	return s
}

func (s *Search) Offset(offset interface{}) api.Searcher {
	s.offset = offset
	return s
}

func (s *Search) Group(query string) (err error) {
	s.group, err = s.getInterfaceAsSQL(query)
	return
}

func (s *Search) getInterfaceAsSQL(value interface{}) (str string, err error) {
	switch value.(type) {
	case string, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		str = fmt.Sprintf("%v", value)
	default:
		return "", errwrap.Wrap(ErrInvalidSQL, "SQL: %q", value)
	}

	if str == "-1" {
		return "", nil
	}
	return
}

func (s *Search) Having(query interface{}, values ...interface{}) api.Searcher {
	if val, ok := query.(*expr); ok {
		s.havingConditions = append(s.havingConditions, map[string]interface{}{"query": val.expr, "args": val.args})
	} else {
		s.havingConditions = append(s.havingConditions, map[string]interface{}{"query": query, "args": values})
	}
	return s
}

func (s *Search) Joins(query string, values ...interface{}) api.Searcher {
	s.joinConditions = append(s.joinConditions, map[string]interface{}{"query": query, "args": values})
	return s
}
