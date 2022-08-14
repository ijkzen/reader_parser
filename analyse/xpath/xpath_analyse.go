package xpath

import (
	"fmt"
	"strings"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

type XpathAnalyse struct {
	Doc *html.Node
}

func (analyse *XpathAnalyse) GetElements(rule string) ([]*html.Node, error) {
	elements, err := htmlquery.QueryAll(analyse.Doc, rule)
	if err != nil {
		return make([]*html.Node, 0), err
	}

	return elements, nil
}

func (analyse *XpathAnalyse) GetValue(rule string) (string, error) {
	atIndex := strings.LastIndex(rule, "@")
	if atIndex < 0 {
		return "", fmt.Errorf("rule %s is invalid", rule)
	} else {
		preRule := rule[0:atIndex]
		if preRule[len(preRule)-1] == '/' {
			preRule = preRule[0 : len(preRule)-1]
		}

		attribute := rule[atIndex+1:]
		element, err := htmlquery.Query(analyse.Doc, preRule)
		if err != nil {
			return "", err
		}

		switch attribute {
		case "text", "textNodes", "ownText", "content":
			{
				return htmlquery.InnerText(element), nil
			}
		default:
			{
				return htmlquery.SelectAttr(element, attribute), nil
			}
		}
	}
}
