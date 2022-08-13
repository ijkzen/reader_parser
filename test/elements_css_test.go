package test

import (
	"os"
	"reader_parse/analyse/css"
	"testing"
)

func TestCssElement(t *testing.T) {
	htmlBytes, err := os.ReadFile("../html/css_test.html")

	if err != nil {
		t.Errorf("can not find gsoup_test.html")
		return
	}

	html := string(htmlBytes)

	analyse, err := css.GetCssAnalyseFromString(html)
	if err != nil {
		t.Errorf(err.Error())
	}

	elements, err := analyse.GetElements("#main > div.novelslistss > li")
	if err != nil {
		t.Errorf(err.Error())
	}

	if len(elements) != 46 {
		t.Errorf("want 46 but %d", len(elements))
	}
}
