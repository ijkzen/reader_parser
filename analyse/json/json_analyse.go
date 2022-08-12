package json

import (
	"github.com/larstos/jsonpath_go"
)

type JsonAnalyse struct {
	Json interface{}
}

func (analyse *JsonAnalyse) GetElements(rule string) (interface{}, error) {

	return jsonpath_go.Lookup(analyse.Json, rule)
}

func (analyse *JsonAnalyse) GetValue(rule string) (interface{}, error) {
	return jsonpath_go.Lookup(analyse.Json, rule)
}