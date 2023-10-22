package orders_repo

import (
	"context"
	"encoding/json"
	//"encoding/json"
	"fmt"
	"log"
	"new/model/order"

	//"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const TableName = "orders"

type Repo struct {
	tableName string
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
	query := fmt.Sprintf("select order_uid, entry, internal_signature, payment, items, locale, customer_id, track_number, delivery_service, shardkey, sm_id from %s", TableName)

	rows, err := r.conn.Query(r.ctx, query)
	if err != nil {
		return order.List{}
	}
	defer rows.Close()

	res := order.List{}

	for rows.Next() {
		data := order.Order{}
		if err := rows.Scan(&data.OrderUID, &data.Entry, &data.InternalSignature,  &data.Payment, &data.Items,
			&data.Locale, &data.CustomerID, &data.TrackNumber, &data.DeliveryService, &data.Shardkey, &data.SmID); err != nil {
			log.Print(err)
			return order.List{}
		}
		res = append(res, data)
	}
	return res
}


func (r *Repo) Insert(ordr order.Order) error {
	// deliveryBuff, err := json.Marshal(ordr.Delivery)
	// if err != nil {
	// 	return err
	// }

	paymentBuff, err := json.Marshal(ordr.Payment)
	if err != nil {
		return err
	}
	paymentStr := string(paymentBuff)

	itemsBuff, err := json.Marshal(ordr.Items)
	if err != nil {
		return err
	}

	ins := fmt.Sprintf("insert into %s (order_uid, track_number, entry, delivery, payment, items, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard) values ('%s', '%s', '%s', '%s', '%s','%s, '%s', '%s', '%s', '%s', '%d')", //5,6,7 -  '%s', '%v', '%s',          13,14 - , '%v', '%s'
		//                                                                                                                                                                                                                              ^^    ^^    ^^
		TableName, ordr.OrderUID, ordr.Entry, ordr.InternalSignature, string(paymentStr),  string(itemsBuff),/*string(deliveryBuff),*/ ordr.Locale, ordr.CustomerID, ordr.TrackNumber, ordr.DeliveryService, ordr.Shardkey, ordr.SmID/*, ordr.DateCreated.Format(time.RFC3339), ordr.OofShard*/)

	// TODO debug
	log.Print(ins)

 	_, err = r.conn.Query(r.ctx, ins)
 	return err
}

func (r *Repo) GET(uid string) order.Order {
	rows, err := r.conn.Query(r.ctx, fmt.Sprintf("select order_uid, entry, internal_signature, payment, items, locale, customer_id, track_number, delivery_service, shardkey, sm_id from %s where order_uid = '%s' limit 1", TableName, uid))
	if err != nil {
		return order.Order{}
	}
	defer rows.Close()

	res := order.Order{}
	for rows.Next() {
		if err := rows.Scan(&res.OrderUID, &res.Entry, &res.InternalSignature,  &res.Payment, &res.Items,
			&res.Locale, &res.CustomerID, &res.TrackNumber, &res.DeliveryService, &res.Shardkey, &res.SmID); err != nil { // добавлен track_number
			return order.Order{}
		}
	}
	return res
}

func (r *Repo) DEL(uid string) order.Order {
	err := r.conn.QueryRow(r.ctx, fmt.Sprintf("delete from %s where order_uid = '%s'", TableName, uid))
	if err != nil {
		return order.Order{}
	}
	return order.Order{}
}
