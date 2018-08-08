package api

import (
	"errors"
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

type Searcher interface {
	Select(query interface{}, args ...interface{}) Searcher
	Where(query interface{}, values ...interface{}) Searcher
	Not(query interface{}, values ...interface{}) Searcher
	Or(query interface{}, values ...interface{}) Searcher
	Limit(limit interface{}) Searcher
	Offset(offset interface{}) Searcher
	Group(query string) (err error)
	Having(query interface{}, values ...interface{}) Searcher
	Joins(query string, values ...interface{}) Searcher
}
