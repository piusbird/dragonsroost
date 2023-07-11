package main

import (
	"io"
	"net/http"
)

var apActor = `{
	"subject": "acct:piusbird@treefort.piusbird.space",
	"aliases": [
	  "https://treefort.piusbird.space/@piusbird",
	  "https://treefort.piusbird.space/users/piusbird"
	],
	"links": [
	
	  {
		"rel": "self",
		"type": "application/activity+json",
		"href": "https://treefort.piusbird.space/u/piusbird"
	  }
	]
  }
  `
var user = `{
	"@context": [
	  "https://www.w3.org/ns/activitystreams",
	  "https://w3id.org/security/v1"
	],
	"id": "https://treefort.piusbird.space/u/piusbird",
	"type": "Person",
	"preferredUsername": "piusbird",
	"name": "Pius the Unscrootched",
	"summary": "Server Server Server",
	"manuallyApprovesFollowers": true,
    "discoverable": true,
    "published": "2019-04-20T00:00:00Z",

	"inbox": "https://treefort.piusbird.space/inbox",
	"outbox": "https://treefort.piusbird.space/outbox",
	
	
  }
  `

func webfingeRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/activity+json")
	io.WriteString(w, apActor)

}

func serveUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/activity+json")
	io.WriteString(w, user)

}

func serveProfileHtml(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://treefort.piusbird.space/index", http.StatusMovedPermanently)

}
