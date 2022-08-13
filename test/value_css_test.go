package test

import (
	"os"
	"reader_parse/analyse/css"
	"testing"
)

func TestCssValue(t *testing.T) {
	htmlBytes, err := os.ReadFile("../html/css_test.html")
	if err != nil {
		t.Errorf(err.Error())
	}

	html := string(htmlBytes)

	// html = ConvertToString(html, "GB18030", "utf-8")

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
	} else {
		element := elements[0]
		subAnalyse := css.CssAnalyse{
			Doc: element,
		}

		value, err := subAnalyse.GetValue("span.s2 > a @text")
		if err != nil {
			t.Errorf(err.Error())
		} else {
			if value != "遮天：成帝的我回到地球当保安" {
				t.Errorf("want `遮天：成帝的我回到地球当保安` but %s", value)
			}
		}
	}
}
