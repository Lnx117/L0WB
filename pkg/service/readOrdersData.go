package service

import (
	l0wb "L0WB"
	"L0WB/pkg/repository"
)

type ReadOrdersDataService struct {
	repo repository.ReadOrdersData
}

func NewReadOrdersDataService(repo repository.ReadOrdersData) *ReadOrdersDataService {
	return &ReadOrdersDataService{repo: repo}
}

func (s *ReadOrdersDataService) ReadAllOrdersData() map[string]l0wb.Order {
	ordersMap := s.repo.ReadAllOrdersData()
	return ordersMap
}

func (s *ReadOrdersDataService) ReadOrderData(orderUid string) l0wb.Order {
	order := s.repo.ReadOrderData(orderUid)
	return order
}
