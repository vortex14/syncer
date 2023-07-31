package client

import (
	"context"
	"log"
	"sync"
	"syncer/errors"
	"time"

	"syncer/interfaces"
)

// Client Main client structure
type Client struct {
	singleTon sync.Once // singleTon guaranty once execution

	Endpoint            string             //  url of external service
	Service             interfaces.Service //  instance of service bridge
	CheckStatusInterval time.Duration      //  interval for checking external service

	ExceptionCallback func(err error)              //  catch exception in callback
	SuccessCallback   func(batch interfaces.Batch) //  catch success batch in callback

	BufferBatch  int64 //  count cache of buffer batch
	BufferInput  int64 //  count cache of buffer input
	BufferResort int64 //  count cache of buffer resort

	input           chan interfaces.Item  // main channel of input items
	resortBatchChan chan interfaces.Batch // resort channel of batch
	batchChan       chan interfaces.Batch // channel of ready batch

	batchCount       uint64        // admissible count for sync
	batchProcessTime time.Duration // duration of batch per time

	success    int64 // success counter
	exceptions int64 // exceptions counter
}

// GetChanStats return len of input channel, len of resortBatchChan, len of batchChan
func (c *Client) GetChanStats() []int {
	return []int{
		len(c.input),
		len(c.resortBatchChan),
		len(c.batchChan),
	}
}

// GetStats return 2 counters: success and exceptions
func (c *Client) GetStats() (int64, int64) {
	return c.success, c.exceptions
}

// prepareBatch - create a new batch of items and put to batch channel
func (c *Client) prepareBatch() {
	var batch interfaces.Batch

	for {
		select {
		case item := <-c.input:
			if uint64(len(batch)) == c.batchCount {
				c.batchChan <- batch
				log.Printf("sended of ready batch: %d", len(batch))
				batch = batch[:0]
			}
			batch = append(batch, item)

		}
	}
}

// resortBatch - resort lost batch and put to main batch channel
func (c *Client) resortBatch() {
	for {
		select {
		case batch := <-c.resortBatchChan:
			c.batchChan <- batch
			log.Printf("batch sended for resorting to main channel")

		}
	}
}

// syncBatchExternalService send batch to external service and processing base sync exceptions
func (c *Client) syncBatchExternalService() {
	log.Print("init sync of batch")

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		for {
			select {
			case batch := <-c.batchChan:
				log.Printf("received batch of %d", len(batch))
				err := c.Service.Process(context.Background(), batch)

				switch err {
				case errors.ErrBlocked:
					log.Printf("catch blocked exception...")
					c.resortBatchChan <- batch
					log.Printf("batch skip for resorting: %d, count: %d", len(batch), len(c.resortBatchChan))
					wg.Done()
					if c.ExceptionCallback != nil {
						c.ExceptionCallback(err)
					}
					c.exceptions += 1
					return
				case nil:
					c.success += 1
					log.Printf("batch sended!")
					if c.SuccessCallback != nil {
						c.SuccessCallback(batch)
					}

				}

			}
		}
	}(wg)

	wg.Wait()
	log.Print("close sync goroutine")

	for range time.Tick(c.CheckStatusInterval) {
		log.Print("check status external service")
		c.batchCount, c.batchProcessTime = c.Service.GetLimits()
		if c.batchCount > 0 && c.batchProcessTime > 0 {

			break
		} else {
			log.Print("service is unavailable")
		}
	}
	log.Printf("ready for sync %d", c.batchCount)

	go c.syncBatchExternalService()

}

// AddNewItem add a new item for processing
func (c *Client) AddNewItem(item interfaces.Item) {
	c.input <- item
}

// Run syncer
func (c *Client) Run() *Client {

	if c.Service == nil {
		panic("service not found")
	}

	c.singleTon.Do(func() {
		c.initExternalLimit()

		c.input = make(chan interfaces.Item, c.BufferInput)
		c.batchChan = make(chan interfaces.Batch, c.BufferBatch)
		c.resortBatchChan = make(chan interfaces.Batch, c.BufferResort)

		go c.prepareBatch()
		go c.syncBatchExternalService()
		go c.resortBatch()

	})

	return c

}

// init external service and base settings
func (c *Client) initExternalLimit() {
	c.batchCount, c.batchProcessTime = c.Service.GetLimits()
}
