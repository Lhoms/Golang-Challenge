package sample1

import (
	"fmt"
	"time"
	"sync"
)

// PriceService is a service that we can use to get prices for the items
// Calls to this service are expensive (they take time)
type PriceService interface {
	GetPriceFor(itemCode string) (float64, error)
}

// price itself with the creation date (in the cache)
type Price struct {
	price        float64
	creationDate time.Time
}

// price request structure to abstract the request to GetPriceFor in the channel
type PriceRequest struct {
	price float64
	err   error
}

// Check price expiration by duration
func (p *Price) checkExpiration(maxAge time.Duration) bool {
	return (time.Now().Sub(p.creationDate)) < maxAge
}

// TransparentCache is a cache that wraps the actual service
// The cache will remember prices we ask for, so that we don't have to wait on every call
// Cache should only return a price if it is not older than "maxAge", so that we don't get stale prices
type TransparentCache struct {
	actualPriceService PriceService
	maxAge             time.Duration
	prices             map[string]Price
}

func NewTransparentCache(actualPriceService PriceService, maxAge time.Duration) *TransparentCache {
	return &TransparentCache{
		actualPriceService: actualPriceService,
		maxAge:             maxAge,
		prices:             map[string]Price{},
	}
}

// Mutex used to sync prices map
var m = sync.Mutex{}

func (c *TransparentCache) sequentialPriceSet(itemCode string, p Price) {
	m.Lock()
	c.prices[itemCode] = p
	m.Unlock()
}

// GetPriceFor gets the price for the item, either from the cache or the actual service if it was not cached or too old
func (c *TransparentCache) GetPriceFor(itemCode string) (float64, error) {
	priceStruct, ok := c.prices[itemCode]

	if ok && priceStruct.checkExpiration(c.maxAge) {
		return priceStruct.price, nil
	}
	price, err := c.actualPriceService.GetPriceFor(itemCode)
	if err != nil {
		return 0, fmt.Errorf("getting price from service : %v", err.Error())
	}

	c.sequentialPriceSet(itemCode, Price{price, time.Now()})

	return price, nil
}

func processRequest(cache *TransparentCache, c chan PriceRequest, itemCode string) {
	price, err := cache.GetPriceFor(itemCode)
	c <- PriceRequest{price, err}
}

// GetPricesFor gets the prices for several items at once, some might be found in the cache, others might not
// If any of the operations returns an error, it should return an error as well
func (c *TransparentCache) GetPricesFor(itemCodes ...string) ([]float64, error) {
	results := []float64{}
	messages := make(chan PriceRequest, len(itemCodes))

	// parallelize sending each get to a go routine
	for _, itemCode := range itemCodes {
		go processRequest(c, messages, itemCode)
	}

	// getting each response from 'messages' channel, until finish or error
	for i := len(itemCodes); i > 0; i-- {
		request := <-messages
		if request.err != nil {
			return []float64{}, request.err
		}
		results = append(results, request.price)
	}

	return results, nil
}
