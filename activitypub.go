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
    "published": "1984-04-20T00:00:00Z",

	"inbox": "https://treefort.piusbird.space/inbox",
	"publicKey": {
	  "id": "https://treefort.piusbird.space/u/piusbird/key",
	  "owner": "https://treefort.piusbird.space/u/piusbird",
	  "publickeyPem": "-----BEGIN PUBLIC KEY-----\\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAozrjVu0XYfjYHyk/eoX6\\nhHfTu+qsq+T299wqTJ/vp7l3pAVttUR6eNS5kZNs4Ugr+MheJf/GS2odH8pDSkRt\\naRCq4Cs0uMHBzN+tLaE1lWR6O+fZX6nyNXwax+XRRuAgyr1ciazWGs+1r50JnbJH\\nYI9tofeiq5UaRUyNP2SLvNmDKsXcdL2qw2UQMfJQj/pKAFNnz31xrFvAAQ3S24aC\\nO7bdZ9fIMPfgwudSkga8Hs0ACYW+AXSSDMnwnJ631hSKicU4QAKCbn+cziGR5ZHY\\nbX3xe96NGXJAcA7y9AcQG2REcUvUeXp6XW7KoaBjNN47d/h+F6DmQ/FJPvQUI9bW\\n0QIDAQAB\\n-----END PUBLIC KEY-----\\n"
	}
  }
  `

func webfingeRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/activity+json")
	io.WriteString(w, apActor)
	return
}

func serveUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/activity+json")
	io.WriteString(w, user)
	return

}

func serveProfileHtml(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://treefort.piusbird.space/index", http.StatusMovedPermanently)
	return

}
