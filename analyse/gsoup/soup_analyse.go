package gsoup

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/anaskhan96/soup"
)

type SoupAnalyse struct {
	Element *soup.Root
}

func (analyse SoupAnalyse) GetValue(rule string) (string, error) {
	var ruleList []string = getRuleList(rule)

	if len(ruleList) > 0 {
		for _, subRule := range ruleList {
			ruleParts := strings.Split(subRule, "@")
			partsSize := len(ruleParts)
			elementList := emptyElementSlice()
			if partsSize < 2 {
				return "", fmt.Errorf("rule %s is invalid", rule)
			} else {
				elementListRule := ""
				for i := 0; i < partsSize-1; i++ {
					if i == partsSize-2 {
						elementListRule += ruleParts[i]
					} else {
						elementListRule += ruleParts[i]
						elementListRule += "@"
					}
				}

				elements, err := analyse.GetElements(elementListRule)
				if err != nil {
					continue
				} else {
					elementList = elements
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

func (analyse SoupAnalyse) GetElements(rule string) ([]soup.Root, error) {
	var ruleList []string = getRuleList(rule)

	if len(ruleList) > 0 {
		for _, subRule := range ruleList {
			ruleParts := strings.Split(subRule, "@")
			partsSize := len(ruleParts)
			var elementList = emptyElementSlice()
			elementList = append(elementList, *analyse.Element)
			for ruleIndex := 0; ruleIndex < partsSize; ruleIndex++ {
				if len(elementList) == 0 {
					break
				} else {
					parentSize := len(elementList)
					tmpList := emptyElementSlice()
					for parentIndex := 0; parentIndex < parentSize; parentIndex++ {
						elements, err := getElements(&elementList[parentIndex], ruleParts[ruleIndex])
						if err == nil {
							tmpList = append(tmpList, elements...)
						}
					}

					elementList = tmpList
				}
			}

			if len(elementList) == 0 {
				continue
			} else {
				return elementList, nil
			}
		}

		return emptyElementSlice(), fmt.Errorf("rule %s is invalid", rule)
	} else {
		return emptyElementSlice(), fmt.Errorf("rule %s is invalid", rule)
	}
}

func getRuleList(rule string) []string {
	var ruleList []string
	if strings.Contains(rule, "||") {
		ruleList = strings.Split(rule, "||")
	} else if strings.Contains(rule, "&&") {
		ruleList = strings.Split(rule, "&&")
	} else {
		ruleList = append(ruleList, rule)
	}

	return ruleList
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
 * [1:9]
 * [1:9:2]
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
				childrenSize := len(element.Children())
				finalPosition, err1 := getSliceIndex(childrenSize, index)
				if err1 != nil {
					return emptyElementSlice(), err1
				} else {
					slice := emptyElementSlice()
					slice = append(slice, element.Children()[finalPosition])
					return slice, nil
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

		children := element.Children()
		startIndex := strings.Index(rule, "[")
		endIndex := strings.Index(rule, "]") + 1

		part3 := rule[startIndex:endIndex]
		if strings.Contains(part3, ":") {
			tmpList := emptyElementSlice()
			for _, child := range children {
				slice, err := getElementsBySlice(child.Children(), part3)
				if err == nil {
					tmpList = append(tmpList, slice...)
				}
			}
			return tmpList, nil
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

// [1:9], [1:9:2]
func getElementsBySlice(children []soup.Root, sliceRule string) ([]soup.Root, error) {
	if strings.Contains(sliceRule, "[") && strings.Contains(sliceRule, "]") && strings.Contains(sliceRule, ":") {
		tmpStr := strings.Replace(sliceRule, "[", "", -1)
		tmpStr = strings.Replace(tmpStr, "]", "", -1)

		sliceParts := strings.Split(tmpStr, ":")
		childrenSize := len(children)
		switch len(sliceParts) {
		case 1:
			{
				start, err := getSliceIndex(childrenSize, sliceParts[0])
				if err != nil {
					return emptyElementSlice(), err
				} else {
					return children[start:], nil
				}
			}
		case 2:
			{

				start, err1 := getSliceIndex(childrenSize, sliceParts[0])
				if err1 != nil {
					return emptyElementSlice(), err1
				} else {
					end, err2 := getSliceIndex(childrenSize, sliceParts[1])
					if err2 != nil {
						return emptyElementSlice(), err2
					} else {
						return children[start : end+1], nil
					}
				}
			}
		case 3:
			{
				start, err1 := getSliceIndex(childrenSize, sliceParts[0])
				end, err2 := getSliceIndex(childrenSize, sliceParts[1])
				step, err3 := getSliceIndex(childrenSize, sliceParts[2])

				if err1 != nil || err2 != nil || err3 != nil {
					return emptyElementSlice(), fmt.Errorf("sliceRule %s is invalid", sliceRule)
				} else {
					tmpSlice := emptyElementSlice()
					for i := start; i <= end; {
						tmpSlice = append(tmpSlice, children[i])
						i += step
					}
					return tmpSlice, nil
				}
			}
		default:
			{
				return emptyElementSlice(), fmt.Errorf("sliceRule %s is invalid", sliceRule)
			}
		}
	} else {
		return emptyElementSlice(), fmt.Errorf("sliceRule %s is invalid", sliceRule)
	}
}

func getSliceIndex(size int, index string) (int, error) {
	result, err := strconv.Atoi(index)
	if err != nil {
		return 0, err
	} else {
		if result < 0 {
			if result+size >= 0 {
				return result + size, nil
			} else {
				return 0, fmt.Errorf("index %d is invalid", result)
			}
		} else {
			if result < size {
				return result, nil
			} else {
				return 0, fmt.Errorf("index %d is invalid", result)
			}
		}
	}
}

func getElementsByPosition(children []soup.Root, positionRule string) ([]soup.Root, error) {
	tmpStr := strings.Replace(positionRule, "[", "", -1)
	tmpStr = strings.Replace(tmpStr, "]", "", -1)

	var positionList []string
	if strings.Contains(tmpStr, ":") {
		positionList = strings.Split(tmpStr, ":")
	} else if strings.Contains(tmpStr, ",") {
		positionList = strings.Split(tmpStr, ",")
	} else {
		positionList = append(positionList, tmpStr)
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
		attributeParts := strings.Split(rule, "##")
		attribute := attributeParts[0]

		for _, element := range elementList {
			switch attribute {
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
			if len(attributeParts) > 1 {
				result = replace(result, attributeParts[1:])
			}

			return result, nil
		}
	}
}

func replace(result string, replaceList []string) string {
	if len(replaceList) == 0 {
		return result
	}

	tmp := result
	for i := 0; i < len(replaceList); i++ {
		if i%2 == 0 {
			regexRule := replaceList[i]
			replaceText := ""
			if i+1 < len(replaceList) {
				i++
				replaceText = replaceList[i]
			}

			regex, err := regexp.Compile(regexRule)
			if err != nil {
				break
			} else {
				tmp = regex.ReplaceAllString(tmp, replaceText)
			}
		}
	}

	return tmp
}
