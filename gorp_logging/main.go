package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/gorp.v1"
)

type EC2 struct {
	Id   int64
	Type string
}

type MyGorpTracer struct{}

func (t *MyGorpTracer) Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func main() {
	db, err := sql.Open("sqlite3", "./ec2.db")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
	tracer := &MyGorpTracer{}
	dbmap.TraceOn("[gorp SQL trace]", tracer)

	t := dbmap.AddTableWithName(EC2{}, "ec2").SetKeys(true, "Id")
	t.ColMap("Id").Rename("id")
	t.ColMap("Type").Rename("type")
	dbmap.DropTables()
	err = dbmap.CreateTablesIfNotExists()
	if err != nil {
		panic(err.Error())
	}

	ec2list := make([]interface{}, 0, 10)
	for i := 0; i < 10; i++ {
		ec2 := &EC2{
			Id:   0,
			Type: "t2.micro" + strconv.Itoa(i),
		}

		ec2list = append(ec2list, ec2)
	}

	dbmap.Insert(ec2list...)

	list, _ := dbmap.Select(EC2{}, "select * from ec2")
	for _, l := range list {
		ec2 := l.(*EC2)
		fmt.Printf("%d, %s\n", ec2.Id, ec2.Type)
	}
}
