package main

import (
	"bytes"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/yuin/goldmark"
)

var exampleURL, _ = url.Parse("http://ex.org/level1/level2/")

func TestRenderAbsoluteLinkUnchanged(t *testing.T) {
	rl := NewRelabs(exampleURL)
	md := goldmark.New(goldmark.WithExtensions(rl))
	var buf bytes.Buffer

	src := `[desc](http://ex.org/ "title")`
	expected := `<p><a href="http://ex.org/" title="title">desc</a></p>
`
	err := md.Convert([]byte(src), &buf)
	require.Nil(t, err, "failed to convert markdown")
	require.Equal(t, expected, buf.String())
}

func TestRenderAbsoluteLinkWithoutTitleUnchanged(t *testing.T) {
	rl := NewRelabs(exampleURL)
	md := goldmark.New(goldmark.WithExtensions(rl))
	var buf bytes.Buffer

	src := `[desc](http://ex.org/)`
	expected := `<p><a href="http://ex.org/">desc</a></p>
`
	err := md.Convert([]byte(src), &buf)
	require.Nil(t, err, "failed to convert markdown")
	require.Equal(t, expected, buf.String())
}

func TestRenderRelativeLinkResolved(t *testing.T) {
	rl := NewRelabs(exampleURL)
	md := goldmark.New(goldmark.WithExtensions(rl))
	var buf bytes.Buffer

	src := `[desc](../index.html "title")`
	expected := `<p><a href="http://ex.org/level1/index.html" title="title">desc</a></p>
`
	err := md.Convert([]byte(src), &buf)
	require.Nil(t, err, "failed to convert markdown")
	require.Equal(t, expected, buf.String())
}

func TestRenderImageSrcUnchanged(t *testing.T) {
	rl := NewRelabs(exampleURL)
	md := goldmark.New(goldmark.WithExtensions(rl))
	var buf bytes.Buffer

	src := `[![img-alt](http://ex.org/eg.jpg "img-title")](http://ex.org/)`
	expected := `<p><a href="http://ex.org/"><img src="http://ex.org/eg.jpg" alt="img-alt" title="img-title"></a></p>
`
	err := md.Convert([]byte(src), &buf)
	require.Nil(t, err)
	require.Equal(t, expected, buf.String())
}

func TestRenderImageSrcResolved(t *testing.T) {
	rl := NewRelabs(exampleURL)
	md := goldmark.New(goldmark.WithExtensions(rl))
	var buf bytes.Buffer

	src := `[![img-alt](eg.jpg "img-title")](../../index.html)`
	expected := `<p><a href="http://ex.org/index.html"><img src="http://ex.org/level1/level2/eg.jpg" alt="img-alt" title="img-title"></a></p>
`
	err := md.Convert([]byte(src), &buf)
	require.Nil(t, err)
	require.Equal(t, expected, buf.String())
}
