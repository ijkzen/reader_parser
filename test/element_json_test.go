package test

import (
	"os"
	gson "reader_parse/analyse/json"
	"reflect"
	"testing"
)

func TestJsonElement(t *testing.T) {
	jsonBytes, err := os.ReadFile("../html/json_test.json")

	if err != nil {
		t.Errorf("can not find json_test.json")
		return
	}

	json := string(jsonBytes)

	analyse := gson.JsonAnalyse{
		Json: json,
	}

	elements, err := analyse.GetElements("$..books[*]")

	if err == nil || reflect.TypeOf(elements).Kind() == reflect.Slice {
		if reflect.ValueOf(elements).Len() != 100 {
			t.Errorf("want 100 but %d", reflect.ValueOf(elements).Len())
		}
	}
}
