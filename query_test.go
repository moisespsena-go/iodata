package iodata

import "testing"

func TestQuery_Parse(t *testing.T) {
	ctx := &QueryContext{}
	var sql string
	 sql = `SELECT (1 + (self.nome * 3) + self.j + (case when self.R = 2 then 7 else 8 end)) as t FROM self t
JOIN a ON a.id = t.ID
WHERE exists (select b.colname as f, (select u.ucol from u) as g, NULLIF(b.jjj, '(none)') as h from (select l.w from l) b)
UNION ALL (select 2 from h)`

//	sql = "select c.g from (select b.u as X, 1 as d from b) c"
	q := &Query{Context: ctx, Source:sql}
	if err := q.Parse(); err != nil {
		t.Error(err)
	}
}
