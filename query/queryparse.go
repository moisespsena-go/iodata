package query

import (
	"reflect"

	"errors"

	"strings"

	"fmt"

	"github.com/lfittl/pg_query_go"
	nodes "github.com/lfittl/pg_query_go/nodes"
	"github.com/moisespsena-go/error-wrap"
)

var (
	ErrTableNotFound  = errors.New("Table does not exists")
	ErrColumnNotFound = errors.New("Column does not exists")
)

type SQLError struct {
	Position int
	SQL      string
	Message  string
}

func (se *SQLError) Error() string {
	pos := se.Position
	lines := strings.Split(se.SQL, "\n")
	if len(lines) > 1 {
		i := 0
		for j, v := range lines {
			if (i + len(v) + 1) >= pos {
				pos = pos - i
				lines = append(lines[0:j+1], strings.Repeat(" ", pos))
				break
			}
			i += len(v) + 1
		}
	} else {
		lines = append(lines, strings.Repeat(" ", pos))
	}

	l := len(lines) - 1
	poss := fmt.Sprint(pos)
	if l > 0 {
		lines[l-2] = fmt.Sprint(l-2+1, strings.Repeat(" ", len(poss)+1), "| ") + lines[l-2]
	}
	d := 20
	left, _ := pos-d, pos+d
	if left > 0 {
		lines[l-1] = "... " + lines[l-1][left:]
		lines[l] = "    " + lines[l][left:]
	}
	if l > 0 {
		if (d * 2) < len(lines[l-2]) {
			if left > 0 {
				lines[l-2] = lines[l-2][0:d*2+left] + " ..."
			}
		}
	}
	if (d * 2) < len(lines[l-1]) {
		lines[l-1] = lines[l-1][0:d*2] + " ..."
	}
	lineNum := fmt.Sprint(l-1+1, ":"+poss+"| ")
	lines[l-1] = lineNum + lines[l-1]
	lines[l] = strings.Repeat(" ", len(lineNum)) + lines[l] + "^"
	lines = append(lines, "CAUSE: "+se.Message)
	data := strings.Join(lines, "\n")
	println(data)
	return data
}

func (q *Query) parseError(node interface{}, err interface{}) error {
	var e error
	if es, ok := err.(string); ok {
		e = errors.New(es)
	} else {
		e = err.(error)
	}

	f := reflect.ValueOf(node).FieldByName("Location")
	if f.IsValid() {
		return &SQLError{int(f.Int()), q.Source, e.Error()}
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
		case []nodes.Node:
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
					if s.IntoClause != nil {
						return q.parseError(s, "Into Clause not implemented.")
					}
					if len(s.WindowClause.Items) > 0 {
						return q.parseError(s, "Window Clause not implemented.")
					}
					if len(s.LockingClause.Items) > 0 {
						return q.parseError(s, "Locking Clause not implemented.")
					}
					err = q.parse(s.DistinctClause, s.TargetList, s.FromClause, s.WhereClause, s.GroupClause, s.
						HavingClause, s.ValuesLists, s.SortClause, s.LimitOffset, s.LimitCount, s.WithClause, s.Larg, s.Rarg)
				}
			case *nodes.WithClause:
				if s != nil {
					err = q.parse(s.Ctes.Items)
				}
			case nodes.CommonTableExpr:
				qt := &QueryTable{Fields: map[string]bool{}}
				q.tables[lower(*s.Ctename)] = qt
				err = q.parse(s.Ctequery)
			case nodes.InsertStmt:
				return q.parseError(s, "Insert Statement not allowed.")
			case nodes.UpdateStmt:
				return q.parseError(s, "Update Statement not allowed.")
			case nodes.DeleteStmt:
				return q.parseError(s, "Delete Statement not allowed.")
			case nodes.FuncCall:
				q.nodes = append(q.nodes, &s)
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
				q.nodes = append(q.nodes, &s)
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
					return q.parseError(s, "ColumnRef without table or alias name. Use: TABLE_OR_ALIAS.COLUMN_NAME. Example: `users.name`.")
				}
				q.nodes = append(q.nodes, &s)
				tbName, colName := s.Fields.Items[0].(nodes.String), s.Fields.Items[1].(nodes.String)
				qt, ok := q.tables[lower(tbName.Str)]
				if !ok {
					qt = &QueryTable{Name: tbName.Str, Fields: map[string]bool{}}
					q.tables[lower(tbName.Str)] = qt
				}
				qt.Fields[lower(colName.Str)] = true
			case nodes.A_Const:
				err = q.parse(s.Val)
			case nodes.BoolExpr:
				err = q.parse(s.Args, s.Xpr)
			}
		}
		if err != nil {
			return
		}
	}
	return
}

func (q *Query) Parse() (err error) {
	q.tree, err = pg_query.Parse(q.Source)
	if err != nil {
		return errwrap.Wrap(err, "Parse")
	}
	q.tables = map[string]*QueryTable{}
	q.funcs = map[string]*QueryFunc{}
	for _, stmt := range q.tree.Statements {
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
	return err
}

func (q *Query) Validate() error {
	for name, qt := range q.tables {
		if ct, ok := q.Context.DataSources[name]; ok {
			h := ct.Header()
			for cname := range qt.Fields {
				if !h.HasField(cname) {
					return errwrap.Wrap(ErrColumnNotFound, "Column %q.%q", name, cname)
				}
			}
		} else {
			return errwrap.Wrap(ErrTableNotFound, "Table %q", name)
		}
	}
	return nil
}
