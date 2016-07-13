package io_ops

import (
	"database/sql"
	"fmt"
	"github.com/d-jo/webserver/structs"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

var db sql.DB
var preparedSelect = "SELECT title, author, content, good_points, idiom_points FROM snippits WHERE id=? LIMIT 1"
var preparedSelectPoints = "SELECT good_points, idiom_points FROM snippits WHERE id=? LIMIT 1"
var preparedInsert = "INSERT INTO snippits (title, author, content, good_points, idiom_points) VALUES (?, ?, ?, ?, ?)"
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

func UpdatePointsInDB(goodPointsDelta, idiomPointsDelta int, id string) (int, int){
	goodPoints, idiomPoints := getPointsForId(id)
	goodPoints, idiomPoints = goodPoints+goodPointsDelta, idiomPoints+idiomPointsDelta
	db.Exec(preparedUpdatePoints, goodPoints, idiomPoints, id)
	return goodPoints, idiomPoints
}

func getPointsForId(id string) (int, int) {
	var goodPointsScan, idiomPointsScan int
	db.QueryRow(preparedSelectPoints, id).Scan(&goodPointsScan, &idiomPointsScan)
	return goodPointsScan, idiomPointsScan
}

func GetCodeSnipFromDB(id string) (*structs.CodeSnip, error) {
	var titleScan, authorScan, contentScan string
	var goodPointsScan, idiomPointsScan int
	err := db.QueryRow(preparedSelect, id).Scan(&titleScan, &authorScan, &contentScan, &goodPointsScan, &idiomPointsScan)
	return &structs.CodeSnip{Title: titleScan, Author: authorScan, Content: contentScan, GoodPoints: goodPointsScan, IdiomPoints: idiomPointsScan}, err
}

func InsertCodeSnipToDB(snip *structs.CodeSnip) (int, error) {
	res, err := db.Exec(preparedInsert, snip.Title, snip.Author, snip.Content, snip.GoodPoints, snip.IdiomPoints)
	if err != nil {
		return -1, err
	}
	lastId, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}
	return int(lastId), nil
}

func Close() {
	db.Close()
}
