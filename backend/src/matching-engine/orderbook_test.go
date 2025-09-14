package main

import (
	"fmt"
	"reflect"
	"testing"
)

func assert(t *testing.T, a, b any) {
	if !reflect.DeepEqual(a, b) {
		t.Errorf("Expected %v, got %v", a, b)
	}
}

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

func TestPlaceLimitOrder(t *testing.T) {
	ob := NewOrderBook()

	sellOrderA := NewOrder(false, 10)
	sellOrderB := NewOrder(false, 5)
	ob.PlaceLimitOrder(11000, sellOrderA)
	ob.PlaceLimitOrder(9000, sellOrderB)

	assert(t, len(ob.asks), 2)
}

func TestPlaceMarketOrder(t *testing.T) {
	ob := NewOrderBook()

	sellOrderA := NewOrder(false, 20)
	ob.PlaceLimitOrder(10000, sellOrderA)

	buyOrder := NewOrder(true, 10)
	matches := ob.PlaceMarketOrder(buyOrder)

	assert(t, len(matches), 1)
	assert(t, len(ob.asks), 1)
	assert(t, ob.AskQuantity(), 10)

	fmt.Printf("%+v\n", matches)

}
