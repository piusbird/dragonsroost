//go:build dev

//

package main

import (
	"io"
	"log"
	"net/http"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"piusbird.space/dragonsroost/models"
)

func setupDatabase(w http.ResponseWriter, req *http.Request) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		http.Error(w, err.Error()+" Before table create", 550)
	}
	log.Println("Create Tables")
	db.AutoMigrate(&models.Key{})
	db.AutoMigrate(&models.Page{})
	db.AutoMigrate(&models.Post{})
	db.AutoMigrate(&models.LogEntry{})
	posts, err := postsToStructs()
	if err != nil {
		http.Error(w, "Import Error Posts "+err.Error(), 550)
	}
	for p := posts.Front(); p != nil; p = p.Next() {
		vl := p.Value.(models.Post)
		db.Create(&vl)

	}
	pages, err := pagesToStructs()
	for p := pages.Front(); p != nil; p = p.Next() {
		vl := p.Value.(models.Page)
		db.Create(&vl)

	}
	io.WriteString(w, "DB Import successful")
	var testActor models.Key
	testActor.Key = testKey
	testActor.User = "Test"
	db.Create(&testActor)
	return

}
