package repository

import (
	l0wb "L0WB"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type SaveDataPostgres struct {
	db *sqlx.DB
}

func NewSaveOrderDataPostgres(db *sqlx.DB) *SaveDataPostgres {
	return &SaveDataPostgres{db}
}

func (r *SaveDataPostgres) SaveOrderData(data l0wb.Order) error {
	var orderUid string
	query := fmt.Sprintf("INSERT INTO %s (order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING order_uid", ordersTable)

	row := r.db.QueryRow(query, data.Order_uid, data.Track_number, data.Entry, data.Locale, data.Internal_signature, data.Customer_id,
		data.Delivery_service, data.Shardkey, data.Sm_id, data.Date_created, data.Oof_shard)

	if err := row.Scan(&orderUid); err != nil {
		logrus.Fatalf("failed to save to %s table: %s", ordersTable, err.Error())
		return err
	}

	query = fmt.Sprintf("INSERT INTO %s (order_uid, name, phone, zip, city, address, region, email) values ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING order_uid", deliveryTable)

	row = r.db.QueryRow(query, data.Order_uid, data.Delivery.Name, data.Delivery.Phone, data.Delivery.Zip, data.Delivery.City, data.Delivery.Address,
		data.Delivery.Region, data.Delivery.Email)

	if err := row.Scan(&orderUid); err != nil {
		logrus.Fatalf("failed to save to %s table: %s", deliveryTable, err.Error())
		return err
	}

	query = fmt.Sprintf("INSERT INTO %s (order_uid, transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING order_uid", paymentTable)

	row = r.db.QueryRow(query, data.Order_uid, data.Payment.Transaction, data.Payment.Request_id, data.Payment.Currency, data.Payment.Provider,
		data.Payment.Amount, data.Payment.Payment_dt, data.Payment.Bank, data.Payment.Delivery_cost, data.Payment.Goods_total, data.Payment.Custom_fee)

	if err := row.Scan(&orderUid); err != nil {
		logrus.Fatalf("failed to save to %s table: %s", paymentTable, err.Error())
		return err
	}

	for _, item := range data.Items {

		query = fmt.Sprintf("INSERT INTO %s (order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status) values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING order_uid", itemsTable)

		row = r.db.QueryRow(query, data.Order_uid, item.Chrt_id, item.Track_number, item.Price, item.Rid, item.Name, item.Sale,
			item.Size, item.Total_price, item.Nm_id, item.Brand, item.Status)

		if err := row.Scan(&orderUid); err != nil {
			logrus.Fatalf("failed to save to %s table: %s", itemsTable, err.Error())
			return err
		}
	}

	return nil
}
