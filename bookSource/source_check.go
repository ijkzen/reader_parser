package booksource

import (
	"encoding/json"
	"strings"
)

type RuleSearch struct {
	Author      string
	BookList    string
	BookUrl     string
	CoverUrl    string
	Intro       string
	Kind        string
	LastChapter string
	Name        string
	WordCount   string
}

func isRuleValid(rule string) bool {
	return !strings.Contains(rule, "@js:") && !strings.Contains(rule, "<js>") && !strings.Contains(rule, "webView")
}

func (rule RuleSearch) IsValid() bool {
	return isRuleValid(rule.Author) && isRuleValid(rule.BookList) &&
		isRuleValid(rule.BookUrl) && isRuleValid(rule.CoverUrl) && isRuleValid(rule.Intro) &&
		isRuleValid(rule.Kind) && isRuleValid(rule.LastChapter) && isRuleValid(rule.Name) &&
		isRuleValid(rule.WordCount)
}

type RuleBookInfo struct {
	Author      string
	Intro       string
	Kind        string
	LastChapter string
	Name        string
	TocUrl      string
	WordCount   string
}

func (rule RuleBookInfo) IsValid() bool {
	return isRuleValid(rule.Author) && isRuleValid(rule.Intro) && isRuleValid(rule.Kind) &&
		isRuleValid(rule.LastChapter) && isRuleValid(rule.Name) && isRuleValid(rule.TocUrl) &&
		isRuleValid(rule.WordCount)
}

type RuleToc struct {
	ChapterList string
	ChapterName string
	ChapterUrl  string
	UpdateTime  string
}

func (rule RuleToc) isValid() bool {
	return isRuleValid(rule.ChapterList) && isRuleValid(rule.ChapterName) && isRuleValid(rule.ChapterUrl) &&
		isRuleValid(rule.UpdateTime)
}

type RuleContent struct {
	Content        string
	NextContentUrl string
	ReplaceRegex   string
}

func (rule RuleContent) IsValid() bool {
	return isRuleValid(rule.Content) && isRuleValid(rule.NextContentUrl) && isRuleValid(rule.ReplaceRegex)
}

type BookSource struct {
	BookSourceGroup string
	BookSourceName  string
	BookSourceUrl   string
	Header          string
	RuleSearch      RuleSearch
	RuleBookInfo    RuleBookInfo
	RuleToc         RuleToc
	RuleContent     RuleContent
}

func (source BookSource) IsValid() bool {
	return source.RuleSearch.IsValid() && source.RuleBookInfo.IsValid() && source.RuleToc.isValid() &&
		source.RuleContent.IsValid()
}

func GetBookSourceList(jsonData []byte) ([]BookSource, []BookSource, error) {
	validList := make([]BookSource, 0)
	invalidList := make([]BookSource, 0)

	var originList []BookSource

	err := json.Unmarshal(jsonData, &originList)

	if err != nil {
		return validList, invalidList, err
	}

	for _, bookSource := range originList {
		if bookSource.IsValid() {
			validList = append(validList, bookSource)
		} else {
			invalidList = append(invalidList, bookSource)
		}
	}

	return validList, invalidList, nil
}
