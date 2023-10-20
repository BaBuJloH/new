package orders_repo

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const TableName = "orders"

type Repo struct {
	TableName string
	ctx       context.Context
	conn      pgxpool.Pool
}

func New(ctx context.Context, conn *pgxpool.Pool) *Repo {
	return &Repo{
		tableName: TableName,
		ctx:       ctx,
		conn:      *conn,
	}
}

func (r *Repo) Select(where string) order.List {
	query := fmt.Sprintf("select order_uid, track_number, entry, delivery, payment, items, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard from %s", TableName)

	rows, err := r.conn.Query(r.ctx, query)
	if err != nil {
		return order.List{}
	}
	defer rows.Close()

	res := order.List{}

	for rows.Next() {
		data := order.Data{}
		if err := rows.Scan(&data.OrderUID, &data.TrackNumber, &data.Entry, &data.Delivery, &data.Payment, &data.Items,
			&data.Locale, &data.InternalSignature, &data.CustomerID, &data.DeliveryService, &data.Shardkey, &data.SmID,
			&data.DateCreated, &data.OofShard); err != nil {
			log.Print(err)
			return order.List{}
		}
		res = append(res, data)
	}
	return res
}

func (r *Repo) Get(query string) user.Data {
	return user.NewUser(0, "", "")
}

func (r *Repo) Insert(ordr order.Data) error {
	deliveryBuff, err := json.Marshal(ordr.Delivery)
	if err != nil {
		return err
	}

	paymentBuff, err := json.Marshal(ordr.Payment)
	if err != nil {
		return err
	}
	paymentStr := string(paymentBuff)

	itemsBuff, err := json.Marshal(ordr.Items)
	if err != nil {
		return err
	}

	ins := fmt.Sprintf("insert into %s (order_uid, track_number, entry, delivery, payment, items, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard) values ('%s', '%s', '%s', '%v', '%s', '%v', '%s', '%s', '%s', '%s', '%s', %d, '%v', '%s')",
		//                                                                                                                                                                                                                              ^^    ^^    ^^
		TableName, ordr.OrderUID, ordr.TrackNumber, ordr.Entry, string(deliveryBuff), string(paymentStr), string(itemsBuff), ordr.Locale, ordr.InternalSignature, ordr.CustomerID, ordr.DeliveryService, ordr.Shardkey, ordr.SmID, ordr.DateCreated.Format(time.RFC3339), ordr.OofShard)

	// TODO debug
	log.Print(ins)

	_, err = r.conn.Query(r.ctx, ins)
	return err
}

func (r *Repo) GET(uid string) order.Data {
	rows, err := r.conn.Query(r.ctx, fmt.Sprintf("select order_uid, track_number, delivery from %s where order_uid = '%s' limit 1", TableName, uid))
	if err != nil {
		return order.Data{}
	}
	defer rows.Close()

	res := order.Data{}
	for rows.Next() {
		if err := rows.Scan(&res.OrderUID, &res.TrackNumber, &res.Delivery); err != nil { // добавлен track_number
			return order.Data{}
		}
	}
	return res
}

func (r *Repo) DEL(uid string) order.Data {
	err := r.conn.QueryRow(r.ctx, fmt.Sprintf("delete from %s where order_uid = '%s'", TableName, uid))
	if err != nil {
		return order.Data{}
	}
	return order.Data{}
}
