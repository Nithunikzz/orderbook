package orderbook

import (
	"orderbook/models"
	"sort"
	"sync"
)

var mu sync.Mutex

func MatchOrders(pair string) {
	mu.Lock()
	defer mu.Unlock()

	buyOrders := models.GetOpenOrders(pair, "buy")
	sellOrders := models.GetOpenOrders(pair, "sell")

	// Sort: Buy -> Highest price first, earliest time
	sort.Slice(buyOrders, func(i, j int) bool {
		return buyOrders[i].Price > buyOrders[j].Price ||
			(buyOrders[i].Price == buyOrders[j].Price && buyOrders[i].CreatedAt.Before(buyOrders[j].CreatedAt))
	})

	// Sort: Sell -> Lowest price first, earliest time
	sort.Slice(sellOrders, func(i, j int) bool {
		return sellOrders[i].Price < sellOrders[j].Price ||
			(sellOrders[i].Price == sellOrders[j].Price && sellOrders[i].CreatedAt.Before(sellOrders[j].CreatedAt))
	})

	for i := range buyOrders {
		buy := &buyOrders[i]
		for j := range sellOrders {
			sell := &sellOrders[j]

			if buy.Price >= sell.Price && buy.RemainingQuantity() > 0 && sell.RemainingQuantity() > 0 {
				matchedQty := min(buy.RemainingQuantity(), sell.RemainingQuantity())
				matchedPrice := sell.Price
				tradeAmount := matchedQty * matchedPrice
				fee := 0.05 * tradeAmount

				buy.FilledQuantity += matchedQty
				buy.Fee += fee
				sell.FilledQuantity += matchedQty
				sell.Fee += fee
				if buy.RemainingQuantity() == 0 {
					buy.Status = "filled"
				} else {
					buy.Status = "partial"
				}

				if sell.RemainingQuantity() == 0 {
					sell.Status = "filled"
				} else {
					sell.Status = "partial"
				}

				models.UpdateOrder(buy)
				models.UpdateOrder(sell)
			}
		}
	}
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

// 5 percent fee
// need to be taken
// for every order
// and put that in the table and sum it
//docker exec -it orderbook-db psql -U postgres -d orderbook
