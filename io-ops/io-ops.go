package io_ops

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"github.com/d-jo/webserver/structs"
	"os"
)

var db sql.DB
var preparedSelect = "SELECT title, content, good_points, idiom_points FROM snippits WHERE id=? LIMIT 1"
var preparedInsert = "INSERT INTO snippits (title, content, good_points, idiom_points) VALUES (?, ?, ?, ?)"

func init() {
	if len(os.Args[1:]) != 3 {
		panic("Missing command line arg(s): user:pass dbserverip dataname")
	}
	userpass := os.Args[1]
	serverip := os.Args[2]
	dataname := os.Args[3]
	fullstring := fmt.Sprintf("%s@tcp(%s)/%s", userpass, serverip, dataname)
	database, err := sql.Open("mysql", fullstring)
	db = *database
	if err != nil {
		panic(err.Error())
	}
}

func GetCodeSnipFromDB(id string) structs.CodeSnip {
	var titleScan, contentScan string
	var goodPointsScan, idiomPointsScan int
	db.QueryRow(preparedSelect, id).Scan(&titleScan, &contentScan, &goodPointsScan, &idiomPointsScan)
	return structs.CodeSnip{Title: titleScan, Content: contentScan, GoodPoints: goodPointsScan, IdiomPoints: idiomPointsScan}
}

func InsertCodeSnipToDB(snip structs.CodeSnip) int {
	res, err := db.Exec(preparedInsert, snip.Title, snip.Content, snip.GoodPoints, snip.IdiomPoints)
	if err != nil {
		panic(err.Error())
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		fmt.Errorf("Error! %s", err.Error())
		return 0
	}
	return int(lastId)
}

func Test() {

	InsertCodeSnipToDB(structs.CodeSnip{Title:"testing with command line", Content:"if this works this should insert with commabnd line args", GoodPoints:5, IdiomPoints:34})
	fmt.Println("asdf")
	db.Close()

}