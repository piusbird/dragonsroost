package main

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"piusbird.space/dragonsroost/models"
	"piusbird.space/dragonsroost/utils"
)

func main() {
	dsn := os.Getenv("DSN")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Database connection error" + err.Error())
	}
	log.Println("Create Tables")
	db.AutoMigrate(&models.Key{})
	db.AutoMigrate(&models.Page{})
	db.AutoMigrate(&models.Post{})
	db.AutoMigrate(&models.LogEntry{})
	posts, err := utils.PostsToStructs()
	if err != nil {
		panic("Database Failure" + err.Error())
	}
	for p := posts.Front(); p != nil; p = p.Next() {
		vl := p.Value.(models.Post)
		db.Create(&vl)

	}
	pages, err := utils.PagesToStructs()
	for p := pages.Front(); p != nil; p = p.Next() {
		vl := p.Value.(models.Page)
		db.Create(&vl)

	}
	fmt.Println("DB Import successful")
	var testActor models.Key
	testActor.Key = utils.TestKey
	testActor.User = "Test"
	db.Create(&testActor)
	return
}
