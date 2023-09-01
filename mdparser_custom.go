package main

// Note i made chatgpt do the boring parts
import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/flosch/pongo2/v6"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
	xhtml "golang.org/x/net/html"
)

func parseRewriteMarkdown(text []byte) ([]byte, error) {
	var tplCreature = pongo2.Must(pongo2.FromFile("templates/creature.html"))
	var output bytes.Buffer
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
	if err := markdown.Convert(text, &buf, parser.WithContext(context)); err != nil {
		return nil, errors.New("could not parse markdown")
	}
	r := bytes.NewReader(buf.Bytes())
	origTree, err := xhtml.Parse(r)
	if err != nil {
		return nil, errors.New("could not parse markdown")
	}
	for c := origTree.FirstChild; c != nil; c = c.NextSibling {
		customNode := findElement(c, creature_tag)
		if customNode != nil {
			msg := fmt.Sprint(customNode.FirstChild.Data)
			ctx, err := gatherAttribs(customNode, creature_tag)
			if err != nil {
				return nil, errors.New("level 2 parser error")

			}
			ctx["message"] = msg
			var ourCTX pongo2.Context = pongo2.Context{"creature": ctx["creature"], "mood": ctx["mood"], "message": ctx["message"]}
			s, _ := tplCreature.Execute(ourCTX)
			repNode, err := xhtml.Parse(strings.NewReader(s))
			if err != nil {
				return nil, errors.New("level 3 parsing error")
			}
			replaceNode(customNode, repNode)

		}
	}
	xhtml.Render(&output, origTree)
	return output.Bytes(), nil

}

func gatherAttribs(node *xhtml.Node, target string) (map[string]string, error) {
	retmap := make(map[string]string)
	if node.Type == xhtml.ElementNode && node.Data == target {
		// Get the value of the href attribute
		for _, attr := range node.Attr {
			retmap[attr.Key] = attr.Val
		}
	} else {
		return retmap, errors.New("couldn't parse custom element")
	}
	return retmap, nil

}

// Boring part

func findElement(node *xhtml.Node, tagName string) *xhtml.Node {
	if node.Type == xhtml.ElementNode && node.Data == tagName {
		return node
	}
	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if found := findElement(child, tagName); found != nil {
			return found
		}
	}
	return nil

}

func replaceNode(targetNode *xhtml.Node, newNode *xhtml.Node) {
	parent := targetNode.Parent
	if parent == nil {
		return
	}

	// Find the position of the target node within its parent's children
	var prevSibling *xhtml.Node
	for child := parent.FirstChild; child != nil; child = child.NextSibling {
		if child == targetNode {
			break
		}
		prevSibling = child
	}

	if prevSibling != nil {
		// Insert the new node after the previous sibling
		parent.InsertBefore(newNode, targetNode)
		parent.RemoveChild(targetNode)
	} else {
		// The target node was the first child, so replace it directly
		parent.RemoveChild(targetNode)
		parent.AppendChild(newNode)
	}
}

// End boring
