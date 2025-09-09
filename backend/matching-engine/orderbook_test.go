package main

import (
	"fmt"
	"testing"
)

func TestLimit(t *testing.T) {
	l := NewLimit(100)
	buyOrder := NewOrder(true, 5)
	buyOrder2 := NewOrder(true, 5)
	buyOrder3 := NewOrder(true, 5)
	/*if l.Price != 100 {
		t.Errorf("Limit price should be 100, got %f", l.Price)
	}
	if buyOrder.Size != 5 {
		t.Errorf("Order size should be 5, got %f", buyOrder.Size)
	}*/
	l.AddOrder(buyOrder)
	l.AddOrder(buyOrder2)
	l.AddOrder(buyOrder3)

	fmt.Println(l)

}

func TestOrderBook(t *testing.T) {

}
