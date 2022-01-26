package main

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"time"
)

func main() {
	nc, _ := nats.Connect(nats.DefaultURL)

	nc.Subscribe("help", func(m *nats.Msg) {
		nc.Publish(m.Reply, []byte("I can help!"))
	})

	msg, err := nc.Request("help", []byte("help me"), 10*time.Millisecond)
	if err != nil {
		fmt.Println("err: ", err.Error())
	}

	fmt.Println("msg: ", string(msg.Data))

}
