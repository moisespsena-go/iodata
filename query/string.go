package query

import (
	"fmt"

	"strings"

	"io"

	"strconv"

	"github.com/lfittl/pg_query_go"
	nodes "github.com/lfittl/pg_query_go/nodes"
)

func isOuterJoin(t nodes.JoinType) bool {
	return (((1 << (t)) &
		((1 << nodes.JOIN_LEFT) |
			(1 << nodes.JOIN_FULL) |
			(1 << nodes.JOIN_RIGHT) |
			(1 << nodes.JOIN_ANTI))) != 0)
}

type QueryString struct {
	Tree  pg_query.ParsetreeList
	Out   io.Writer
	size  int
	error error
	q     int
}

func (q *QueryString) where(node nodes.Node) {
	if node != nil {
		q.twp("\n", "WHERE\n")
		q.tab(func() {
			q.tw()
			q.gen(node)
		})
	}
}
func (q *QueryString) gen(nods ...interface{}) {
	for i, node := range nods {
		if node == nil {
			continue
		}
		switch nt := node.(type) {
		case nodes.List:
			for _, e := range nt.Items {
				q.gen(e)
			}
		case []nodes.Node:
			for _, e := range nt {
				q.gen(e)
			}
		case [][]nodes.Node:
			for _, e := range nt {
				q.gen(e)
			}
		case nodes.Node:
			switch s := nt.(type) {
			case nodes.ResTarget:
				q.gen(s.Val)
				if s.Name != nil {
					q.w(" AS " + *s.Name)
				}
			case *nodes.WithClause:
				if s != nil {
					q.tw("WITH")
					if s.Recursive {
						q.w(" RECURSIVE")
					}
					q.w("\n")
					q.tab(func() {
						for i, item := range s.Ctes.Items {
							cte := item.(nodes.CommonTableExpr)
							if i != 0 {
								q.w(",\n")
							}
							q.tw(*cte.Ctename)
							if len(cte.Aliascolnames.Items) > 0 {
								q.w("(")
								for j, cn := range cte.Aliascolnames.Items {
									if j != 0 {
										q.w(", ")
									}
									q.w(cn.(nodes.String).Str)
								}
								q.w(")")
							}
							q.w(" AS (\n")
							q.tab(func() {
								q.gen(cte.Ctequery)
							})
							q.twp("\n", ")")
						}
					})
					q.w("\n")
				}
			case *nodes.IntoClause:
				if s != nil {
					panic(fmt.Errorf("IntoClause Not Implemented %q", nt))
				}
			case nodes.FuncCall:
				var names []string
				for _, n := range s.Funcname.Items {
					names = append(names, n.(nodes.String).Str)
				}
				q.w(strings.Join(names, "."))
				q.w("(")
				q.join(", ", s.Args.Items...)
				q.w(")")
			case nodes.RawStmt:
				q.gen(s.Stmt)
			case nodes.A_Expr:
				switch s.Kind {
				case nodes.AEXPR_OP:
					q.w("(")
					q.gen(s.Lexpr)
					q.w(" ", q.joinStrings(".", s.Name.Items...), " ")
					q.gen(s.Rexpr)
					q.w(")")
				case nodes.AEXPR_NULLIF:
					q.w("NULLIF(")
					q.gen(s.Lexpr)
					q.w(", ")
					q.gen(s.Rexpr)
					q.w(")")
				default:
					panic("A_Expr.Kind Not Implemented error!")
				}
			case nodes.BoolExpr:
				switch s.Boolop {
				case nodes.AND_EXPR, nodes.OR_EXPR:
					op := "OR"
					if s.Boolop == nodes.AND_EXPR {
						op = "AND"
					}
					for i, e := range s.Args.Items {
						if i != 0 {
							q.twp("\n", op, " ")
							q.tab(func() {
								q.gen(e)
							})
						} else {
							q.gen(e)
						}
					}
				default:
					q.w("NOT ")
					q.gen(s.Args.Items[0])
				}
			case nodes.CaseExpr:
				q.gen(s.Xpr, s.Arg, s.Defresult, s.Args)
			case nodes.CaseWhen:
				q.gen(s.Xpr, s.Result, s.Expr)
			case nodes.JoinExpr:
				q.gen(s.Larg)
				q.twp("\n", "")
				if s.IsNatural {
					q.w("NATURAL ")
				}
				switch s.Jointype {
				case nodes.JOIN_INNER:
				case nodes.JOIN_LEFT:
					q.w("LEFT ")
				case nodes.JOIN_RIGHT:
					q.w("RIGHT ")
				case nodes.JOIN_FULL:
					q.w("FULL")
				}

				if isOuterJoin(s.Jointype) {
					q.w("OUTER ")
				}
				q.w("JOIN ")
				switch r := s.Rarg.(type) {
				case nodes.RangeVar:
					q.w(*r.Relname, " ")
					if s.Alias != nil {
						q.w("AS ", *r.Alias.Aliasname, " ")
					}
				}

				q.tab(func() {
					q.twp("\n", "ON ")
					q.gen(s.Quals)
				})
			case nodes.SubLink:
				switch s.SubLinkType {
				case nodes.EXISTS_SUBLINK:
					q.w("EXISTS")
				case nodes.ALL_SUBLINK:
					q.w("ALL")
				case nodes.ANY_SUBLINK:
					q.w("ANY")
				case nodes.ROWCOMPARE_SUBLINK:
				case nodes.ARRAY_SUBLINK:
					q.w("ARRAY")
				case nodes.EXPR_SUBLINK:
				default:
					panic(fmt.Errorf("SubLinkType Not Implemented %q", nt))
				}
				q.w(" (\n")
				q.tab(func() {
					q.gen(s.Xpr, s.OperName, s.Subselect, s.Testexpr)
				})
				q.w(")")
			case nodes.RangeSubselect:
				q.w("(")
				q.gen(s.Subquery)
				q.w(")")
				if s.Alias != nil {
					q.w(" AS ", *s.Alias.Aliasname)
				}
			case nodes.SelectStmt:
				q.gen(&s)
			case *nodes.SelectStmt:
				if s != nil {
					if s.Op == nodes.SETOP_NONE {
						if len(s.ValuesLists) > 0 {
							for i, values := range s.ValuesLists {
								if i != 0 {
									q.w(",")
								}
								q.tw("VALUES (")
								q.join(",", values...)
								q.w(")")
							}
						} else {
							q.gen(s.WithClause)

							q.tw("SELECT")

							q.tab(func() {
								if s.All {
									q.w(" ALL")
								}
								if len(s.DistinctClause.Items) > 0 {
									if s.DistinctClause.Items[0] != nil {
										q.twp("\n", "DISTINCT ON (")
										q.tab(func() {
											for i, n := range s.DistinctClause.Items {
												if i != 0 {
													q.twp(",\n")
												} else {
													q.twp("\n")
												}
												q.gen(n)
											}
										})
										q.twp("\n", ")")
									} else {
										q.w(" DISTINCT")
									}
								}

								for i, n := range s.TargetList.Items {
									if i == 0 {
										q.twp("\n")
									} else {
										q.twp(",\n")
									}
									q.gen(n)
								}
								// TODO: IntoClause
								// TODO: WindowClause
								// TODO: LockingClause

								q.gen(s.FromClause)

								q.where(s.WhereClause)

								if len(s.GroupClause.Items) > 0 {
									q.twp("\n", "GROUP BY ")
									q.gen(s.GroupClause)
								}
								if s.HavingClause != nil {
									q.twp("\n", "HAVING ")
									q.gen(s.HavingClause)
								}

								// having
								// window
								// { UNION | INTERSECT | EXCEPT } [ ALL | DISTINCT ] select

								if len(s.SortClause.Items) > 0 {
									q.twp("\n", "ORDER BY ")
									q.gen(s.SortClause.Items)
								}

								if s.LimitCount != nil {
									q.twp("\n", "LIMIT ")
									q.gen(s.LimitCount)

									if s.LimitOffset != nil {
										q.w(" OFFSET ")
										q.gen(s.LimitOffset)
									}
								} else if s.LimitOffset != nil {
									q.twp("\n", "OFFSET ")
									q.gen(s.LimitOffset)
								}

								// FETCh
								// FOR
							})
						}
					} else {
						q.tab(func() {
							q.gen(s.Larg)
						})
						q.w("\n")
						switch s.Op {
						case nodes.SETOP_UNION:
							q.tw("UNION")
						case nodes.SETOP_INTERSECT:
							q.tw("INTERSECT")
						case nodes.SETOP_EXCEPT:
							q.tw("EXCEPT")
						}
						if s.All {
							q.w(" ALL")
						}
						q.w("\n")
						q.tab(func() {
							q.gen(s.Rarg)
						})
					}
				}
			case nodes.RangeVar:
				q.twp("\n", "FROM ", *s.Relname)
				if s.Alias != nil {
					q.w(" AS ", *s.Alias.Aliasname)
				}
			case nodes.ColumnRef:
				tbName, colName := s.Fields.Items[0].(nodes.String), s.Fields.Items[1].(nodes.String)
				q.w(tbName.Str, ".", colName.Str)
			case nodes.A_Const:
				switch vt := s.Val.(type) {
				case nodes.Integer:
					q.w(strconv.Itoa(int(vt.Ival)))
				case nodes.Float:
					q.w(vt.Str)
				case nodes.String:
					q.w("'", strings.Replace(vt.Str, "'", "''", -1), "'")
				default:
					panic(fmt.Errorf("A_Const Not Implemented %q", nt))
				}
			case nodes.String:
				q.tw(s.Str)
			case *nodes.TypeName:
				println()
			case nodes.TypeCast:
				var typName string

				if s.TypeName != nil {
					if s.TypeName.Names.Items[0].(nodes.String).Str == "pg_catalog" {
						typName = s.TypeName.Names.Items[1].(nodes.String).Str
					} else {
						typName = q.joinStrings(".", s.TypeName.Names.Items...)
					}
				}
				var skip bool
				if typName == "bool" {
					if v, ok := s.Arg.(nodes.A_Const); ok {
						if v, ok := v.Val.(nodes.String); ok {
							if v.Str == "t" {
								q.w("TRUE")
								skip = true
							} else if v.Str == "f" {
								q.w("FALSE")
								skip = true
							}
						}
					}
				}
				if !skip {
					q.w("(")
					q.gen(s.Arg)
					q.w(")::")
					q.w(typName)
				}
			case nodes.ParamRef:
				q.w("$",fmt.Sprint(s.Number))
			default:
				panic(fmt.Errorf("Not Implemented %q", nt))
			}
		default:
			panic(fmt.Errorf("Not Implemented %d -> %q", i, nt))
		}
	}
}

