package db

import (
	"testing"

	"database/sql"
	"log"
	"math/rand"
	"time"

	"fmt"

	"io"

	_ "github.com/lib/pq"
	"github.com/moisespsena-go/iodata"
	"github.com/moisespsena-go/iodata/api/datatypes"
	"github.com/moisespsena-go/iodata/modelstruct"
)

func TestReader_Read(t *testing.T) {
	testDB(t, func(db *sql.DB) {
		testReader(db, 10, func(r *Reader) {
			v := struct {
				a int64
				b float64
				c time.Time
			}{}
			old := v
			for i := 0; i < 10; i++ {
				_, err := r.ReadOne(&v.a, &v.b, &v.c)
				if err != nil {
					t.Fatal(err)
				}
				if v.a == 0 || v.b == 0.0 || v.c.IsZero() {
					t.Fatal("Result %d: field is empty", i+1)
				}
				if i > 0 && fmt.Sprint(old) == fmt.Sprint(v) {
					t.Fatal("Result %d = Result %d", i, i+1)
				}
				old = v
			}
		})
	})
}

func TestReader_ReadModelStruct(t *testing.T) {
	testDB(t, func(db *sql.DB) {
		testReader(db, 10, func(r *Reader) {
			model := modelstruct.Get(&User{})
			results := make([]interface{}, 10)
			count, err := r.ReadModelStruct(model, results...)
			assertCountError(t, 10, count, err)
			assertError(t, err)
			for _, res := range results {
				assertUser(t, res.(*User))
			}
		})
	})
}

func TestReader_ReadModelStruct2(t *testing.T) {
	testDB(t, func(db *sql.DB) {
		testReader(db, 1, func(r *Reader) {
			var a, b User
			m := modelstruct.Get(&a)
			count, err := r.ReadModelStruct(m, &a, &b)
			if err != io.EOF {
				t.Fatal("Expected EOF")
			} else {
				err = nil
			}
			assertCountError(t, 1, count, nil)
			assertError(t, err)
			assertUser(t, &a)
			assertUserEmpty(t, &b)
		})
	})
}

func TestReader_ReadStruct(t *testing.T) {
	testDB(t, func(db *sql.DB) {
		testReader(db, 10, func(r *Reader) {
			var user User
			count, err := r.ReadStruct(&user)
			assertCountError(t, 1, count, err)
			assertError(t, err)
			assertUser(t, &user)
		}, func(r *Reader) {
			var user1, user2 User
			count, err := r.ReadStruct(&user1, &user2)
			assertCountError(t, 2, count, err)
			assertError(t, err)
			assertUser(t, &user1, &user2)
		}, func(r *Reader) {
			var users = make([]*User, 2)
			count, err := r.ReadStruct(users)
			assertCountError(t, 2, count, err)
			assertError(t, err)
			assertUser(t, users...)
		})
	})
}

type User struct {
	B float64
	C time.Time
}

func assertCountError(t *testing.T, expected, got int, err error) {
	if (err == nil && expected != got) || err == io.EOF {
		t.Fatalf("Expected %d, but got %d.", expected, got)
	}
}

func assertError(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func assertUser(t *testing.T, users ...*User) {
	for _, user := range users {
		if user.B == 0.0 {
			t.Fatal("user.B is zero")
		}
		if user.C.IsZero() {
			t.Fatal("user.C is zero")
		}
	}
}

func assertUserEmpty(t *testing.T, users ...*User) {
	for _, user := range users {
		if user.B != 0.0 {
			t.Fatal("user.B is not zero")
		}
		if !user.C.IsZero() {
			t.Fatal("user.C is not zero")
		}
	}
}

func testReader(db *sql.DB, limit int, funcs ...func(r *Reader)) {
	fields := []iodata.Field{{"a", datatypes.INT64}, {"v", datatypes.FLOAT64}, {"d", datatypes.DATE}}
	r := &Reader{DB: db, DataHeader: iodata.NewDataHeader(fields), SQL: fmt.Sprintf("select * from teste order by a limit %d", limit)}

	for _, f := range funcs {
		f(r)
	}
}

func testDB(t *testing.T, funcs ...func(db *sql.DB)) {
	db := initDB(t)
	defer db.Close()
	for _, f := range funcs {
		f(db)
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
