package handy

import (
	"testing"
	"fmt"
	"time"
	"math/rand"
)

type AskBid int
type MarketLimit int

const (
	Ask AskBid = iota
	Bid
)

const (
	Market MarketLimit = iota
	Limit
)

type Order struct {
	AskBid      AskBid
	MarketLimit MarketLimit
	Amount      int
	Price       float64 //used by Market Tran
}

var Orders chan Order

func OrderBinarySearch(array []Order, first, last int, value Order) int {
	for first < last {
		mid := first + (last-first)/2
		if array[mid].Price < value.Price {
			first = mid + 1
		} else {
			last = mid
		}
	}
	return first
}

func InsertOrder(orders []Order, order Order) []Order {
	i := OrderBinarySearch(orders, 0, len(orders), order)
	orders = append(orders, Order{})
	copy(orders[i+1:], orders[i:])
	orders[i] = order
	return orders
}

func BrokerMainLoop() {
	Asks := make([]Order, 0)
	Bids := make([]Order, 0)
	for {
		select {
		case order := <-Orders:
			switch order.AskBid {
			case Ask:
				Asks = InsertOrder(Asks, order)
			case Bid:
				Bids = InsertOrder(Bids, order)
			}
		}
		//trade off the Asks and Bids

		for len(Asks) > 0 && len(Bids) > 0 && Asks[0].Price <= Bids[len(Bids)-1].Price {
			if Asks[0].Amount < Bids[len(Bids)-1].Amount {
				Bids[len(Bids)-1].Amount = Bids[len(Bids)-1].Amount - Asks[0].Amount
				fmt.Println("Deal!!!", Asks[0])
				Asks = Asks[1:]
				continue
			}
			if Asks[0].Amount > Bids[len(Bids)-1].Amount {
				Asks[0].Amount = Asks[0].Amount - Bids[len(Bids)-1].Amount
				fmt.Println("Deal!!!", Bids[len(Bids)-1])
				Bids = Bids[:len(Bids)-1]
				continue
			}
		}

		if len(Asks) > 5 && len(Bids) > 5 {
			fmt.Println("Loweast 5 Asks", Asks[0:5])
			fmt.Println("Highest 5 Bids", Bids[len(Bids)-5:])
			fmt.Println("-----------------------------------")
		}
	}
}

func DealerMainLoop() {
	for {
		select {
		case <-time.After(time.Millisecond * 20):
			//create dummy orders and send to order chan
			Orders <- Order{AskBid: AskBid(rand.Intn(2)), MarketLimit: Market, Amount: rand.Intn(1000), Price: float64(rand.Intn(1000))}
		}
	}
}

func TestPlayWithChan(t *testing.T) {
	//Init Orders
	Orders = make(chan Order, 1000000)
	go BrokerMainLoop()
	go DealerMainLoop()
	select {}
}