func (q *QueryString) w(v ...string) {
	s, err := q.Out.Write([]byte(strings.Join(v, "")))
	q.size += s
	if err != nil {
		panic(err)
	}
}

func (q *QueryString) tw(v ...string) {
	q.w(strings.Repeat("  ", q.q), strings.Join(v, ""))
}

func (q *QueryString) twp(prefix string, v ...string) {
	q.w(prefix, strings.Repeat("  ", q.q), strings.Join(v, ""))
}

func (q *QueryString) tab(f func()) {
	q.q++
	defer func() {
		q.q--
	}()
	f()
}

func (q *QueryString) Generate() {
	q.join("\n\n", q.Tree.Statements...)
}

func (q *QueryString) Compile(out io.Writer) {
	defer func() {
		q.Out = nil
		q.size = 0
	}()
	q.Out = out
	q.Generate()
}

func (q *QueryString) joinStrings(sep string, args ...nodes.Node) string {
	values := make([]string, len(args))
	for i, arg := range args {
		values[i] = arg.(nodes.String).Str
	}
	return strings.Join(values, sep)
}

func (q *QueryString) join(sep string, args ...nodes.Node) {
	l := len(args)
	if l == 0 {
		return
	}
	last, args := args[l-1], args[0:l-1]
	for _, arg := range args {
		q.gen(arg)
		q.w(sep)
	}
	q.gen(last)
}

func (q *QueryString) joint(sep string, args ...nodes.Node) {

}
