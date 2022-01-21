package common

import "github.com/alexZaicev/go-vtd-xml/vtdxml/erroring"

type List interface {
	Get(index int) ([]interface{}, error)
	Set(index int, value []interface{}) error
	Add(value []interface{})
	Clear()
	Size() int
}

type ArrayList struct {
	oa   [][]interface{}
	size int
}

func NewArrayList() ArrayList {
	return ArrayList{
		oa: [][]interface{}{},
	}
}

func (al *ArrayList) Get(index int) ([]interface{}, error) {
	if index < 0 || index >= al.size {
		return nil, erroring.NewInvalidArgumentError("index", erroring.IndexOutOfRange, nil)
	}
	return al.oa[index], nil
}

func (al *ArrayList) Set(index int, value []interface{}) error {
	if index < 0 || index >= al.size {
		return erroring.NewInvalidArgumentError("index", erroring.IndexOutOfRange, nil)
	}
	al.oa[index] = value
	return nil
}

func (al *ArrayList) Add(value []interface{}) {
	al.oa = append(al.oa, value)
	al.size++
}

func (al *ArrayList) Clear() {
	al.oa = nil
}

func (al *ArrayList) Size() int {
	return al.size
}
