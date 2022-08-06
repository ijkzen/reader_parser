package main

import (
	"fmt"
	"github.com/anaskhan96/soup"
	"io/ioutil"
	"reader_parse/analyse/soup_like"
)

func main() {
	htmlBytes, err := ioutil.ReadFile("./html/起点-搜索.html")

	if err != nil {
		fmt.Println(err)
		return
	}

	html := string(htmlBytes)

	doc := soup.HTMLParse(html)

	elements := doc.Find("", "class", "book-img-text").FindAll("li")

	for _, element := range elements {
		analyse := soup_like.SoupAnalyse{
			Element: &element,
		}
		title, err1 := analyse.GetValue("tag.h4@a@text||tag.a.1@text")
		if err1 == nil {
			fmt.Println("书名：", title)
		} else {
			fmt.Println("书名规则错误，" + err1.Error())
		}
		author, err2 := analyse.GetValue("class.author@class.name.0@text||tag.a.2@text||tag.span@text")
		if err2 == nil {
			fmt.Println("作者：", author)
		} else {
			fmt.Println("作者规则错误，" + err2.Error())
		}
		introduction, err3 := analyse.GetValue("class.intro@textNodes")
		if err3 == nil {
			fmt.Println("简介：", introduction)
		} else {
			fmt.Println("简介规则错误，" + err3.Error())
		}
		coverImgUrl, err4 := analyse.GetValue("class.book-img-box@tag.img@src")
		if err4 == nil {
			fmt.Println("封面：", coverImgUrl)
		}
		fmt.Println()
	}

}
