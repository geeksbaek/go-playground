package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/elastic/apm-agent-go/module/apmhttp"
	"github.com/elastic/apm-agent-go/module/apmsql"
	apmsqlite3 "github.com/elastic/apm-agent-go/module/apmsql/sqlite3"
	sqlite3 "github.com/mattn/go-sqlite3"
)

func sleep(n int64) int64 {
	time.Sleep(time.Duration(n))
	return n
}
func main() {
	apmsql.Register("sqlite3_custom", &sqlite3.SQLiteDriver{
		ConnectHook: func(conn *sqlite3.SQLiteConn) error {
			return conn.RegisterFunc("sleep", sleep, true)
		},
	}, apmsql.WithDSNParser(apmsqlite3.ParseDSN))
	db, err := apmsql.Open("sqlite3_custom", ":memory:")
	if err != nil {
		panic(err)
	}
	http.HandleFunc("/sleep", func(w http.ResponseWriter, req *http.Request) {
		delay := time.Duration(rand.ExpFloat64() * float64(time.Millisecond*50))
		if _, err := db.ExecContext(req.Context(), "SELECT sleep(?)", int64(delay)); err != nil {
			panic(err)
		}
		fmt.Fprintln(w, "slept for", delay)
	})
	http.ListenAndServe(":8080", apmhttp.Wrap(http.DefaultServeMux))
}
