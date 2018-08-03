package iodata

import (
	"reflect"

	"errors"

	"strings"

	"github.com/lfittl/pg_query_go"
	nodes "github.com/lfittl/pg_query_go/nodes"
	"github.com/moisespsena/go-error-wrap"
)

var (
	ErrTableNotFound  = errors.New("Table does not exists")
	ErrColumnNotFound = errors.New("Column does not exists")
)

func (q *Query) parseError(node interface{}, err interface{}) error {
	var e error
	if es, ok := err.(string); ok {
		e = errors.New(es)
	} else {
		e = err.(error)
	}

	f := reflect.ValueOf(node).FieldByName("Location")
	if f.IsValid() {
		return errwrap.Wrap(e, "SQL:\n"+q.Source+"\n"+strings.Repeat(" ", int(f.Int()))+"^\n")
	}
	return e
}

var lower = strings.ToLower

func (q *Query) parse(nods ...interface{}) (err error) {
	for _, node := range nods {
		if node == nil {
			continue
		}
		switch nt := node.(type) {
		case nodes.List:
			for _, e := range nt.Items {
				if err = q.parse(e); err != nil {
					break
				}
			}
		case [][]nodes.Node:
			for _, e := range nt {
				if err = q.parse(e); err != nil {
					break
				}
			}
		case nodes.Node:
			switch s := nt.(type) {
			case nodes.ResTarget:
				err = q.parse(s.Val)
			case nodes.SelectStmt:
				err = q.parse(&s)
			case *nodes.SelectStmt:
				if s != nil {
					q.hasSelect = true
					err = q.parse(s.DistinctClause, s.IntoClause, s.TargetList, s.FromClause, s.WhereClause, s.GroupClause, s.
						HavingClause, s.WindowClause, s.ValuesLists, s.SortClause, s.LimitOffset, s.
						LimitCount, s.LockingClause, s.WithClause, s.Op, s.All, s.Larg, s.Rarg)
				}
			case *nodes.WithClause:
				if s != nil {
					q.parse(s.Ctes)
				}
			case nodes.InsertStmt:
				return q.parseError(s, "Insert Statement not allowed.")
			case nodes.UpdateStmt:
				return q.parseError(s, "Update Statement not allowed.")
			case nodes.DeleteStmt:
				return q.parseError(s, "Delete Statement not allowed.")
			case nodes.FuncCall:
				q.hasFunc = true
				var names []string
				for _, n := range s.Funcname.Items {
					names = append(names, n.(nodes.String).Str)
				}

				funcName := strings.Join(names, ".")
				if _, ok := q.funcs[funcName]; !ok {
					f := &QueryFunc{}
					for _, arg := range s.Args.Items {
						switch at := arg.(type) {
						case nodes.TypeCast:
							names = []string{}
							for _, n := range at.TypeName.Names.Items {
								names = append(names, n.(nodes.String).Str)
							}
							f.Args = append(f.Args, strings.Join(names, "."))
						default:
							return q.parseError(at, "Call Function require type cast for arg.")
						}
					}
					q.funcs[funcName] = f
				}
			case nodes.RawStmt:
				err = q.parse(s.Stmt)
			case nodes.A_Expr:
				err = q.parse(s.Lexpr, s.Rexpr)
			case nodes.CaseExpr:
				err = q.parse(s.Xpr, s.Arg, s.Defresult, s.Args)
			case nodes.CaseWhen:
				err = q.parse(s.Xpr, s.Result, s.Expr)
			case nodes.JoinExpr:
				err = q.parse(s.Larg, s.Rarg, s.Quals)
			case nodes.SubLink:
				err = q.parse(s.Xpr, s.OperName, s.Subselect, s.Testexpr)
			case *nodes.IntoClause:
				if s != nil {
					err = q.parse(s.Options, s.ColNames, s.Rel, s.ViewQuery)
				}
			case nodes.RangeSubselect:
				if s.Alias != nil {
					qt := &QueryTable{Fields: map[string]bool{}}
					q.tables[lower(*s.Alias.Aliasname)] = qt
					if sub, ok := s.Subquery.(nodes.SelectStmt); ok {
						for _, target := range sub.TargetList.Items {
							switch tt := target.(type) {
							case nodes.ResTarget:
								if tt.Name == nil {
									if colRef, ok := tt.Val.(nodes.ColumnRef); ok {
										name := colRef.Fields.Items[len(colRef.Fields.Items)-1].(nodes.String).Str
										qt.Fields[lower(name)] = true
									} else {
										return q.parseError(tt, "Exported field is unamed")
									}
								} else {
									qt.Fields[lower(*tt.Name)] = true
								}
							}
						}
					}
				}
				err = q.parse(s.Subquery)
			case nodes.RangeVar:
				if s.Schemaname != nil {
					return q.parseError(s, "Schemamame not is empty")
				}
				if s.Catalogname != nil {
					return q.parseError(s, "Catalogname not is empty")
				}
				qt, ok := q.tables[lower(*s.Relname)]
				if !ok {
					qt = &QueryTable{Name: *s.Relname}
					q.tables[lower(*s.Relname)] = qt
				}
				if qt.Fields == nil {
					qt.Fields = map[string]bool{}
				}
				if s.Alias != nil {
					lqt, ok := q.tables[lower(*s.Alias.Aliasname)]
					if !ok {
						lqt = &QueryTable{Name: *s.Alias.Aliasname, Alias: qt}
						q.tables[lower(*s.Alias.Aliasname)] = lqt
					} else {
						lqt.Alias = qt
					}
					if lqt.Fields == nil {
						lqt.Fields = map[string]bool{}
					}
				}
			case nodes.ColumnRef:
				if len(s.Fields.Items) != 2 {
					return q.parseError(s, "ColumnRef without table or alias name.")
				}

				tbName, colName := s.Fields.Items[0].(nodes.String), s.Fields.Items[1].(nodes.String)
				qt, ok := q.tables[lower(tbName.Str)]
				if !ok {
					qt = &QueryTable{Name: tbName.Str, Fields: map[string]bool{}}
					q.tables[lower(tbName.Str)] = qt
				}
				qt.Fields[lower(colName.Str)] = true
			case nodes.A_Const:
				err = q.parse(s.Val)
			}
		}
		if err != nil {
			return
		}
	}
	return
}

func (q *Query) Parse() error {
	tree, err := pg_query.Parse(q.Source)
	if err != nil {
		return errwrap.Wrap(err, "Parse")
	}
	q.tables = map[string]*QueryTable{}
	q.funcs = map[string]*QueryFunc{}
	for _, stmt := range tree.Statements {
		err = q.parse(stmt)
		if err != nil {
			return err
		}
	}

	var aliases []string

	real := func(qt *QueryTable) *QueryTable {
		for qt.Alias != nil {
			qt = qt.Alias
		}
		return qt
	}

	for name, qt := range q.tables {
		if qt.Alias != nil {
			rqt := real(qt)
			for fname := range qt.Fields {
				rqt.Fields[fname] = true
			}
			aliases = append(aliases, name)
		}
		if qt.Name == "" {
			aliases = append(aliases, name)
		}
	}

	for _, alias := range aliases {
		delete(q.tables, alias)
	}

	for name, qt := range q.tables {
		if ct, ok := q.Context.Tables[name]; ok {
			for cname := range qt.Fields {
				if _, ok = ct.ByName[cname]; !ok {
					return errwrap.Wrap(ErrColumnNotFound, "Column %q.%q", name, cname)
				}
			}
		} else {
			return errwrap.Wrap(ErrTableNotFound, "Table %q", name)
		}
	}
	return err
}
