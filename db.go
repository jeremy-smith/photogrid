package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var conn *sql.DB

func dbConnect() *sql.DB {
	var err error
	conn, err = sql.Open("mysql", fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8",
		Conf.Db.DbUser,
		Conf.Db.DbPass,
		Conf.Db.DbHost,
		Conf.Db.DbPort,
		Conf.Db.DbName,
	))
	if err != nil {
		log.Fatalln(err)
	}

	return conn
}

func dbCreateAccount(email string, password string) error {
	password = HashPassword(password)

	_, err := conn.Exec(
		"INSERT INTO users (username, password) VALUES (?, ?)",
		email,
		password)
	if err != nil {
		log.Println(err)
	}
	return err
}

func validateLogin(username string, password string) bool {
	row := conn.QueryRow("SELECT password FROM users WHERE username = ? LIMIT 1", username)
	var dbPsw string
	if err := row.Scan(&dbPsw); err == sql.ErrNoRows {
		log.Println(err)
		return false
	}
	return ValidatePassword(dbPsw, password)
}

//	res, err := db.Query("select * from amigos")
//	if err != nil {
//		log.Fatalln(err)
//	}
//	defer res.Close()
//	for res.Next() {
//		id := new(int)
//		name := new(string)
//		err := res.Scan(id, name)
//		if err != nil {
//			log.Fatalln(err)
//		}
//
//		fmt.Println(*id, *name)
//	}
//}
