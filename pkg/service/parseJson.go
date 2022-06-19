package service

import (
	l0wb "L0WB"
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
)

type ParseJson struct {
}

func NewParseJson() *ParseJson {
	return &ParseJson{}
}

func (s *ParseJson) ParseJSON(data []byte) (l0wb.Order, error) {
	var order l0wb.Order
	var err error

	err = json.Unmarshal([]byte(data), &order)
	if err != nil {
		logrus.Fatalf("ERROR: %s", err.Error())
		return order, err
	}
	return order, nil
}

func (s *ParseJson) Hui(data string) {
	fmt.Println(data)
}
