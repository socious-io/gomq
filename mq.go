package gomq

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"
)

type MessageQueue struct {
	client    *nats.Conn
	consumers map[string]func(interface{})
}

var NatsClient *nats.Conn
var Mq MessageQueue

func (mq *MessageQueue) subscribe(channel string, consumer func(interface{})) {
	client := mq.client

	mq.consumers[channel] = consumer

	_, err := client.Subscribe(channel, func(msg *nats.Msg) {
		var dest interface{}

		err := json.Unmarshal(msg.Data, &dest)
		if err != nil {
			fmt.Printf("received invalid JSON payload: %s\n", msg.Data)
		} else {
			fmt.Printf("received valid JSON payload: %+v\n", dest)
		}
		consumer(dest)
	})

	if err != nil {
		fmt.Printf("Channel '%s' failed to be subscribed, Error: %s", channel, err.Error())
	}
}

func (mq *MessageQueue) queueSubscribe(channel string, count int, consumer func(interface{})) {
	client := mq.client

	mq.consumers[channel] = consumer
	queue := "main"

	for i := 0; i < count; i++ {
		go func() {
			_, err := client.QueueSubscribe(channel, queue, func(msg *nats.Msg) {
				var dest interface{}

				err := json.Unmarshal(msg.Data, &dest)
				if err != nil {
					fmt.Printf("received invalid JSON payload: %s\n", msg.Data)
				} else {
					fmt.Printf("received valid JSON payload: %+v\n", dest)
				}
				consumer(dest)
			})

			if err != nil {
				fmt.Printf("Channel '%s' failed to be subscribed, Error: %s", channel, err.Error())
			}
		}()
	}
}

func (mq *MessageQueue) SendJson(channel string, message interface{}) {
	client := mq.client
	channel = categorizeChannel(config.ChannelDir, channel)

	payload, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Couldn't parse/marshal JSON data", message)
	}

	err = client.Publish(channel, payload)

	if err != nil {
		fmt.Println("Couldn't publish JSON data, Error", err.Error())
	}
}

func Connect() {
	NatsClient, err := nats.Connect(
		config.Url,
		nats.Token(config.Token),
		nats.DisconnectHandler(func(nc *nats.Conn) {
			fmt.Println("[NATS] Disconnected")
		}),
		nats.ReconnectHandler(func(nc *nats.Conn) {
			fmt.Println("[NATS] Reconnected to", nc.ConnectedUrl())
		}),
		nats.ClosedHandler(func(nc *nats.Conn) {
			fmt.Println("[NATS] Connection closed")
		}),
	)

	if err != nil {
		fmt.Printf("Nats failed to connect, Error: %s", err)
		os.Exit(1)
	} else {
		fmt.Printf("Nats is connected to: %s\n", NatsClient.ConnectedAddr())
	}

	Mq = MessageQueue{
		NatsClient,
		map[string]func(interface{}){},
	}
}

func Init() {

	Connect()
	registerConsumers()

	// Wait for termination signal (CTRL+C, Docker stop, etc.)
	terminationSignal := make(chan os.Signal, 1)
	signal.Notify(terminationSignal, syscall.SIGINT, syscall.SIGTERM)

	<-terminationSignal
	fmt.Println("[NATS] Shutting down gracefully...")
	Mq.client.Drain()
	time.Sleep(1 * time.Second) // give time to drain before exit

}

// Register Services
func registerConsumers() {
	for channel, worker := range config.Consumers {
		Mq.subscribe(channel, worker)
	}
}
