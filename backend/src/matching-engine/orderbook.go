package main

import (
	"fmt"
	"sort"
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
type Orders []*Order

func (o Orders) Len() int           { return len(o) }
func (o Orders) Swap(i, j int)      { o[i], o[j] = o[j], o[i] }
func (o Orders) Less(i, j int) bool { return o[i].Timestamp < o[j].Timestamp }

func (o *Order) String() string {
	if o.Limit == nil {
		return fmt.Sprintf("Order: %f %t (no limit) %d", o.Size, o.Bid, o.Timestamp)
	}
	return fmt.Sprintf("Order: %f %t %f %d", o.Size, o.Bid, o.Limit.Price, o.Timestamp)
}

type OrderBook struct {
	asks []*Limit
	bids []*Limit

	AskLimits map[float64]*Limit
	BidLimits map[float64]*Limit
}

type Limits []*Limit

type ByBestAsk struct {
	Limits
}

func (a ByBestAsk) Len() int           { return len(a.Limits) }
func (a ByBestAsk) Swap(i, j int)      { a.Limits[i], a.Limits[j] = a.Limits[j], a.Limits[i] }
func (a ByBestAsk) Less(i, j int) bool { return a.Limits[i].Price < a.Limits[j].Price }

type ByBestBid struct{ Limits }

func (b ByBestBid) Len() int           { return len(b.Limits) }
func (b ByBestBid) Swap(i, j int)      { b.Limits[i], b.Limits[j] = b.Limits[j], b.Limits[i] }
func (b ByBestBid) Less(i, j int) bool { return b.Limits[i].Price > b.Limits[j].Price }

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

func (l *Limit) String() string {
	return fmt.Sprintf("Limit: %f %f", l.Price, l.Quantity)
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

func (l *Limit) Fill(o *Order) []Match {
	matches := []Match{}
	for _, order := range l.Orders {
		match := l.fillOrder(order, o)
		matches = append(matches, match)

		l.Quantity -= match.SizedFilled

		if o.isFilled() {
			break
		}
	}
	return matches

}

func (o *Order) isFilled() bool {
	return o.Size == 0
}

func (l *Limit) fillOrder(a, b *Order) Match {
	var (
		bid         *Order
		ask         *Order
		SizedFilled float64
	)

	if a.Bid {
		bid = a
		ask = b
	} else {
		bid = b
		ask = a
	}

	if a.Size >= b.Size {
		SizedFilled = b.Size
		a.Size -= SizedFilled
		b.Size = 0
	} else {
		SizedFilled = a.Size
		b.Size -= SizedFilled
		a.Size = 0
	}
	return Match{
		Ask:         ask,
		Bid:         bid,
		SizedFilled: SizedFilled,
		Price:       ask.Limit.Price,
	}
}

func NewOrderBook() *OrderBook {
	return &OrderBook{
		asks:      []*Limit{},
		bids:      []*Limit{},
		AskLimits: make(map[float64]*Limit),
		BidLimits: make(map[float64]*Limit),
	}
}

func (ob *OrderBook) PlaceLimitOrder(price float64, o *Order) {
	var limit *Limit

	if o.Bid {
		limit = ob.BidLimits[price]
	} else {
		limit = ob.AskLimits[price]
	}

	if limit == nil {
		limit = NewLimit(price)
		limit.AddOrder(o)

		if o.Bid {
			ob.BidLimits[price] = limit
			ob.bids = append(ob.bids, limit)
		} else {
			ob.AskLimits[price] = limit
			ob.asks = append(ob.asks, limit)
		}
	}

}

func (ob *OrderBook) PlaceMarketOrder(o *Order) []Match {

	matches := []Match{}
	if o.Bid {
		if o.Size > ob.AskQuantity() {
			panic(fmt.Errorf("not enough volume {%.2f} for market order{%.2f}", ob.AskQuantity(), o.Size))
		}
		for _, limit := range ob.asks {
			limitmatches := limit.Fill(o)
			matches = append(matches, limitmatches...)
		}
	} else {

	}

	return matches
}

//func (ob *OrderBook) PlaceOrder(price float64, o *Order) []Match {
// Try to match the ask and bid orders in the order book
//Mathcing Logic
//if o.Size > 0.0 {
//	ob.Add(price, o)
//}
//return []Match{}
//}

func (ob *OrderBook) BidQuantity() float64 {
	quantity := 0.0

	for i := 0; i < len(ob.bids); i++ {
		quantity += ob.bids[i].Quantity
	}
	return quantity
}

func (ob *OrderBook) AskQuantity() float64 {
	quantity := 0.0

	for i := 0; i < len(ob.asks); i++ {
		quantity += ob.asks[i].Quantity
	}
	return quantity
}
func (ob *OrderBook) Add(price float64, o *Order) {

}

func (ob *OrderBook) BestAsk() *Limit {
	sort.Sort(ByBestAsk{ob.asks})
	return ob.asks[0]

}

func (ob *OrderBook) BestBid() *Limit {
	sort.Sort(ByBestBid{ob.bids})
	return ob.bids[0]

}
