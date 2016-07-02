package io_ops

import (
	"database/sql"
	"fmt"
	"github.com/d-jo/webserver/structs"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

var db sql.DB
var preparedSelect = "SELECT title, content, good_points, idiom_points FROM snippits WHERE id=? LIMIT 1"
var preparedSelectPoints = "SELECT good_points, idiom_points FROM snippits WHERE id=? LIMIT 1"
var preparedInsert = "INSERT INTO snippits (title, content, good_points, idiom_points) VALUES (?, ?, ?, ?)"
var preparedUpdatePoints = "UPDATE snippits SET good_points=?,idiom_points=? WHERE id=?"

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

func UpdatePointsInDB(goodPointsDelta, idiomPointsDelta int, id string) {
	goodPoints, idiomPoints := getPointsForId(id)
	goodPoints, idiomPoints = goodPoints+goodPointsDelta, idiomPoints+idiomPointsDelta
	db.Exec(preparedUpdatePoints, goodPoints, idiomPoints, id)
}

func getPointsForId(id string) (int, int) {
	var goodPointsScan, idiomPointsScan int
	db.QueryRow(preparedSelectPoints, id).Scan(&goodPointsScan, &idiomPointsScan)
	return goodPointsScan, idiomPointsScan
}

func GetCodeSnipFromDB(id string) (structs.CodeSnip, error) {
	var titleScan, contentScan string
	var goodPointsScan, idiomPointsScan int
	err := db.QueryRow(preparedSelect, id).Scan(&titleScan, &contentScan, &goodPointsScan, &idiomPointsScan)
	return structs.CodeSnip{Title: titleScan, Content: contentScan, GoodPoints: goodPointsScan, IdiomPoints: idiomPointsScan}, err
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

func Close() {
	db.Close()
}

func Test() {

	UpdatePointsInDB(-5, -34, "7")
	fmt.Println("asdf")
	Close()
}
