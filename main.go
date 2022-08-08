package main

import (
	"fmt"
	"github.com/anaskhan96/soup"
	"io/ioutil"
	"reader_parse/analyse/gsoup"
)

func main() {
	htmlBytes, err := ioutil.ReadFile("./html/红袖-搜索.html")

	if err != nil {
		fmt.Println(err)
		return
	}

	html := string(htmlBytes)

	doc := soup.HTMLParse(html)

	rootAnalyse := gsoup.SoupAnalyse{
		Element: &doc,
	}

	elements, err := rootAnalyse.GetElements("id.result-list@tag.li")
	if err != nil {
		return
	}

	for _, element := range elements {
		analyse := gsoup.SoupAnalyse{
			Element: &element,
		}
		title, err1 := analyse.GetValue("class.book-mid-info@tag.a.0@text")
		if err1 == nil {
			fmt.Println("书名：", title)
		} else {
			fmt.Println("书名规则错误，" + err1.Error())
		}
		author, err2 := analyse.GetValue("class.author@tag.a.0@text")
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
		lastChapter, err5 := analyse.GetValue("class.update@tag.a@text##最新更新")
		if err5 == nil {
			fmt.Println("最新章节：", lastChapter)
		}
		fmt.Println()
	}

}
