package repository

import (
	l0wb "L0WB"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type ReadDataPostgres struct {
	db *sqlx.DB
}

func NewReadDataPostgres(db *sqlx.DB) *ReadDataPostgres {
	return &ReadDataPostgres{db}
}

func (r *ReadDataPostgres) ReadAllOrdersData() map[string]l0wb.Order {
	query := fmt.Sprintf("SELECT order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard FROM %s", ordersTable)

	row, err := r.db.Queryx(query)

	if err != nil {
		logrus.Info("failed to read data from %s table: %s", ordersTable, err.Error())
	}

	ordersMap := make(map[string]l0wb.Order)

	for row.Next() {
		var p l0wb.Order
		err := row.StructScan(&p)
		if err != nil {
			logrus.Info("failed to scan data from %s table: %s", ordersTable, err.Error())
		}

		query = fmt.Sprintf("SELECT name, phone, zip, city, address, region, email FROM %s WHERE order_uid=$1", deliveryTable)

		row := r.db.QueryRowx(query, p.Order_uid)

		err = row.StructScan(&p.Delivery)
		if err != nil {
			logrus.Info("failed to scan data from %s table: %s", deliveryTable, err.Error())
		}

		query = fmt.Sprintf("SELECT transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee FROM %s WHERE order_uid=$1", paymentTable)

		row = r.db.QueryRowx(query, p.Order_uid)

		err = row.StructScan(&p.Payment)
		if err != nil {
			logrus.Info("failed to scan data from %s table: %s", paymentTable, err.Error())
		}

		query := fmt.Sprintf("SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status FROM %s WHERE order_uid=$1", itemsTable)

		newRow, _ := r.db.Queryx(query, p.Order_uid)

		var item l0wb.Item
		count := 0
		for newRow.Next() {
			err := newRow.StructScan(&item)
			if err != nil {
				logrus.Info("failed to scan data from %s table: %s", itemsTable, err.Error())
			}
			p.Items = append(p.Items, item)
			count++
		}

		ordersMap[p.Order_uid] = p
	}

	return ordersMap
}

func (r *ReadDataPostgres) ReadOrderData(orderUid string) l0wb.Order {
	var p l0wb.Order

	query := fmt.Sprintf("SELECT order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard FROM %s WHERE order_uid=$1", ordersTable)

	row := r.db.QueryRowx(query, orderUid)

	err := row.StructScan(&p)
	if err != nil {
		logrus.Info("failed to scan data from %s table: %s", ordersTable, err.Error())
	}

	query = fmt.Sprintf("SELECT name, phone, zip, city, address, region, email FROM %s WHERE order_uid=$1", deliveryTable)

	row = r.db.QueryRowx(query, orderUid)

	err = row.StructScan(&p.Delivery)
	if err != nil {
		logrus.Info("failed to scan data from %s table: %s", deliveryTable, err.Error())
	}

	query = fmt.Sprintf("SELECT transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee FROM %s WHERE order_uid=$1", paymentTable)

	row = r.db.QueryRowx(query, orderUid)

	err = row.StructScan(&p.Payment)
	if err != nil {
		logrus.Info("failed to scan data from %s table: %s", paymentTable, err.Error())
	}

	query = fmt.Sprintf("SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status FROM %s WHERE order_uid=$1", itemsTable)

	newRow, _ := r.db.Queryx(query, orderUid)

	var item l0wb.Item
	count := 0
	for newRow.Next() {
		err := newRow.StructScan(&item)
		if err != nil {
			logrus.Info("failed to scan data from %s table: %s", itemsTable, err.Error())
		}
		p.Items = append(p.Items, item)
		count++
	}

	return p
}
