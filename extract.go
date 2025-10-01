package main

import (
	"log/slog"
	"regexp"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
)

// extract parses the provided markdown content and extracts k6 extension module list from it.
func extract(contents []byte) ([]string, error) {
	parser := parser.NewParser(
		parser.WithBlockParsers(parser.DefaultBlockParsers()...),
		parser.WithInlineParsers(parser.DefaultInlineParsers()...),
		parser.WithParagraphTransformers(parser.DefaultParagraphTransformers()...),
	)
	root := parser.Parse(text.NewReader(contents)).OwnerDocument()
	builder := newBuilder(contents)

	err := ast.Walk(root, builder.walk)
	if err != nil {
		return nil, err
	}

	return builder.build()
}

//nolint:gochecknoglobals
var (
	reLink = regexp.MustCompile(
		`^https://(?P<host>github\.com|gitlab\.com)/(?P<owner>[^/]+)/(?P<repo>[^/]+)(/releases/tag/(?P<tag>[^/]+))?$`,
	)

	idxHost  = reLink.SubexpIndex("host")
	idxOwner = reLink.SubexpIndex("owner")
	idxRepo  = reLink.SubexpIndex("repo")
	idxTag   = reLink.SubexpIndex("tag")
)

type builder struct {
	found bool

	insideList          bool
	hasNonmatchingItems bool

	insideItem         bool
	hasNonmatchingLink bool
	numberOfLinks      int

	links [][]byte

	source []byte
}

func newBuilder(source []byte) *builder {
	b := new(builder)

	b.source = source

	return b
}

func (b *builder) walk(node ast.Node, entering bool) (ast.WalkStatus, error) {
	if b.found {
		return ast.WalkStop, nil
	}

	if list := asList(node); list != nil {
		return b.handleList(entering)
	}

	if item := asListItem(node); item != nil {
		return b.handleListItem(entering)
	}

	if (asLink(node) != nil || asAutoLink(node) != nil) && b.insideItem {
		return b.handleLink(node, entering)
	}

	if text := asText(node); text != nil {
		return b.handleText(text, entering)
	}

	return ast.WalkContinue, nil
}

func (b *builder) handleList(entering bool) (ast.WalkStatus, error) {
	if entering {
		slog.Debug("entering list")

		b.insideList = true
		b.hasNonmatchingItems = false
		b.links = nil

		return ast.WalkContinue, nil
	}

	slog.Debug("exiting list")

	if !b.hasNonmatchingItems {
		b.found = true

		return ast.WalkStop, nil
	}

	b.insideList = false

	return ast.WalkContinue, nil
}

func (b *builder) handleListItem(entering bool) (ast.WalkStatus, error) {
	if entering {
		b.insideItem = true
		b.hasNonmatchingLink = false
		b.numberOfLinks = 0

		return ast.WalkContinue, nil
	}

	b.insideItem = false

	if b.numberOfLinks != 1 {
		slog.Debug("list item does not have exactly one link", "links", b.numberOfLinks)

		b.hasNonmatchingItems = true

		return ast.WalkContinue, nil
	}

	return ast.WalkContinue, nil
}

func (b *builder) handleLink(node ast.Node, entering bool) (ast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}

	var url []byte

	if link := asLink(node); link != nil {
		url = link.Destination
	} else if link := asAutoLink(node); link != nil {
		url = link.URL(b.source)
	}

	b.numberOfLinks++
	b.links = append(b.links, url)

	if !isLinkMatching(url) {
		slog.Debug("found non-matching link", "url", string(url))

		b.hasNonmatchingLink = true
	}

	return ast.WalkContinue, nil
}

func (b *builder) handleText(text *ast.Text, entering bool) (ast.WalkStatus, error) {
	if !entering || !b.insideItem || b.numberOfLinks > 0 {
		return ast.WalkContinue, nil
	}

	value := text.Segment.Value(b.source)
	if isLinkMatching(value) {
		b.links = append(b.links, value)
		b.numberOfLinks++
	}

	return ast.WalkContinue, nil
}

func isLinkMatching(url []byte) bool {
	return reLink.Match(url)
}

func link2module(link []byte) (bool, string) {
	matches := reLink.FindSubmatch(link)

	if len(matches) == 0 {
		return false, ""
	}

	path := string(matches[idxHost]) + "/" +
		string(matches[idxOwner]) + "/" +
		string(matches[idxRepo])

	version := ""

	if tag := matches[idxTag]; len(tag) > 0 {
		version = "@" + string(tag)
	}

	return true, path + version
}

func (b *builder) build() ([]string, error) {
	if !b.found || len(b.links) == 0 {
		return []string{}, nil
	}

	mods := make([]string, 0, len(b.links))

	for _, link := range b.links {
		ok, mod := link2module(link)
		if ok {
			mods = append(mods, mod)
		}
	}

	return mods, nil
}

func asList(node ast.Node) *ast.List {
	if node.Kind() != ast.KindList {
		return nil
	}

	if list, ok := node.(*ast.List); ok {
		return list
	}

	return nil
}

func asListItem(node ast.Node) *ast.ListItem {
	if node.Kind() != ast.KindListItem {
		return nil
	}

	if item, ok := node.(*ast.ListItem); ok {
		return item
	}

	return nil
}

func asLink(node ast.Node) *ast.Link {
	if node.Kind() != ast.KindLink {
		return nil
	}

	if link, ok := node.(*ast.Link); ok {
		return link
	}

	return nil
}

func asAutoLink(node ast.Node) *ast.AutoLink {
	if node.Kind() != ast.KindAutoLink {
		return nil
	}

	if link, ok := node.(*ast.AutoLink); ok {
		return link
	}

	return nil
}

func asText(node ast.Node) *ast.Text {
	if node.Kind() != ast.KindText {
		return nil
	}

	if text, ok := node.(*ast.Text); ok {
		return text
	}

	return nil
}
