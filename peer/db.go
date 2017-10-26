package main

import (
	"database/sql"
	"fmt"
	scmn "github.com/kenmazsyma/soila/peer/common"
	_ "github.com/lib/pq"
)

type TEST struct {
	A  string
	ID int
}

func _main() {
	scmn.ReadEnv("soila.conf")
	db, err := sql.Open(scmn.Env.DB_DRIVER, scmn.Env.DB_ARGS)
	if err != nil {
		fmt.Printf("ERROR:%s\n", err.Error())
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT ID, A FROM TEST WHERE ID = $1;", 15)
	if err != nil {
		fmt.Printf("ERROR:%s\n", err.Error())
		return
	}

	column := TEST{}
	for rows.Next() {
		err = rows.Scan(&column.ID, &column.A)
		if err != nil {
			fmt.Printf("ERROR:%s\n", err.Error())
			return
		}
		fmt.Printf("ID:%d, A:%s\n", column.ID, column.A)
	}
}
