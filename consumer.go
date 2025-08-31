package gomq

import (
	"encoding/json"
	"log"

	"github.com/go-playground/validator/v10"
)

type AddConsumerParams struct {
	Channel       string
	Consumer      func(interface{})
	IsCategorized bool
}

func NewConsumer[T any](fn func(T) error) func(interface{}) {
	return func(data interface{}) {
		var form T
		b, err := json.Marshal(data)
		if err != nil {
			log.Printf("Consumer: marshal failed: %v", err)
			return
		}
		if err := json.Unmarshal(b, &form); err != nil {
			log.Printf("Consumer: unmarshal failed: %v", err)
			return
		}
		if err := validator.New().Struct(form); err != nil {
			log.Printf("Consumer: validation failed: %v", err)
			return
		}

		if err := fn(form); err != nil {
			log.Printf("Consumer: handler failed: %v", err)
		}
	}
}

func AddConsumer(params AddConsumerParams) {
	if params.IsCategorized {
		params.Channel = categorizeChannel(config.ChannelDir, params.Channel)
		return
	}

	config.Consumers[params.Channel] = params.Consumer
}
