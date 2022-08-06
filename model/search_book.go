package model

/* 搜索时的搜索结果 */
type SearchBook struct {
	// 书名
	Name string
	// 作者名
	Author string
	// 书籍类型
	Kind string
	// 字数
	WordCount string
	// 最新章节名
	LastChapter string
	// 简介
	Introduction string
	// 封面链接
	CoverImgUrl string
	// 书籍详情页链接
	DetailUrl string
}
