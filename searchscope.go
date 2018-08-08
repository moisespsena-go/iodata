package iodata

import (
	"reflect"

	"fmt"
	"strings"

	"bytes"
	"database/sql"
	"database/sql/driver"
	"regexp"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/moisespsena-go/iodata/api"
	"github.com/moisespsena/go-error-wrap"
)

type SearchField struct {
	Index int
	Value reflect.Value
}

type Scope struct {
	Header          api.DataHeader
	Search          *Search
	Dialect         gorm.Dialect
	Value           interface{}
	SQL             string
	SQLVars         []interface{}
	instanceID      string
	primaryKeyField int
	skipLeft        bool
	fields          *[]*SearchField
	fieldsByName    map[string]*SearchField
	selectAttrs     *[]string
	primaryValue    interface{}
	primaryValueSet bool
	counter         bool
	Table           string
	KeyField        string
}

// Quote used to quote string to escape them for database
func (scope *Scope) Quote(str string) string {
	if strings.Index(str, ".") != -1 {
		newStrs := []string{}
		for _, str := range strings.Split(str, ".") {
			newStrs = append(newStrs, scope.Dialect.Quote(str))
		}
		return strings.Join(newStrs, ".")
	}

	return scope.Dialect.Quote(str)
}

// AddToVars add value as sql's vars, used to prevent SQL injection
func (scope *Scope) AddToVars(value interface{}) string {
	if expr, ok := value.(*expr); ok {
		exp := expr.expr
		for _, arg := range expr.args {
			exp = strings.Replace(exp, "?", scope.AddToVars(arg), 1)
		}
		return exp
	}

	scope.SQLVars = append(scope.SQLVars, value)
	return scope.Dialect.BindVar(len(scope.SQLVars))
}

// SelectAttrs return selected attributes
func (scope *Scope) SelectAttrs() []string {
	if scope.selectAttrs == nil {
		attrs := []string{}
		for _, value := range scope.Search.selects {
			if str, ok := value.(string); ok {
				attrs = append(attrs, str)
			} else if strs, ok := value.([]string); ok {
				attrs = append(attrs, strs...)
			} else if strs, ok := value.([]interface{}); ok {
				for _, str := range strs {
					attrs = append(attrs, fmt.Sprintf("%v", str))
				}
			}
		}
		scope.selectAttrs = &attrs
	}
	return *scope.selectAttrs
}

// CombinedConditionSql return combined condition sql
func (scope *Scope) CombinedConditionSql() (sql string, err error) {
	var (
		joinSQL, whereSQL, havingSQL string
	)
	if joinSQL, err = scope.joinsSQL(); err != nil {
		err = errwrap.Wrap(err, "Join SQL")
		return
	}
	if whereSQL, err = scope.whereSQL(); err != nil {
		err = errwrap.Wrap(err, "Where SQL")
		return
	}
	if havingSQL, err = scope.havingSQL(); err != nil {
		err = errwrap.Wrap(err, "Having SQL")
		return
	}
	if scope.Search.raw {
		whereSQL = strings.TrimSuffix(strings.TrimPrefix(whereSQL, "WHERE ("), ")")
	}
	return joinSQL + whereSQL + scope.groupSQL() +
		havingSQL + scope.orderSQL() + scope.limitAndOffsetSQL(), nil
}

// Raw set raw sql
func (scope *Scope) Raw(sql string) *Scope {
	scope.SQL = strings.Replace(sql, "$$$", "?", -1)
	return scope
}

var (
	columnRegexp        = regexp.MustCompile("^[a-zA-Z\\d]+(\\.[a-zA-Z\\d]+)*$") // only match string like `name`, `users.name`
	isNumberRegexp      = regexp.MustCompile("^\\s*\\d+\\s*$")                   // match if string is number
	comparisonRegexp    = regexp.MustCompile("(?i) (=|<>|(>|<)(=?)|LIKE|IS|IN) ")
	countingQueryRegexp = regexp.MustCompile("(?i)^count(.+)$")
)

