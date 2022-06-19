package service

import (
	l0wb "L0WB"
	"L0WB/pkg/repository"
)

type SaveOrderDataService struct {
	repo repository.SaveOrderData
}

func NewSaveOrderDataService(repo repository.SaveOrderData) *SaveOrderDataService {
	return &SaveOrderDataService{repo: repo}
}

func (s *SaveOrderDataService) SaveOrderData(data l0wb.Order) error {
	if err := s.repo.SaveOrderData(data); err != nil {
		return err
	}
	return nil
}
