package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"time"

	"github.com/gorilla/mux"
	"github.com/metal3d/go-slugify"
	"golang.org/x/crypto/nacl/sign"

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
		tmp := "Lol"
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
	db.AutoMigrate(&models.Key{})
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
	var testActor models.Key
	testActor.Key = testKey
	testActor.User = "Test"
	db.Create(&testActor)
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
func UpdatePost(w http.ResponseWriter, r *http.Request) {
	nBytes, _ := io.ReadAll(r.Body)
	log.Printf("%v", nBytes)
	actor := r.Header.Get("X-Username")
	if actor == "" {
		http.Error(w, "Must Set Username", http.StatusFailedDependency)
		return
	}
	content := r.Header.Get("Content-Type")
	if content != "application/json+signed" {
		http.Error(w, "Must sign all requests", http.StatusFailedDependency)
		return
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		http.Error(w, "Cannot open database connection", http.StatusInternalServerError)
		return
	}
	var actorInfo models.Key
	res := db.Where("user = ?", actor).First(&actorInfo)
	log.Printf("%v", actorInfo)
	if res.Error != nil {
		http.Error(w, "Unable to fetch user", http.StatusInternalServerError)
		return
	}

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "body fetch error", http.StatusFailedDependency)
		return
	}
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Upload Error", http.StatusInternalServerError)
		return
	}
	log.Printf("%s %d", "Len of data", len(data))

	buf := string(nBytes)
	msg, err := base64.StdEncoding.DecodeString(buf)
	if err != nil {
		http.Error(w, "Upload Error", http.StatusInternalServerError)
		return
	}

	log.Printf("%d", len(msg))
	var pubKey [32]byte
	pubKeyData, err := base64.StdEncoding.DecodeString(testKey)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Server side key decode error", http.StatusFailedDependency)
		return
	}
	if len(pubKeyData) != 32 {
		http.Error(w, "keydat for user is bad", http.StatusInternalServerError)
		return
	}
	copy(pubKey[:], pubKeyData)

	raw_json, ok := sign.Open(nil, msg, &pubKey)
	if !ok {
		http.Error(w, "Crypto error ", http.StatusForbidden)
		return
	}
	log.Printf(string(raw_json))
	var upload JsonUpload
	err = json.Unmarshal(raw_json, &upload)
	if err != nil {
		http.Error(w, "Json error "+err.Error(), http.StatusInternalServerError)
		return
	}
	oType, _ := strconv.ParseInt(upload.Type, 10, 16)
	if oType == Undefined {
		http.Error(w, "Undefined Type", http.StatusForbidden)
		return
	}
	if oType == typePost {
		var newPost models.Post
		var oldPost models.Post
		newPost.Title = upload.Title
		newPost.Slug = slugify.Marshal(upload.Title)
		newPost.CreatedAt, _ = time.Parse(timelayout, upload.Date)
		newPost.UpdatedAt = time.Now()
		newPost.Text = []byte(upload.Body)
		result := db.Model(&newPost).Where("title = ?", newPost.Title).First(&oldPost)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			result = db.Create(&newPost)
			defer r.Body.Close()
			w.Header().Set("Content-Type", "text/plain")
			io.WriteString(w, "ok: "+string(result.RowsAffected))
			return

		}
		oldPost.CreatedAt, _ = time.Parse(timelayout, upload.Date)
		oldPost.UpdatedAt = time.Now()
		oldPost.Text = []byte(upload.Body)
		result = db.Save(&oldPost)
		if result.Error != nil {
			http.Error(w, "Submit Error!", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		io.WriteString(w, "ok: "+string(result.RowsAffected)+"\n")
		return

	}
	if oType == typePage {
		var newPost models.Page
		var oldPost models.Page
		newPost.Title = upload.Title
		newPost.ShortName = slugify.Marshal(upload.Title)
		newPost.CreatedAt, _ = time.Parse(timelayout, upload.Date)
		newPost.UpdatedAt = time.Now()
		newPost.Text = []byte(upload.Body)
		result := db.Model(&newPost).Where("title = ?", newPost.Title).First(&oldPost)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			result = db.Create(&newPost)
			defer r.Body.Close()
			w.Header().Set("Content-Type", "text/plain")
			io.WriteString(w, "ok: "+string(result.RowsAffected)+"\n")
			return

		}
		oldPost.CreatedAt, _ = time.Parse(timelayout, upload.Date)
		oldPost.UpdatedAt = time.Now()
		oldPost.Text = []byte(upload.Body)
		result = db.Save(&oldPost)
		if result.Error != nil {
			http.Error(w, "Submit Error!", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		io.WriteString(w, "ok: "+string(result.RowsAffected)+"\n")
		return

	}

}
func main() {
	r := mux.NewRouter()
	r.HandleFunc("/{page}", renderPage)
	r.HandleFunc("/", renderPage)
	r.HandleFunc("/images/", imageRedir)
	r.HandleFunc("/api/upload", UpdatePost).Methods("POST")

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