func (scope *Scope) quoteIfPossible(str string) string {
	if columnRegexp.MatchString(str) {
		return scope.Quote(str)
	}
	return str
}

func (scope *Scope) primaryCondition(value interface{}) string {
	return fmt.Sprintf("(%v.%v = %v)", scope.Table, scope.Quote(scope.KeyField), value)
}

func (scope *Scope) buildCondition(clause map[string]interface{}, include bool) (str string, err error) {
	var (
		quotedTableName  = scope.Table
		quotedPrimaryKey = scope.Quote(scope.KeyField)
		equalSQL         = "="
		inSQL            = "IN"
	)

	// If building not conditions
	if !include {
		equalSQL = "<>"
		inSQL = "NOT IN"
	}

	switch value := clause["query"].(type) {
	case sql.NullInt64:
		str = fmt.Sprintf("(%v.%v %s %v)", quotedTableName, quotedPrimaryKey, equalSQL, value.Int64)
		return
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		str = fmt.Sprintf("(%v.%v %s %v)", quotedTableName, quotedPrimaryKey, equalSQL, value)
		return
	case []int, []int8, []int16, []int32, []int64, []uint, []uint8, []uint16, []uint32, []uint64, []string, []interface{}:
		if !include && reflect.ValueOf(value).Len() == 0 {
			return
		}
		str = fmt.Sprintf("(%v.%v %s (?))", quotedTableName, quotedPrimaryKey, inSQL)
		clause["args"] = []interface{}{value}
	case string:
		if isNumberRegexp.MatchString(value) {
			str = fmt.Sprintf("(%v.%v %s %v)", quotedTableName, quotedPrimaryKey, equalSQL, scope.AddToVars(value))
			return
		}

		if value != "" {
			if !include {
				if comparisonRegexp.MatchString(value) {
					str = fmt.Sprintf("NOT (%v)", value)
				} else {
					str = fmt.Sprintf("(%v.%v NOT IN (?))", quotedTableName, scope.Quote(value))
				}
			} else {
				str = fmt.Sprintf("(%v)", value)
			}
		}
	case map[string]interface{}:
		var sqls []string
		for key, value := range value {
			if value != nil {
				sqls = append(sqls, fmt.Sprintf("(%v.%v %s %v)", quotedTableName, scope.Quote(key), equalSQL, scope.AddToVars(value)))
			} else {
				if !include {
					sqls = append(sqls, fmt.Sprintf("(%v.%v IS NOT NULL)", quotedTableName, scope.Quote(key)))
				} else {
					sqls = append(sqls, fmt.Sprintf("(%v.%v IS NULL)", quotedTableName, scope.Quote(key)))
				}
			}
		}
		str = strings.Join(sqls, " AND ")
		return
	default:
		err = fmt.Errorf("invalid query condition: %v", value)
		return
	}

	replacements := []string{}
	args := clause["args"].([]interface{})
	for _, arg := range args {
		switch reflect.ValueOf(arg).Kind() {
		case reflect.Slice: // For where("id in (?)", []int64{1,2})
			if scanner, ok := interface{}(arg).(driver.Valuer); ok {
				arg, err = scanner.Value()
				replacements = append(replacements, scope.AddToVars(arg))
			} else if b, ok := arg.([]byte); ok {
				replacements = append(replacements, scope.AddToVars(b))
			} else if as, ok := arg.([][]interface{}); ok {
				var tempMarks []string
				for _, a := range as {
					var arrayMarks []string
					for _, v := range a {
						arrayMarks = append(arrayMarks, scope.AddToVars(v))
					}

					if len(arrayMarks) > 0 {
						tempMarks = append(tempMarks, fmt.Sprintf("(%v)", strings.Join(arrayMarks, ",")))
					}
				}

				if len(tempMarks) > 0 {
					replacements = append(replacements, strings.Join(tempMarks, ","))
				}
			} else if values := reflect.ValueOf(arg); values.Len() > 0 {
				var tempMarks []string
				for i := 0; i < values.Len(); i++ {
					tempMarks = append(tempMarks, scope.AddToVars(values.Index(i).Interface()))
				}
				replacements = append(replacements, strings.Join(tempMarks, ","))
			} else {
				replacements = append(replacements, scope.AddToVars(Expr("NULL")))
			}
		default:
			if valuer, ok := interface{}(arg).(driver.Valuer); ok {
				arg, err = valuer.Value()
			}

			replacements = append(replacements, scope.AddToVars(arg))
		}
	}

	buff := bytes.NewBuffer([]byte{})
	i := 0
	for _, s := range str {
		if s == '?' {
			buff.WriteString(replacements[i])
			i++
		} else {
			buff.WriteRune(s)
		}
	}

	str = buff.String()

	return
}

