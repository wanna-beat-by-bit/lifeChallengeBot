package consumer

import (
	"log"
	ev "tgBot/events"
	"time"
)

type Consumer struct {
	fetcher   ev.Fethcer
	processor ev.Processor
	batchSize int
}

func NewConsumer(fetcher ev.Fethcer, processor ev.Processor, batchSize int) *Consumer {
	return &Consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
	}
}

func (c *Consumer) Consume() {
	for {
		events, err := c.fetcher.Fetch(c.batchSize)

		if err != nil {
			log.Printf("[ERR]: %s", err.Error())
			continue
		}

		if len(events) == 0 {
			time.Sleep(time.Second * 1)
			continue
		}
		for _, event := range events {
			if err := c.processor.Process(event); err != nil {
				log.Printf("[ERR]: %s", err)
			}
		}

	}
}
