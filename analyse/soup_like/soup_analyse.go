package soup_like

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/anaskhan96/soup"
)

type SoupAnalyse struct {
	Element *soup.Root
}

//todo 替换过滤
func (analyse SoupAnalyse) GetValue(rule string) (string, error) {
	var ruleList []string
	if strings.Contains(rule, "||") {
		ruleList = strings.Split(rule, "||")
	} else if strings.Contains(rule, "&&") {
		ruleList = strings.Split(rule, "&&")
	} else {
		ruleList = append(ruleList, rule)
	}

	if len(ruleList) > 0 {
		for _, subRule := range ruleList {
			ruleParts := strings.Split(subRule, "@")
			partsSize := len(ruleParts)
			var elementList = emptyElementSlice()
			elementList = append(elementList, *analyse.Element)
			for i := 0; i < partsSize-1; i++ {
				if len(elementList) == 0 {
					break
				} else {
					elements, err := getElements(&elementList[0], ruleParts[0])
					if err != nil {
						return "", err
					} else {
						elementList = elements
					}
				}
			}

			if len(elementList) == 0 {
				continue
			} else {
				return getAttribute(elementList, ruleParts[partsSize-1])
			}
		}
		return "", fmt.Errorf("rule %s is invalid", rule)
	} else {
		return "", fmt.Errorf("rule %s is invalid", rule)
	}
}

func emptyElementSlice() []soup.Root {
	return make([]soup.Root, 0)
}

/*
 * 根据规则获取节点列表，这里的规则有如下形式
 * tag.span.-1
 * class.odd.0
 * class.odd
 * tag.li!0:1:2:3:4:5
 * tag.li.0:1
 * .1
 * [0,1,2,!3]
 * [1:9] // 暂不适配
 * [1:9:2] // 暂不适配
 * children[1]
 */
func getElements(element *soup.Root, rule string) ([]soup.Root, error) {
	if strings.Contains(rule, ".") {
		parts := strings.Split(rule, ".")
		switch len(parts) {
		// 适配 .1 规则
		case 1:
			{
				index := parts[0]
				position, err := strconv.Atoi(index)
				if err != nil {
					return make([]soup.Root, 0), err
				} else {
					childrenSize := len(element.Children())
					if position < childrenSize {
						slice := emptyElementSlice()
						slice = append(slice, *&element.Children()[position])
						return slice, nil
					} else {
						return emptyElementSlice(), fmt.Errorf("child size %d, but position is %d", childrenSize, position)
					}
				}

			}
		// 适配 class.odd，tag.li!0:1:2:3:4:5 规则
		case 2:
			{
				var part1, part2, part3 string
				part1 = parts[0]
				positionStartIndex := strings.Index(parts[1], "!")
				if positionStartIndex >= 0 {
					// 适配 tag.li!0:1:2:3:4:5
					part2 = parts[1][0:positionStartIndex]
					part3 = parts[1][positionStartIndex:]
					children, err := getChildren(element, part1, part2)
					if err != nil {
						return emptyElementSlice(), err
					} else {
						return getElementsByPosition(children, part3)
					}
				} else {
					// 适配 class.odd
					part2 = parts[1]
					return getChildren(element, part1, part2)
				}
			}
		case 3:
			{
				part1 := parts[0]
				part2 := parts[1]
				part3 := parts[2]
				children, err := getChildren(element, part1, part2)
				if err != nil {
					return emptyElementSlice(), err
				} else {
					return getElementsByPosition(children, part3)
				}
			}
		default:
			{
				return make([]soup.Root, 0), fmt.Errorf("rule %s is invalid", rule)
			}
		}

	} else if strings.Contains(rule, "[") && strings.Contains(rule, "]") {
		//todo
		children := element.Children()
		startIndex := strings.Index(rule, "[")
		endIndex := strings.Index(rule, "]")

		part3 := rule[startIndex:endIndex]
		if strings.Contains(part3, ":") {
			return emptyElementSlice(), fmt.Errorf("暂不支持区间写法")
		}

		return getElementsByPosition(children, part3)
	} else {
		return emptyElementSlice(), fmt.Errorf("rule %s is invalid", rule)
	}
}

func getChildren(element *soup.Root, part1 string, part2 string) ([]soup.Root, error) {
	switch part1 {
	case "class":
		{
			slice := element.FindAll("", "class", part2)
			return slice, nil
		}
	case "tag":
		{
			slice := element.FindAll(part2)
			return slice, nil
		}
	case "id":
		{
			slice := element.FindAll("", "id", part2)
			return slice, nil
		}
	default:
		{
			return make([]soup.Root, 0), fmt.Errorf("rule %s is invalid", part1+"."+part2)
		}
	}
}

func getElementsByPosition(children []soup.Root, positionRule string) ([]soup.Root, error) {
	var positionList []string
	if strings.Contains(positionRule, ":") {
		positionList = strings.Split(positionRule, ":")
	} else if strings.Contains(positionRule, ",") {
		positionList = strings.Split(positionRule, ",")
	} else {
		positionList = append(positionList, positionRule)
	}

	positionSize := len(positionList)
	childrenSize := len(children)
	slice := emptyElementSlice()
	for i := 0; i < positionSize; i++ {
		singlePosition := positionList[i]
		if realPosition, ok := getRealPosition(singlePosition, childrenSize); ok {
			slice = append(slice, children[realPosition])
		}
	}
	if len(slice) == 0 {
		return slice, fmt.Errorf("positionRule %s is invalid", positionRule)
	}

	return slice, nil
}

/* 获取指定的下标
 * -1 -> listSize -1
 * !x -> -1
 * 1 -> 1
 */
func getRealPosition(position string, listSize int) (int, bool) {
	if strings.Contains(position, "!") {
		index, err := strconv.Atoi(position[1:])
		if err != nil {
			return -1, false
		} else {
			return index, false
		}
	} else {
		index, err := strconv.Atoi(position)
		if err != nil {
			return -1, false
		} else {
			if index < 0 {
				finalIndex := listSize + index
				if finalIndex < 0 {
					return -1, false
				} else {
					return finalIndex, true
				}
			} else {
				if index < listSize {
					return index, true
				} else {
					return -1, false
				}
			}
		}
	}
}

func getAttribute(elementList []soup.Root, rule string) (string, error) {
	var result = ""
	if len(elementList) == 0 {
		return "", fmt.Errorf("elementList size is 0")
	} else {
		for _, element := range elementList {
			switch rule {
			case "text", "textNodes", "ownText":
				{
					result += element.FullText()
				}
			default:
				{
					value := element.Attrs()[rule]

					if value != "" {
						result += value
					}
				}
			}
		}

		if result == "" {
			return result, fmt.Errorf("attributes is not exists")
		} else {
			return result, nil
		}
	}
}