func (scope *Scope) buildSelectQuery(clause map[string]interface{}) (str string) {
	switch value := clause["query"].(type) {
	case string:
		str = value
	case []string:
		str = strings.Join(value, ", ")
	}

	args := clause["args"].([]interface{})
	replacements := []string{}
	for _, arg := range args {
		switch reflect.ValueOf(arg).Kind() {
		case reflect.Slice:
			values := reflect.ValueOf(arg)
			var tempMarks []string
			for i := 0; i < values.Len(); i++ {
				tempMarks = append(tempMarks, scope.AddToVars(values.Index(i).Interface()))
			}
			replacements = append(replacements, strings.Join(tempMarks, ","))
		default:
			if valuer, ok := interface{}(arg).(driver.Valuer); ok {
				arg, _ = valuer.Value()
			}
			replacements = append(replacements, scope.AddToVars(arg))
		}
	}

	buff := bytes.NewBuffer([]byte{})
	i := 0
	for pos := range str {
		if str[pos] == '?' {
			buff.WriteString(replacements[i])
			i++
		} else {
			buff.WriteByte(str[pos])
		}
	}

	str = buff.String()

	return
}

func (scope *Scope) whereSQL() (sql string, err error) {
	var (
		quotedTableName                                = scope.Table
		primaryConditions, andConditions, orConditions []string
	)

	if scope.primaryValueSet {
		sql = fmt.Sprintf("%v.%v = %v", quotedTableName, scope.Quote(scope.KeyField), scope.AddToVars(scope.primaryValue))
		primaryConditions = append(primaryConditions, sql)
	}

	for i, clause := range scope.Search.whereConditions {
		if sql, err = scope.buildCondition(clause, true); err != nil {
			return "", errwrap.Wrap(err, "Build Where Condition %d", i)
		}
		if sql != "" {
			andConditions = append(andConditions, sql)
		}
	}

	for i, clause := range scope.Search.orConditions {
		if sql, err = scope.buildCondition(clause, true); err != nil {
			return "", errwrap.Wrap(err, "Build Or Condition %d", i)
		} else if sql != "" {
			orConditions = append(orConditions, sql)
		}
	}

	for i, clause := range scope.Search.notConditions {
		if sql, err = scope.buildCondition(clause, false); err != nil {
			return "", errwrap.Wrap(err, "Build Not Condition %d", i)
		} else if sql != "" {
			andConditions = append(andConditions, sql)
		}
	}

	sql = ""
	orSQL := strings.Join(orConditions, " OR ")
	combinedSQL := strings.Join(andConditions, " AND ")
	if len(combinedSQL) > 0 {
		if len(orSQL) > 0 {
			combinedSQL = combinedSQL + " OR " + orSQL
		}
	} else {
		combinedSQL = orSQL
	}

	if len(primaryConditions) > 0 {
		sql = "WHERE " + strings.Join(primaryConditions, " AND ")
		if len(combinedSQL) > 0 {
			sql = sql + " AND (" + combinedSQL + ")"
		}
	} else if len(combinedSQL) > 0 {
		sql = "WHERE " + combinedSQL
	}
	return
}

