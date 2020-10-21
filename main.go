package main

import (
	"net/url"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	gast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"
	"github.com/yuin/goldmark/util"
)

type RelabsHTMLRenderer struct {
	html.Config
	base *url.URL
}

func NewRelabsHTMLRenderer(b *url.URL, opts ...html.Option) renderer.NodeRenderer {
	r := &RelabsHTMLRenderer{
		Config: html.NewConfig(),
		base:   b,
	}
	for _, opt := range opts {
		opt.SetHTMLOption(&r.Config)
	}
	return r
}

func (r *RelabsHTMLRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	reg.Register(gast.KindLink, r.renderLink)
	reg.Register(gast.KindImage, r.renderImage)
}

// copy & pasted from goldmark/renderer/html.go
func (r *RelabsHTMLRenderer) renderImage(w util.BufWriter, source []byte, node gast.Node, entering bool) (gast.WalkStatus, error) {
	if !entering {
		return ast.WalkContinue, nil
	}
	n := node.(*ast.Image)
	_, _ = w.WriteString("<img src=\"")
	if r.Unsafe || !html.IsDangerousURL(n.Destination) {
		dst := r.resolveURL(util.EscapeHTML(util.URLEscape(n.Destination, true)))
		_, _ = w.Write(dst)
	}
	_, _ = w.WriteString(`" alt="`)
	_, _ = w.Write(util.EscapeHTML(n.Text(source)))
	_ = w.WriteByte('"')
	if n.Title != nil {
		_, _ = w.WriteString(` title="`)
		r.Writer.Write(w, n.Title)
		_ = w.WriteByte('"')
	}
	if n.Attributes() != nil {
		html.RenderAttributes(w, n, html.ImageAttributeFilter)
	}
	if r.XHTML {
		_, _ = w.WriteString(" />")
	} else {
		_, _ = w.WriteString(">")
	}
	return ast.WalkSkipChildren, nil
}

func (r *RelabsHTMLRenderer) resolveURL(o []byte) []byte {
	u, err := url.Parse(string(o))
	if err == nil && !u.IsAbs() {
		return []byte(r.base.ResolveReference(u).String())
	}
	return o
}

// copy & pasted from goldmark/renderer/html.go
func (r *RelabsHTMLRenderer) renderLink(w util.BufWriter, source []byte, node gast.Node, entering bool) (gast.WalkStatus, error) {
	n := node.(*ast.Link)
	if entering {
		_, _ = w.WriteString("<a href=\"")
		if r.Unsafe || !html.IsDangerousURL(n.Destination) {
			dst := r.resolveURL(util.EscapeHTML(util.URLEscape(n.Destination, true)))
			_, _ = w.Write(dst)
		}
		_ = w.WriteByte('"')
		if n.Title != nil {
			_, _ = w.WriteString(` title="`)
			r.Writer.Write(w, n.Title)
			_ = w.WriteByte('"')
		}
		if n.Attributes() != nil {
			html.RenderAttributes(w, n, html.LinkAttributeFilter)
		}
		_ = w.WriteByte('>')
	} else {
		_, _ = w.WriteString("</a>")
	}
	return ast.WalkContinue, nil
}

type relabs struct {
	BaseURL *url.URL
}

func NewRelabs(baseURL *url.URL) *relabs {
	return &relabs{baseURL}
}

func (e *relabs) Extend(m goldmark.Markdown) {
	m.Renderer().AddOptions(renderer.WithNodeRenderers(
		util.Prioritized(NewRelabsHTMLRenderer(e.BaseURL), 500),
	))
}
