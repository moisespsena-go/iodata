package main

import (
	"os"

	"github.com/moisespsena-go/iodata/query"
)

func main() {
	ctx := &query.QueryContext{}

	var sql string
	sql = `
SELECT (1 + (self.nome * 3) + self.j) as t 
FROM self t
JOIN a ON a.id = t.ID
	WHERE exists (select b.colname as f, (select u.ucol from u) as g, NULLIF(b.jjj, '(none)') as h from (select l.w from l) b)
	UNION ALL (select 2 from h)`
	q := &query.Query{Context: ctx, Source: sql}
	if err := q.Parse(); err != nil {
		panic(err)
	}
	os.Stderr.WriteString("-----------------------\n")
	c := q.NewCompiler()
	c.Compile(os.Stderr)
	os.Stderr.WriteString("\n-----------------------\n")
	os.Stderr.Sync()
}
