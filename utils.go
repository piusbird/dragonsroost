package main

import (
	"bytes"
	"container/list"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/metal3d/go-slugify"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
	"piusbird.space/dragonsroost/models"
)

var dsn = "web:changethis@unix(/var/lib/mysql/mysql.sock)/website?charset=utf8mb4&parseTime=True&loc=Local"
var timelayout = "2006-01-02"

type Page struct {
	Title string
	Html  string
}
type PostMeta struct {
	Title     string
	Date      string
	ShortDesc string
	Url       string
}
type EnityType int16

const (
	Undefined  = 0
	typePage   = 1
	typePost   = 2
	typeFailed = 3
)
const RSS_NUMPOST = 10

type JsonUpload struct {
	Type  string
	Title string
	Date  string
	Body  string
}

var testKey = "nR6FC9GHk+olScO5FPpUYoppgo95SHvd5UJKKFt4Crs="

type PostMetaByDate []PostMeta

func (a PostMetaByDate) Len() int      { return len(a) }
func (a PostMetaByDate) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a PostMetaByDate) Less(i, j int) bool {
	timeI, _ := time.Parse(time.RFC1123, a[i].Date)
	timeJ, _ := time.Parse(a[j].Date, time.RFC1123)
	return timeI.Before(timeJ)
}
func postsToStructs() (list.List, error) {
	files, err := ioutil.ReadDir("content/posts")
	if err != nil {
		return *list.New(), err

	}
	postList := list.New()

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".md") {
			continue
		}
		var absPath = filepath.Join("content/posts", file.Name())
		var post = models.Post{}
		log.Println("BINGO! " + absPath)
		raw, err := os.ReadFile(absPath)
		if err != nil {
			panic("Disk Error!")

		}
		markdown := goldmark.New(
			goldmark.WithExtensions(
				meta.Meta,
			),
		)
		var buf bytes.Buffer
		context := parser.NewContext()
		if err := markdown.Convert(raw, &buf, parser.WithContext(context)); err != nil {
			return *list.New(), err
		}
		metaData := meta.Get(context)
		post.Title = "Unknown Post"
		if metaData["title"] != nil {
			post.Title = metaData["title"].(string)

		}
		st_time := file.ModTime()
		if metaData["date"] != nil {
			st_time, _ = time.Parse(timelayout, metaData["date"].(string))
		}
		post.CreatedAt = st_time
		post.UpdatedAt = time.Now()
		post.Text = raw
		post.Public = true
		post.Slug = slugify.Marshal(post.Title)
		postList.PushBack(post)

	}

	return *postList, nil
}

func pagesToStructs() (list.List, error) {
	files, err := ioutil.ReadDir("content/pages")
	if err != nil {
		return *list.New(), err

	}
	postList := list.New()

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".md") {
			continue
		}
		var post = models.Page{}
		var absPath = filepath.Join("content/pages", file.Name())
		raw, _ := os.ReadFile(absPath)
		markdown := goldmark.New(
			goldmark.WithExtensions(
				meta.Meta,
			),
		)
		var buf bytes.Buffer
		context := parser.NewContext()
		if err := markdown.Convert(raw, &buf, parser.WithContext(context)); err != nil {
			return *list.New(), err
		}
		metaData := meta.Get(context)
		post.Title = "Unknown"
		if metaData["Title"] != nil {
			post.Title = metaData["Title"].(string)
		}

		post.CreatedAt = time.Now()
		post.UpdatedAt = time.Now()
		post.Text = raw
		post.Sidebar = true
		post.ShortName = strings.TrimSuffix(file.Name(), ".md")
		postList.PushBack(post)

	}

	return *postList, nil
}
