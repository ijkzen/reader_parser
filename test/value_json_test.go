package test

import (
	"os"
	gson "reader_parse/analyse/json"
	"reflect"
	"testing"
)

func TestJsonValue(t *testing.T) {
	jsonBytes, err := os.ReadFile("../html/json_test.json")

	if err != nil {
		t.Errorf("can not find json_test.json")
		return
	}

	json := string(jsonBytes)

	analyse := gson.JsonAnalyse{
		Json: json,
	}

	elements, err := analyse.GetElements("$..books[0].author")

	if err == nil || reflect.TypeOf(elements).Kind() == reflect.Slice {
		if reflect.ValueOf(elements).Len() != 1 {
			t.Errorf("want 1 but %d", reflect.ValueOf(elements).Len())
		} else {
			element := elements[0]

			if element != "囧囧有妖" {
				t.Errorf("want 囧囧有妖 but %s", element)
			}
		}
	}

	elements, err = analyse.GetElements("$..books[0]..lastChapter")

	if err == nil || reflect.TypeOf(elements).Kind() == reflect.Slice {
		if reflect.ValueOf(elements).Len() != 1 {
			t.Errorf("want 1 but %d", reflect.ValueOf(elements).Len())
		} else {
			element := elements[0]

			if element != "2 【新书《月亮在怀里》】" {
				t.Errorf("want `2 【新书《月亮在怀里》】` but %s", element)
			}
		}
	}

}
