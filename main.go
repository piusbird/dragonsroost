package main

import (
	"bytes"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"time"

	"github.com/gorilla/mux"

	"github.com/flosch/pongo2/v6"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"piusbird.space/dragonsroost/models"
)

func renderMarkdownPage(pagename string) (Page, error) {

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return Page{}, err
	}
	var result models.Page
	db.Where("short_name = ?", pagename).First(&result)
	if result.Text == nil {
		return Page{}, errors.New("Not found in db")

	}
	markdown := goldmark.New(
		goldmark.WithExtensions(
			meta.Meta,
		),
	)
	var buf bytes.Buffer
	context := parser.NewContext()
	if err := markdown.Convert(result.Text, &buf, parser.WithContext(context)); err != nil {
		return Page{}, err
	}
	metaData := meta.Get(context)
	title := metaData["Title"]
	render := buf.String()
	retPage := Page{}
	retPage.Title = title.(string)
	retPage.Html = render
	return retPage, nil
}

var tplPage = pongo2.Must(pongo2.FromFile("templates/page.html"))

func imageRedir(w http.ResponseWriter, r *http.Request) {
	oldpath := r.URL.Path
	assetPath := strings.TrimPrefix("/blog/", oldpath)
	newurl, _ := url.JoinPath(r.Host, "/assets/", assetPath)
	io.WriteString(w, newurl)
	//http.Redirect(w, r, newurl, http.StatusMovedPermanently)
	return
}
func renderPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if vars["page"] == "" {
		p, err := renderMarkdownPage("index")
		if err != nil {
			http.Error(w, err.Error(), http.StatusExpectationFailed)
		}
		tplPage.ExecuteWriter(pongo2.Context{"title": p.Title, "page": p}, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusExpectationFailed)

		}
	}
	p, err := renderMarkdownPage(vars["page"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusExpectationFailed)
		return
	}
	tplPage.ExecuteWriter(pongo2.Context{"title": p.Title, "page": p}, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusExpectationFailed)
		return
	}

}

var tplBlogLanding = pongo2.Must(pongo2.FromFile("templates/blog.html"))

func blogLanding(w http.ResponseWriter, req *http.Request) {
	var schema = "http://"
	if req.TLS != nil {
		schema = "https://"
	}
	baseURL := schema + req.Host
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		http.Error(w, err.Error()+" Failed to open database", 550)
	}
	var allPost []models.Post
	result := db.Order("created_at DESC").Find(&allPost)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusBadRequest)
		return
	}
	var newPost []PostMeta
	// You might be wondering oh why doesn't he just put the url in the db
	// STATELESSNESS is the answer
	for _, post := range allPost {
		finalUrl, _ := url.JoinPath(baseURL, "/blog/", post.Slug)
		pm := PostMeta{}
		pm.Date = post.CreatedAt.Format(time.RFC1123)
		pm.Title = post.Title
		tmp := bytes.NewBuffer(post.Text[0:499]).String()
		pm.ShortDesc = tmp
		pm.Url = finalUrl
		newPost = append(newPost, pm)

	}
	//sort.Sort(PostMetaByDate(newPost))
	tplBlogLanding.ExecuteWriter(pongo2.Context{"postList": newPost, "title": "Notes from the treefort", "total": len(newPost)}, w)

}

func setupDatabase(w http.ResponseWriter, req *http.Request) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		http.Error(w, err.Error()+" Before table create", 550)
	}
	log.Println("Create Tables")

	db.AutoMigrate(&models.Page{})
	db.AutoMigrate(&models.Post{})
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
	return

}

var tplBlog = pongo2.Must(pongo2.FromFile("templates/post.html"))

func getBlogPost(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	slug := vars["slug"]
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		http.Error(w, "Cannot open database connection", http.StatusInternalServerError)
		return
	}
	var result models.Post
	ok := db.Where("slug = ?", slug).First(&result)
	if ok.Error != nil {
		http.Error(w, ok.Error.Error(), http.StatusNotFound)
	}
	markdown := goldmark.New(
		goldmark.WithExtensions(
			meta.Meta,
		),
	)
	var buf bytes.Buffer
	context := parser.NewContext()
	if err := markdown.Convert(result.Text, &buf, parser.WithContext(context)); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var retPage Page
	retPage.Title = result.Title
	retPage.Html = buf.String()
	tplBlog.ExecuteWriter(pongo2.Context{"title": retPage.Title, "post": retPage, "metadata": result}, w)

}
func main() {
	r := mux.NewRouter()
	r.HandleFunc("/{page}", renderPage)
	r.HandleFunc("/", renderPage)
	r.HandleFunc("/images/", imageRedir)

	r.HandleFunc("/setup/", setupDatabase)
	blogroute := r.PathPrefix("/blog").Subrouter()
	blogroute.HandleFunc("/", blogLanding)
	blogroute.PathPrefix("/images/").Handler(http.StripPrefix("/images/", http.FileServer(http.Dir("./assets/images"))))
	blogroute.HandleFunc("/{slug}", getBlogPost)

	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))
	port := os.Getenv("PORT")
	if port == "" {
		port = "12345"
	}

	http.Handle("/", r)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
