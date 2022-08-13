package css

import (
	"bytes"
	"fmt"
	"github.com/samccone/css"
	"golang.org/x/net/html"
	"strings"
)

type CssAnalyse struct {
	Doc *html.Node
}

func (analyse *CssAnalyse) GetElements(rule string) ([]*html.Node, error) {
	selector, err := css.Compile(rule)
	if err != nil {
		return make([]*html.Node, 0), err
	}

	return selector.Select(analyse.Doc), nil
}

func (analyse *CssAnalyse) GetValue(rule string) (string, error) {
	if !strings.Contains(rule, "@") {
		return "", fmt.Errorf("rule %s is invalid", rule)
	}

	atIndex := strings.Index(rule, "@")
	if atIndex > 0 {
		preRule := rule[0:atIndex]
		attribute := rule[atIndex+1:]

		elements, err := analyse.GetElements(preRule)
		if err != nil {
			return "", err
		}

		elementSize := len(elements)

		if elementSize != 1 {
			return "", fmt.Errorf("elements size is not 1 but %d", elementSize)
		} else {
			targetElement := elements[0]

			switch attribute {
			case "text", "textNodes", "ownText":
				{
					return FullText(targetElement), nil
				}
			default:
				{
					attributes := getKeyValue(targetElement.Attr)
					value := attributes[attribute]

					if value == "" {
						return "", fmt.Errorf("value is empty")
					} else {
						return value, nil
					}
				}
			}
		}

	} else {
		return "", fmt.Errorf("@ index is %d", atIndex)
	}
}

func GetCssAnalyseFromString(str string) (CssAnalyse, error) {
	node, err := html.Parse(strings.NewReader(str))
	if err != nil {
		return CssAnalyse{}, err
	}

	analyse := CssAnalyse{
		Doc: node,
	}

	return analyse, nil
}

func getKeyValue(attributes []html.Attribute) map[string]string {
	var keyvalues = make(map[string]string)
	for i := 0; i < len(attributes); i++ {
		_, exists := keyvalues[attributes[i].Key]
		if exists == false {
			keyvalues[attributes[i].Key] = attributes[i].Val
		}
	}
	return keyvalues
}

func FullText(root *html.Node) string {
	var buf bytes.Buffer

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n == nil {
			return
		}
		if n.Type == html.TextNode {
			buf.WriteString(n.Data)
		}
		if n.Type == html.ElementNode {
			f(n.FirstChild)
		}
		if n.NextSibling != nil {
			f(n.NextSibling)
		}
	}

	f(root.FirstChild)

	return buf.String()
}
