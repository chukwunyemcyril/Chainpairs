package main

import (
	"fmt"
	"time"
)

type Limit struct {
	Price    float64
	Quantity float64
	Orders   []*Order
}

type Match struct {
	Ask         *Order
	Bid         *Order
	SizedFilled float64
	Price       float64
}

type Order struct {
	Size      float64
	Bid       bool
	Limit     *Limit
	Timestamp int64
}

func (o *Order) String() string {
	return fmt.Sprintf("Order: %f %t %f %d", o.Size, o.Bid, o.Limit.Price, o.Timestamp)
}

type OrderBook struct {
	Asks []*Limit
	Bids []*Limit

	AskLimits map[float64]*Limit
	BidLimits map[float64]*Limit
}

func NewLimit(price float64) *Limit {
	return &Limit{
		Price:  price,
		Orders: []*Order{},
	}
}

func NewOrder(bid bool, size float64) *Order {
	return &Order{
		Size:      size,
		Bid:       bid,
		Timestamp: time.Now().UnixNano(),
	}

}

func (l *Limit) AddOrder(o *Order) {
	o.Limit = l
	l.Orders = append(l.Orders, o)
	l.Quantity += o.Size

}

func (l *Limit) DeleteOrder(o *Order) {
	for i := 0; i < len(l.Orders); i++ {
		if l.Orders[i] == o {
			l.Orders[i] = l.Orders[len(l.Orders)-1]
			l.Orders = l.Orders[:len(l.Orders)-1]
			l.Quantity -= o.Size
			break
		}
	}
}

func NewOrderBook() *OrderBook {
	return &OrderBook{
		Asks:      []*Limit{},
		Bids:      []*Limit{},
		AskLimits: make(map[float64]*Limit),
		BidLimits: make(map[float64]*Limit),
	}
}

func (ob *OrderBook) PlaceOrder(price float64, o *Order) []Match {
	if o.Size > 0.0 {
		ob.Add(price, o)
	}
	return []Match{}
}

func (ob *OrderBook) Add(price float64, o *Order) {
	var limit *Limit

	if o.Bid {
		limit = ob.BidLimits[price]
	} else {
		limit = ob.AskLimits[price]
	}

	if limit == nil {
		limit = NewLimit(price)
		if o.Bid {
			ob.BidLimits[price] = limit
			ob.Bids = append(ob.Bids, limit)
		} else {
			ob.AskLimits[price] = limit
			ob.Asks = append(ob.Asks, limit)
		}
	}

}
