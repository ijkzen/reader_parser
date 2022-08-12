package test

import (
	"os"
	"testing"

	"reader_parse/analyse/gsoup"

	"github.com/anaskhan96/soup"
)

func TestElements(t *testing.T) {
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

	if elements, err := analyse.GetElements("class.class-ul@class.class-li"); err == nil {
		if len(elements) != 12 {
			t.Errorf("want 12 but %d", len(elements))
		}
	}

	if elements1, err1 := analyse.GetElements("id.id-ul@[0:9]"); err1 == nil {
		if len(elements1) != 10 {
			t.Errorf("want 10 but %d", len(elements1))
		}
	}

	if elements2, err2 := analyse.GetElements("id.id-ul@[0,9]"); err2 == nil {
		if len(elements2) != 2 {
			t.Errorf("want 2 but %d", len(elements2))
		}
	}

	if elements2, err2 := analyse.GetElements("id.id-ul@id.id-li.0:1:2:3"); err2 == nil {
		if len(elements2) != 4 {
			t.Errorf("want 4 but %d", len(elements2))
		}
	}
	if elements2, err2 := analyse.GetElements("id.id-ul@id.id-li!0:1:2:3"); err2 == nil {
		if len(elements2) != 3 {
			t.Errorf("want 3 but %d", len(elements2))
		}
	}
	if elements2, err2 := analyse.GetElements("id.id-ul@id.id-li.-1"); err2 == nil {
		if len(elements2) != 1 {
			t.Errorf("want 1 but %d", len(elements2))
		}
	}
	if elements2, err2 := analyse.GetElements("id.id-ul@id.id-li.0"); err2 == nil {
		if len(elements2) != 1 {
			t.Errorf("want 1 but %d", len(elements2))
		}
	}
	if elements2, err2 := analyse.GetElements("id.id-ul@children[0:9]"); err2 == nil {
		if len(elements2) != 10 {
			t.Errorf("want 10 but %d", len(elements2))
		}
	}
	if elements2, err2 := analyse.GetElements("id.id-ul@children[0:9:2]"); err2 == nil {
		if len(elements2) != 5 {
			t.Errorf("want 5 but %d", len(elements2))
		}
	}
}
