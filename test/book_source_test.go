package test

import (
	"os"
	"testing"

	"reader_parse/booksource"
)

func TestBookSourceValid(t *testing.T) {
	jsonBytes, err := os.ReadFile("../html/bookSource.json")

	if err != nil {
		t.Errorf(err.Error())
	} else {
		validList, invalidList, err := booksource.GetBookSourceList(jsonBytes)
		if err != nil {
			t.Errorf(err.Error())
		} else {
			t.Logf("\nvalidList size: %d, invalidList size: %d", len(validList), len(invalidList))
		}
	}
}
