package main

import (
	"database/sql"
	"flag"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/howeyc/gopass"
)

func main() {
	var host, user, database string
	var port int
	flag.StringVar(&host, "h", "localhost", "MySQL host")
	flag.IntVar(&port, "P", 3306, "MySQL port")
	flag.StringVar(&user, "u", "root", "MySQL user")
	flag.StringVar(&database, "D", "", "MySQL database name")
	flag.Parse()

	fmt.Print("Enter password:")
	password, _ := gopass.GetPasswd()

	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?interpolateParams=true&charset=utf8&timeout=1s",
		user, password, host, port, database)
	var conn *sql.DB
	var err error
	if conn, err = sql.Open("mysql", connStr); err != nil {
		panic(err)
	}
	newAnalyzer(conn, database).analyseAndCompress()
}
