package main

import (
		"fmt"
		"log"
		"net/url"
//	"os"
	"time"

		mqtt "github.com/eclipse/paho.mqtt.golang"
)


func connect(clientId string, uri *url.URL) mqtt.Client{
	fmt.Printf("Connecting-...\n")
	//opts := createClientOptions("pub", uri)
	opts := createClientOptions("pub", uri)
	client := mqtt.NewClient(opts)

	token := client.Connect()

	for !token.WaitTimeout(3*time.Second){

	}
	if err := token.Error(); err != nil{
		log.Fatal(err)
	}

	return client
}


func createClientOptions(clientId string, uri *url.URL) *mqtt.ClientOptions{
	opts := mqtt.NewClientOptions()
	//opts.AddBroker(fmt.Sprintf("tcps://%s", uri.Host))
	//opts.SetUsername(uri.User.Username())
	opts.AddBroker(fmt.Sprintf("quic://%s", uri.Host))
	opts.SetUsername("")
	//password, _ := uri.User.Password()
	opts.SetPassword("")
	opts.SetClientID(clientId)

	return opts
}

func listen(uri *url.URL, topic string){
	fmt.Printf("Hola 1\n")
	client := connect("sub",uri)
	client.Subscribe(topic,0, func(client mqtt.Client, msg mqtt.Message){
		fmt.Printf("* [%s] %s\n", msg.Topic(),string(msg.Payload()))
	})

}
func main(){
	//mqtt.DEBUG = log.New(os.Stdout, "", 0)
	//mqtt.ERROR = log.New(os.Stdout, "", 0)
	//mqtt.WARN = log.New(os.Stdout, "", 0)
	//mqtt.CRITICAL = log.New(os.Stdout, "", 0)
	//	uri, err := url.Parse("quic://127.0.0.1:1883	/test")
	uri, err := url.Parse("quic://127.0.0.1:1883/test")

	if err != nil {
		log.Fatal(err)
	}

	topic := uri.Path[1:len(uri.Path)]

	if topic == "" {
		topic = "test"
	}


	client := connect("pub", uri)
	fmt.Printf("Connected !!\n")
	timer := time.NewTicker(2 *  time.Second)

	for t := range timer.C {
		msg := t.String()
		fmt.Printf("Publishing, %s\n", msg)
		client.Publish(topic, 0, false, msg)
	}
}
