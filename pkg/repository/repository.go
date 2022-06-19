package repository

import (
	l0wb "L0WB"

	"github.com/jmoiron/sqlx"
)

type SaveOrderData interface {
	SaveOrderData(data l0wb.Order) error
}

type ReadOrdersData interface {
	ReadAllOrdersData() map[string]l0wb.Order
	ReadOrderData(orderUid string) l0wb.Order
}

type Repository struct {
	SaveOrderData
	ReadOrdersData
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		SaveOrderData:  NewSaveOrderDataPostgres(db),
		ReadOrdersData: NewReadDataPostgres(db),
	}
}
