package json

import (
	"github.com/ohler55/ojg/jp"
	"github.com/ohler55/ojg/oj"
)

type JsonAnalyse struct {
	Json string
}

func (analyse *JsonAnalyse) GetElements(rule string) ([]interface{}, error) {
	obj, err := oj.ParseString(analyse.Json)
	if err != nil {
		return nil, err
	}

	x, err := jp.ParseString(rule)
	if err != nil {
		return nil, err
	}

	return x.Get(obj), nil
}

func (analyse *JsonAnalyse) GetValue(rule string) ([]interface{}, error) {
	obj, err := oj.ParseString(analyse.Json)
	if err != nil {
		return nil, err
	}

	x, err := jp.ParseString(rule)
	if err != nil {
		return nil, err
	}

	return x.Get(obj), nil
}
