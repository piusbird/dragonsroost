package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"

	"time"

	"github.com/gorilla/mux"
	"github.com/metal3d/go-slugify"
	"golang.org/x/crypto/nacl/sign"

	"github.com/flosch/pongo2/v6"
	"github.com/gorilla/feeds"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
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
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
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

func renderPage(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	switch vars["page"] {
	case "":
		p, err := renderMarkdownPage("index")
		if err != nil {
			http.Error(w, err.Error(), http.StatusExpectationFailed)
			return
		}
		tplPage.ExecuteWriter(pongo2.Context{"title": p.Title, "page": p}, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusExpectationFailed)

		}
	case "blog":
		http.Redirect(w, r, "/blog/", http.StatusMovedPermanently)
		return
	default:
		p, err := renderMarkdownPage(vars["page"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		tplPage.ExecuteWriter(pongo2.Context{"title": p.Title, "page": p}, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusExpectationFailed)

		}
	}

}

func rssFeed(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		http.Error(w, "Error establishing database connection", http.StatusInternalServerError)
		return
	}
	var schema = "http://"
	if req.TLS != nil {
		schema = "https://"
	}
	baseURL := schema + req.Host

	feed := &feeds.Feed{
		Title:       "Notes From the Treefort",
		Link:        &feeds.Link{Href: "https://treefort.piusbird.space/rss"},
		Description: "The Pius Bird's blog",
		Author:      &feeds.Author{Name: "Matt Arnold"},
		Created:     time.Now(),
	}
	var feedItems []*feeds.Item
	var feedData []models.Post
	markdown := goldmark.New(
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
		goldmark.WithExtensions(
			meta.Meta,
		),
	)
	result := db.Where("public = ?", true).Order("created_at DESC").Limit(RSS_NUMPOST).Find(&feedData)
	if result.Error != nil {
		http.Error(w, "Database Error", http.StatusInternalServerError)
		return

	}
	for i := 0; i < len(feedData); i++ {
		item := feedData[i]
		var buf bytes.Buffer
		context := parser.NewContext()
		if err := markdown.Convert(item.Text, &buf, parser.WithContext(context)); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		finalUrl, _ := url.JoinPath(baseURL, "/blog/", item.Slug)

		feedItems = append(feedItems,
			&feeds.Item{
				Id:          string(item.ID),
				Title:       item.Title,
				Link:        &feeds.Link{Href: finalUrl},
				Description: buf.String(),
				Created:     item.CreatedAt,
			})
	}
	feed.Items = feedItems
	var feedCtype string
	switch vars["type"] {
	case "json":
		finalFeed, _ := feed.ToJSON()
		feedCtype = "application/json"
		w.Header().Set("Content-Type", feedCtype)
		io.WriteString(w, finalFeed)
		return
	case "atom":
		finalFeed, _ := feed.ToAtom()
		feedCtype = "application/atom+xml"
		w.Header().Set("Content-Type", feedCtype)
		io.WriteString(w, finalFeed)
		return
	case "rss":
		finalFeed, _ := feed.ToRss()
		feedCtype = "application/rss+xml"
		w.Header().Set("Content-Type", feedCtype)
		io.WriteString(w, finalFeed)
		return

	}

	return
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
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
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
	// TODO Break this into seperate functions later on
	// Protocal Inspired by ssh+ Debian Package uploads
	// The Idea is you toss a signed bit of Json
	// At the server describing what you want to create, delete or unpublish
	// And if you're in the keyring the server does the thing you want to do

	// Uses Ed25519 keys from nacl/sodium for signing/auth

	// TODO Add delete/unpublish verbs
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		http.Error(w, "Cannot open database connection", http.StatusInternalServerError)
		return
	}
	var faildUploadLog models.LogEntry
	faildUploadLog.Type = typeFailed
	faildUploadLog.ForUser = string(r.RemoteAddr)
	faildUploadLog.Affected = "EUSER"
	actor := r.Header.Get("X-Username")
	if actor == "" {
		db.Create(&faildUploadLog)
		http.Error(w, "Must Set Username", http.StatusFailedDependency)
		return
	}
	content := r.Header.Get("Content-Type")
	if content != "application/json+signed" {
		faildUploadLog.Affected = "ESIG"
		db.Create(&faildUploadLog)
		http.Error(w, "Must sign all requests", http.StatusFailedDependency)
		return
	}

	// Only copy into the buffer once we are sure all preconditions are met

	uploadBuffer, _ := io.ReadAll(r.Body)
	var actorInfo models.Key
	res := db.Where("user = ?", actor).First(&actorInfo)
	log.Printf("%v", actorInfo)
	if res.Error != nil {
		faildUploadLog.Affected = "ENOKEY"
		db.Create(&faildUploadLog)
		http.Error(w, "Unable to fetch user", http.StatusInternalServerError)
		return
	}

	buf := string(uploadBuffer)
	msg, err := base64.StdEncoding.DecodeString(buf)
	if err != nil {
		faildUploadLog.Affected = "EIEIO"
		db.Create(&faildUploadLog)
		http.Error(w, "Upload Error", http.StatusInternalServerError)
		return
	}

	log.Printf("%d", len(msg))
	// The public key comes out of the database as a base64 encoded string
	// Needs to go in a 32 byte non encoded buffer for nacl to verify
	var pubKey [32]byte
	pubKeyData, err := base64.StdEncoding.DecodeString(actorInfo.Key)

	if err != nil {
		faildUploadLog.Affected = "EIOKEY"
		db.Create(&faildUploadLog)
		http.Error(w, "Server side key decode error", http.StatusFailedDependency)
		return
	}
	if len(pubKeyData) != 32 {
		http.Error(w, "keydata for user is bad", http.StatusInternalServerError)
		return
	}
	copy(pubKey[:], pubKeyData)

	// The signed data will be encrypted to the public key, because
	// Signing 💡 this returns the unencrypted form if it can

	raw_json, ok := sign.Open(nil, msg, &pubKey)
	if !ok {
		faildUploadLog.Affected = "EBADSIG"
		db.Create(&faildUploadLog)
		http.Error(w, "Crypto error ", http.StatusForbidden)
		return
	}

	var upload JsonUpload
	err = json.Unmarshal(raw_json, &upload)
	if err != nil {
		faildUploadLog.Affected = "EIOJSON"
		db.Create(&faildUploadLog)
		http.Error(w, "Json error "+err.Error(), http.StatusInternalServerError)
		return
	}
	var logTransact models.LogEntry
	logTransact.ForUser = actorInfo.User // Better to log who the user actually was,
	// rather then who they said that they were

	oType, _ := strconv.ParseInt(upload.Type, 10, 16)
	logTransact.Type = int16(oType)
	if oType == Undefined {
		faildUploadLog.Affected = "EWRONGVERB"
		db.Create(&faildUploadLog)
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
			logTransact.Affected = newPost.Slug
			db.Create(&logTransact)
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
		logTransact.Affected = oldPost.Slug
		db.Create(&logTransact)
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
			logTransact.Affected = newPost.ShortName
			db.Create(&logTransact)
			w.Header().Set("Content-Type", "text/plain")
			io.WriteString(w, "ok: "+string(result.RowsAffected)+"\n")
			return

		}
		oldPost.CreatedAt, _ = time.Parse(timelayout, upload.Date)
		oldPost.UpdatedAt = time.Now()
		oldPost.Text = []byte(upload.Body)
		result = db.Save(&oldPost)
		if result.Error != nil {
			faildUploadLog.ForUser = actorInfo.User
			faildUploadLog.Affected = "ESUBMIT"
			db.Create(&faildUploadLog)
			http.Error(w, "Submit Error!", http.StatusInternalServerError)
			return
		}
		logTransact.Affected = oldPost.ShortName
		db.Create(&logTransact)
		w.Header().Set("Content-Type", "text/plain")
		io.WriteString(w, "ok: "+string(result.RowsAffected)+"\n")
		return

	}

}
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		http.Error(w, "Cannot open database connection", http.StatusInternalServerError)
		return
	}
	var faildUploadLog models.LogEntry
	faildUploadLog.Type = typeFailed
	faildUploadLog.ForUser = string(r.RemoteAddr)
	faildUploadLog.Affected = "EUSER"
	actor := r.Header.Get("X-Username")
	if actor == "" {
		db.Create(&faildUploadLog)
		http.Error(w, "Must Set Username", http.StatusFailedDependency)
		return
	}
	signedHash := r.Header.Get("X-Upload-Sig")
	if signedHash == "" {
		faildUploadLog.Affected = "ESIG"
		db.Create(&faildUploadLog)
		http.Error(w, "Must sign all requests", http.StatusFailedDependency)
		return
	}
	fileName := r.Header.Get("X-Upload-Filename")
	if fileName == "" {
		faildUploadLog.Affected = "ENOFILENAME"
		db.Create(&faildUploadLog)
		http.Error(w, "Must set file name", http.StatusBadRequest)
		return
	}

	uploadBuffer, _ := io.ReadAll(r.Body)
	var actorInfo models.Key
	res := db.Where("user = ?", actor).First(&actorInfo)
	log.Printf("%v", actorInfo)
	if res.Error != nil {
		faildUploadLog.Affected = "ENOKEY"
		db.Create(&faildUploadLog)
		http.Error(w, "Unable to fetch user", http.StatusInternalServerError)
		return
	}

	buf := string(signedHash)
	msg, err := base64.StdEncoding.DecodeString(buf)
	if err != nil {
		faildUploadLog.Affected = "EIEIO"
		db.Create(&faildUploadLog)
		http.Error(w, "Upload Error", http.StatusInternalServerError)
		return
	}

	log.Printf("%d", len(msg))
	// The public key comes out of the database as a base64 encoded string
	// Needs to go in a 32 byte non encoded buffer for nacl to verify
	var pubKey [32]byte
	pubKeyData, err := base64.StdEncoding.DecodeString(actorInfo.Key)

	if err != nil {
		faildUploadLog.Affected = "EIOKEY"
		db.Create(&faildUploadLog)
		http.Error(w, "Server side key decode error", http.StatusFailedDependency)
		return
	}
	if len(pubKeyData) != 32 {
		http.Error(w, "keydata for user is bad", http.StatusInternalServerError)
		return
	}
	copy(pubKey[:], pubKeyData)
	submitHash, ok := sign.Open(nil, msg, &pubKey)
	if !ok {
		faildUploadLog.Affected = "ECRYPTO"
		db.Create(&faildUploadLog)
		http.Error(w, "Signiture failed verify", http.StatusBadRequest)
		return
	}
	actualHash := sha256.New()
	io.Copy(actualHash, bytes.NewReader(uploadBuffer))
	shStr := hex.EncodeToString(submitHash)
	ahStr := hex.EncodeToString(actualHash.Sum(nil))
	if ahStr != shStr {
		faildUploadLog.Affected = "EHASH"
		log.Printf("%s vs %s", ahStr, shStr)
		db.Create(&faildUploadLog)
		http.Error(w, "Hash validation failed", http.StatusBadRequest)
		return
	}
	targetPath := filepath.Join("assets", "uploads", fileName)
	err = os.RemoveAll(targetPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInsufficientStorage)
		return
	}
	fh, err := os.Create(targetPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInsufficientStorage)
		return
	}
	io.Copy(fh, bytes.NewReader(uploadBuffer))
	fh.Close()
	w.WriteHeader(http.StatusAccepted)
	io.WriteString(w, "Ok\n")
	return

}
func main() {
	r := mux.NewRouter()
	r.HandleFunc("/{page}", renderPage)
	r.HandleFunc("/", renderPage)
	r.HandleFunc("/api/post", UpdatePost).Methods("POST")
	r.HandleFunc("/api/upload", uploadHandler).Methods("POST")
	r.HandleFunc("/feeds/{type}", rssFeed)

	r.HandleFunc("/setup/", setupDatabase)
	blogroute := r.PathPrefix("/blog").Subrouter()
	blogroute.HandleFunc("/", blogLanding)
	blogroute.HandleFunc("/{slug}", getBlogPost)

	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))
	f, err := os.OpenFile("blog-backend.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {

		log.Fatalf("error opening file: %v", err)

	}

	defer f.Close()

	port := os.Getenv("PORT")
	newdsn := os.Getenv("DSN")
	unixsock := os.Getenv("SOCKPATH")
	if newdsn != "" {
		dsn = newdsn
	}
	if port == "" {
		port = "12345"
	}

	server := http.Server{Handler: r}
	if unixsock != "" {
		if err := os.RemoveAll(unixsock); err != nil {
			log.Fatalln(err)

		}
		unixListener, err := net.Listen("unix", unixsock)

		if err != nil {
			log.Fatalln(err)

		}
		err = os.Chmod(unixsock, 0770)
		if err != nil {
			log.Fatalln(err)
		}
		log.SetOutput(f)
		server.Serve(unixListener)

	} else {
		server.Addr = ":" + port
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
		log.SetOutput(f)

	}

}
