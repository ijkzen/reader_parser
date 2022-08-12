package test

import (
	"os"
	"reader_parse/analyse/gsoup"
	"testing"

	"github.com/anaskhan96/soup"
)

func TestValue(t *testing.T) {
	htmlBytes, err := os.ReadFile("../html/gsoup_test.html")

	if err != nil {
		t.Errorf("can not find gsoup_test.html")
		return
	}

	html := string(htmlBytes)

	doc := soup.HTMLParse(html)

	analyse := gsoup.SoupAnalyse{
		Element: &doc,
	}

	if value, err := analyse.GetValue("tag.ul.0@tag.li.0@tag.a.0@text"); err == nil {
		if value != "第一书名" {
			t.Errorf("want 第一书名 but %s", value)
		}
	}

	if value, err := analyse.GetValue("tag.ul.0@tag.li.0@tag.a.0@href"); err == nil {
		if value != "link" {
			t.Errorf("want link but %s", value)
		}
	}
}
