package service

import (
	l0wb "L0WB"
	"L0WB/pkg/repository"
)

type SaveOrderData interface {
	SaveOrderData(l0wb.Order) error
}

type ReadOrdersData interface {
	ReadAllOrdersData() map[string]l0wb.Order
	ReadOrderData(orderUid string) l0wb.Order
}

type ParseJSON interface {
	ParseJSON([]byte) (l0wb.Order, error)
}

type Service struct {
	Cache map[string]l0wb.Order
	SaveOrderData
	ReadOrdersData
	ParseJSON
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Cache:          make(map[string]l0wb.Order),
		SaveOrderData:  NewSaveOrderDataService(repos.SaveOrderData),
		ReadOrdersData: NewReadOrdersDataService(repos.ReadOrdersData),
		ParseJSON:      NewParseJson(),
	}
}
