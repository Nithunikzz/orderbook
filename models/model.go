package models

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB

type CurrencyPair struct {
	Base  string `json:"base"`
	Quote string `json:"quote"`
}

type Order struct {
	ID             int       `json:"id"`
	UserID         int       `json:"user_id"`
	Pair           string    `json:"pair"`
	Side           string    `json:"side"`
	Price          float64   `json:"price"`
	Quantity       float64   `json:"quantity"`
	FilledQuantity float64   `json:"filled_quantity"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	Fee            float64   `json:"fee"`
}

func (o *Order) RemainingQuantity() float64 {
	return o.Quantity - o.FilledQuantity
}

func ConnectDB() {
	var err error
	DB, err = sql.Open("postgres", "postgres://postgres:password@localhost:5432/orderbook?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
}

func AddPair(pair CurrencyPair) error {
	_, err := DB.Exec("INSERT INTO currency_pairs (base, quote) VALUES ($1, $2)", pair.Base, pair.Quote)
	if err != nil {
		log.Println(" AddPair error:", err)
	}
	return err
}

func CreateOrder(o *Order) error {
	o.Fee = 0.05 * o.Price * o.Quantity

	return DB.QueryRow(`INSERT INTO orders (user_id, pair, side, price, quantity, status,fee)
                        VALUES ($1, $2, $3, $4, $5, 'open',$6) RETURNING id, created_at`,
		o.UserID, o.Pair, o.Side, o.Price, o.Quantity, o.Fee).Scan(&o.ID, &o.CreatedAt)
}

func UpdateOrder(o *Order) error {
	_, err := DB.Exec(`UPDATE orders SET filled_quantity=$1, status=$2 WHERE id=$3`,
		o.FilledQuantity, o.Status, o.ID)
	return err
}

func GetOpenOrders(pair, side string) []Order {
	rows, _ := DB.Query(`SELECT id, user_id, pair, side, price, quantity, filled_quantity, status, created_at
                         FROM orders WHERE pair=$1 AND side=$2 AND status IN ('open', 'partial')`, pair, side)
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var o Order
		rows.Scan(&o.ID, &o.UserID, &o.Pair, &o.Side, &o.Price, &o.Quantity, &o.FilledQuantity, &o.Status, &o.CreatedAt)
		orders = append(orders, o)
	}
	return orders
}

func GetUserOrders(userID int) []Order {
	rows, _ := DB.Query(`SELECT id, user_id, pair, side, price, quantity, filled_quantity, status, created_at
                         FROM orders WHERE user_id=$1`, userID)
	defer rows.Close()

	var orders []Order
	for rows.Next() {
		var o Order
		rows.Scan(&o.ID, &o.UserID, &o.Pair, &o.Side, &o.Price, &o.Quantity, &o.FilledQuantity, &o.Status, &o.CreatedAt)
		orders = append(orders, o)
	}
	return orders
}
