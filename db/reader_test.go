package db

import (
	"testing"

	"database/sql"
	"log"
	"math/rand"
	"time"

	"fmt"

	_ "github.com/lib/pq"
	"github.com/moisespsena-go/iodata"
	"github.com/moisespsena-go/iodata/api/datatypes"
)

func TestReader_Read(t *testing.T) {
	db := initDB(t)
	defer db.Close()
	fields := []iodata.Field{{"a", datatypes.INT64}, {"v", datatypes.FLOAT64}, {"d", datatypes.DATE}}
	r := &Reader{DB: db, DataHeader: iodata.NewDataHeader(fields), SQL: "select * from teste order by a"}
	v := struct {
		a int64
		b float64
		c time.Time
	}{}
	for i := 0; i < 10; i++ {
		_, err := r.ReadOne(&v.a, &v.b, &v.c)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(i, ".", v)
	}
}

func initDB(t *testing.T) *sql.DB {
	db, err := sql.Open("postgres", "password=123456")
	if err != nil {
		t.Errorf("Open Connection: %v", err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS teste (a SERIAL, b float, c date)")
	if err != nil {
		t.Errorf("Create Table: %v", err)
	}
	rows, err := db.Query("SELECT count(*) FROM teste")
	if err != nil {
		t.Errorf("Count: %v", err)
	}
	rows.Next()
	var count int
	if err = rows.Scan(&count); err != nil {
		t.Errorf("Count scan: %v", err)
	}
	if count == 0 {
		func() {
			tx, err := db.Begin()
			if err != nil {
				log.Fatal(err)
			}
			defer func() {
				if err != nil {
					tx.Rollback()
				} else {
					tx.Commit()
				}
			}()
			var stmt *sql.Stmt
			stmt, err = tx.Prepare("INSERT INTO teste (b, c) VALUES ($1, $2)")
			if err != nil {
				t.Fatal(err)
			}
			defer stmt.Close()
			for i := 0; i < 10; i++ {
				_, err = stmt.Exec(rand.Float64(), time.Now())
				if err != nil {
					t.Errorf("Insert[%d]: %v", i, err)
				}
				<-time.After(1 * time.Second)
			}
		}()
	}
	return db
}
