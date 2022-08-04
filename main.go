package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/anaskhan96/soup"
)

type SimpleBook struct {
	Title        string
	Author       string
	Introduction string
	CoverImgUrl  string
}

func main() {
	htmlBytes, err := ioutil.ReadFile("./html/红袖-搜索.html")

	if err != nil {
		fmt.Println(err)
		return
	}

	html := string(htmlBytes)

	doc := soup.HTMLParse(html)
	// list id.result-list@tag.li
	// title class.book-mid-info@tag.a.0@text
	// author class.author@tag.a.0@text
	// introdution class.intro@textNodes
	// coverImgUrl class.book-img-box@tag.img@src

	elements := doc.Find("", "id", "result-list").FindAll("li")

	for _, element := range elements {
		title, err1 := getPropertyFromElement("class.book-mid-info@tag.a.0@text", &element)
		if err1 == nil {
			fmt.Println("书名：", title)
		}
		author, err2 := getPropertyFromElement("class.author@tag.a.0@text", &element)
		if err2 == nil {
			fmt.Println("作者：", author)
		}
		introduction, err3 := getPropertyFromElement("class.intro@textNodes", &element)
		if err3 == nil {
			fmt.Println("简介：", introduction)
		}
		coverImgUrl, err4 := getPropertyFromElement("class.book-img-box@tag.img@src", &element)
		if err4 == nil {
			fmt.Println("封面：", coverImgUrl)
		}
		fmt.Println()
	}

}

func getPropertyFromElement(rule string, element *soup.Root) (string, error) {
	ruleList := strings.Split(rule, "@")
	var targetElement = element
	var attribute = ""
	if len(ruleList) > 0 {
		for i := 0; i < len(ruleList); i++ {
			subRule := ruleList[i]
			slice, err := getPropertyWithSubRule(subRule, targetElement)
			if err != nil {
				if err.Error() == "element is result element" {
					attribute = subRule
					break
				}

				return "", err
			} else {
				if len(slice) > 0 {
					targetElement = &slice[0]
				}
			}
		}

		switch attribute {
		case "text", "textNodes", "ownText":
			{
				return targetElement.FullText(), nil
			}

		default:
			{
				value := targetElement.Attrs()[attribute]

				if value == "" {
					return "", errors.New("do not has attribute")
				} else {
					return value, nil
				}
			}

		}

	} else {
		return "", errors.New("fail to analyse rule")
	}
}

func getPropertyWithSubRule(subRule string, element *soup.Root) ([]soup.Root, error) {
	subRuleParts := strings.Split(subRule, ".")
	switch len(subRuleParts) {
	case 1:
		{
			return make([]soup.Root, 0), errors.New("element is result element")
		}
	case 2:
		{
			subRuleParts = append(subRuleParts, "")
			return getPropertyWithSubRuleParts(subRuleParts, element)
		}
	case 3:
		{
			return getPropertyWithSubRuleParts(subRuleParts, element)
		}
	default:
		{
			return make([]soup.Root, 0), errors.New("fail to analyse subRule")
		}
	}
}

func getPropertyWithSubRuleParts(subRuleParts []string, element *soup.Root) ([]soup.Root, error) {
	if len(subRuleParts) != 3 {
		return make([]soup.Root, 0), errors.New("fail to analyse subRuleParts")
	} else {
		pickElement := func(parts []soup.Root, part3 string) ([]soup.Root, error) {
			if part3 == "" {
				return parts, nil
			} else {
				position, err := strconv.Atoi(part3)
				if err != nil {
					return make([]soup.Root, 0), errors.New("fail to parse position")
				} else {
					if len(parts)-1 >= position {
						slice := make([]soup.Root, 0)
						slice = append(slice, parts[position])
						return slice, nil
					} else {
						return make([]soup.Root, 0), errors.New("position is outOfIndex")
					}
				}
			}
		}

		part1 := subRuleParts[0]
		part2 := subRuleParts[1]
		part3 := subRuleParts[2]

		switch part1 {
		case "class":
			{
				parts := element.FindAll("", "class", part2)
				return pickElement(parts, part3)
			}
		case "tag":
			{
				parts := element.FindAll(part2)
				return pickElement(parts, part3)
			}
		case "id":
			{
				parts := element.FindAll("", "id", part2)
				return pickElement(parts, part3)
			}
		default:
			{
				return make([]soup.Root, 0), errors.New("unknown part1")
			}
		}
	}
}
