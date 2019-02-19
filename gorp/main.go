package main

import (
	"database/sql"
	"fmt"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/gorp.v1"
)

type EC2 struct {
	Id   int64
	SKU  string
	Type string
}

func main() {
	db, err := sql.Open("sqlite3", "./ec2.db")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
	t := dbmap.AddTableWithName(EC2{}, "ec2").SetKeys(true, "Id")
	t.ColMap("Id").Rename("id")
	t.ColMap("SKU").Rename("sku")
	t.ColMap("Type").Rename("type")
	dbmap.DropTables()
	err = dbmap.CreateTablesIfNotExists()
	if err != nil {
		panic(err.Error())
	}

	ec2list := make([]interface{}, 0, 100)
	for i := 0; i < 100; i++ {
		ec2 := &EC2{
			Id:   0,
			SKU:  "abcde",
			Type: "t2.micro" + strconv.Itoa(i),
		}

		ec2list = append(ec2list, ec2)
	}

	dbmap.Insert(ec2list...)

	records := make([]EC2, 0)
	_, err = dbmap.Select(&records, "select * from ec2")
	if err != nil {
		panic(err)
	}

	for _, record := range records {
		fmt.Printf("%d, %s, %s\n", record.Id, record.SKU, record.Type)
	}
}