func (scope *Scope) selectSQL() (sql string) {
	if len(scope.Search.selects) == 0 {
		names := scope.Header.Names()
		columns := make([]string, len(names))
		for i, name := range names {
			columns[i] = fmt.Sprintf("%v.%v", scope.Table, name)
		}
		sql = strings.Join(columns, ", ")
	} else {
		sql = scope.buildSelectQuery(scope.Search.selects)
	}
	return
}

func (scope *Scope) orderSQL() string {
	if len(scope.Search.orders) == 0 || scope.Search.ignoreOrderQuery {
		return ""
	}

	var orders []string
	for _, order := range scope.Search.orders {
		if str, ok := order.(string); ok {
			orders = append(orders, scope.quoteIfPossible(str))
		} else if expr, ok := order.(*expr); ok {
			exp := expr.expr
			for _, arg := range expr.args {
				exp = strings.Replace(exp, "?", scope.AddToVars(arg), 1)
			}
			orders = append(orders, exp)
		}
	}
	return " ORDER BY " + strings.Join(orders, ",")
}

func (scope *Scope) limitAndOffsetSQL() string {
	return scope.Dialect.LimitAndOffsetSQL(scope.Search.limit, scope.Search.offset)
}

func (scope *Scope) groupSQL() string {
	if len(scope.Search.group) == 0 {
		return ""
	}
	return " GROUP BY " + scope.Search.group
}

func (scope *Scope) havingSQL() (sql string, err error) {
	if len(scope.Search.havingConditions) == 0 {
		return "", nil
	}

	var andConditions []string
	for i, clause := range scope.Search.havingConditions {
		if sql, err = scope.buildCondition(clause, true); err != nil {
			return "", errwrap.Wrap(err, "Build Having Condition %d", i)
		} else if sql != "" {
			andConditions = append(andConditions, sql)
		}
	}

	combinedSQL := strings.Join(andConditions, " AND ")
	if len(combinedSQL) == 0 {
		return "", nil
	}

	return " HAVING " + combinedSQL, nil
}

func (scope *Scope) joinsSQL() (sql string, err error) {
	var joinConditions []string
	for i, clause := range scope.Search.joinConditions {
		if sql, err = scope.buildCondition(clause, true); err != nil {
			return "", errwrap.Wrap(err, "Build Join Condition %d", i)
		} else if sql != "" {
			joinConditions = append(joinConditions, strings.TrimSuffix(strings.TrimPrefix(sql, "("), ")"))
		}
	}

	return strings.Join(joinConditions, " ") + " ", nil
}

func (scope Scope) PrepareQuerySQL() (*Scope, error) {
	scope.Search = &(*scope.Search)
	sql, err := scope.CombinedConditionSql()
	if err != nil {
		return nil, err
	}
	if scope.Search.raw {
		scope.Raw(sql)
	} else {
		scope.Raw(fmt.Sprintf("SELECT %v FROM %v %v", scope.selectSQL(), scope.Table, sql))
	}
	return &scope, nil
}

func (scope *Scope) inlineCondition(values ...interface{}) *Scope {
	if len(values) > 0 {
		scope.Search.Where(values[0], values[1:]...)
	}
	return scope
}

func (scope Scope) Count(value interface{}) (*Scope, error) {
	scope.Search = &(*scope.Search)
	scope.counter = true
	if query, ok := scope.Search.selects["query"]; !ok || !countingQueryRegexp.MatchString(fmt.Sprint(query)) {
		if len(scope.Search.group) != 0 {
			scope.Search.Select("count(*) FROM ( SELECT count(*) as name ")
			scope.Search.group += " ) AS count_table"
		} else {
			scope.Search.Select("count(*)")
		}
	}
	scope.Search.ignoreOrderQuery = true
	return &scope, nil
}

func (scope Scope) Counter() bool {
	return scope.counter
}
