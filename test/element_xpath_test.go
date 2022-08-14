package test

import (
	"os"
	"reader_parse/analyse/xpath"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func TestXpathElement(t *testing.T) {
	htmlBytes, err := os.ReadFile("../html/xpath_test.html")

	if err != nil {
		t.Errorf("can not find gsoup_test.html")
		return
	}

	htmlStr := string(htmlBytes)

	doc, err := html.Parse(strings.NewReader(htmlStr))
	if err != nil {
		t.Errorf(err.Error())
	}

	analyse := xpath.XpathAnalyse{
		Doc: doc,
	}

	elements, err := analyse.GetElements("//div[5]/div[2]/div/ul/li")
	if err != nil {
		t.Errorf(err.Error())
	} else {
		if len(elements) != 20 {
			t.Errorf("want 20 but %d", len(elements))
		} else {
			elememt := elements[0]

			subAnalyse := xpath.XpathAnalyse{
				Doc: elememt,
			}

			title, err := subAnalyse.GetValue("//a/@content")
			if err != nil {
				t.Errorf(err.Error())
			} else {
				if title != "霍格沃兹之幕后大bosstxt下载" {
					t.Errorf("want `霍格沃兹之幕后大bosstxt下载` but %s", title)
				}
			}
		}
	}
}
