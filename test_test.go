package syncer

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	. "syncer/client"
	"syncer/errors"
	"syncer/interfaces"
	"syncer/service"
)

func TestCoolService(t *testing.T) {
	Convey("test with even number of items ", t, func() {

		client := (&Client{
			BufferBatch:         10,
			BufferInput:         50,
			BufferResort:        3,
			Endpoint:            "https://external-service.local",
			CheckStatusInterval: 5 * time.Second,
			Service:             &service.ExternalCoolService{ProcessLimit: 2, Duration: 10 * time.Second},
		}).Run()

		for i := 0; i <= 10; i++ {
			client.AddNewItem(interfaces.Item{})
		}

		for _ = range time.Tick(1 * time.Second) {

			success, _ := client.GetStats()

			if success == 5 {
				So(success, ShouldEqual, 5)
				break
			}

		}

	})

}

func TestCoolServiceOddItems(t *testing.T) {
	Convey("test with odd number of items", t, func(c C) {

		client := (&Client{
			BufferBatch:         10,
			BufferInput:         50,
			BufferResort:        3,
			Endpoint:            "https://external-service.local",
			CheckStatusInterval: 5 * time.Second,
			Service:             &service.ExternalCoolService{ProcessLimit: 3, Duration: 10 * time.Second},
		}).Run()

		for i := 0; i <= 10; i++ {
			client.AddNewItem(interfaces.Item{})
		}

		client.SuccessCallback = func(batch interfaces.Batch) {
			c.So(len(batch), ShouldEqual, 3)
		}

		for _ = range time.Tick(1 * time.Second) {

			success, _ := client.GetStats()
			if success == 3 {
				So(success, ShouldEqual, 3)
				break
			}

		}

	})

}

func TestExceptionHandler(t *testing.T) {
	Convey("test exception handler", t, func(c C) {
		status := false

		client := (&Client{
			BufferBatch:         10,
			BufferInput:         50,
			BufferResort:        3,
			Endpoint:            "https://external-service.local",
			CheckStatusInterval: 5 * time.Second,
			Service:             &service.ExternalBadService{ProcessLimit: 3, Duration: 10 * time.Second},
		}).Run()

		for i := 0; i <= 10; i++ {
			client.AddNewItem(interfaces.Item{})
		}

		client.ExceptionCallback = func(err error) {
			c.So(err, ShouldEqual, errors.ErrBlocked)
			status = true
		}

		for _ = range time.Tick(1 * time.Second) {

			if status {
				break
			}

		}

	})

}

func TestException(t *testing.T) {
	Convey("test unavailable service", t, func(c C) {

		client := (&Client{
			BufferBatch:         10,
			BufferInput:         50,
			BufferResort:        3,
			Endpoint:            "https://external-service.local",
			CheckStatusInterval: 5 * time.Second,
			Service:             &service.ExternalBadService{ProcessLimit: 5, Duration: 10 * time.Second},
		}).Run()

		for i := 0; i <= 20; i++ {
			client.AddNewItem(interfaces.Item{})
		}

		for _ = range time.Tick(1 * time.Second) {

			batchCount := client.GetChanStats()[2]

			if batchCount == 4 {
				c.So(batchCount, ShouldEqual, 4)
				break
			}

		}

	})

}
