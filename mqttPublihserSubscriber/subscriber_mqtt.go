package main

import (
	"sync"
	"fmt"
	"log"
	"net/url"
	//"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)


func connect(clientId string, uri *url.URL) mqtt.Client{
	fmt.Printf("Connecting \n")
	opts := createClientOptions(clientId, uri)
	client := mqtt.NewClient(opts)
	token := client.Connect()

	for !token.WaitTimeout(3*time.Second){

	}
	if err := token.Error(); err != nil{
		log.Fatal(err)
	}
	fmt.Printf("Connected \n")
	return client
}


func createClientOptions(clientId string, uri *url.URL) *mqtt.ClientOptions{
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("quic://%s", uri.Host))
	//opts.SetUsername(uri.User.Username())
	//password, _ := uri.User.Password()
	//opts.SetPassword(password)
	opts.SetClientID(clientId)
	fmt.Printf("Creating opts... \n")
	return opts
}

func listen(uri *url.URL, topic string){
	client := connect("sub", uri)
	client.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message){
		fmt.Printf("* [%s] %s\n", msg.Topic(), string(msg.Payload()))
	})
}

func main(){
	//uri, err := url.Parse("tcp://127.0.0.1:1883/test")

	//if err != nil {
	//	log.Fatal(err)
	//}

	//topic := uri.Path[1:len(uri.Path)]

	//if topic == "" {
	//	topic = "test"
	//}

	//go listen(uri, topic)

	//fmt.Printf("Trying... Connected %s\n",uri.Host)
	//client := connect("sub", uri)
	//fmt.Printf("Connected")
	//timer := time.NewTicker(1 *  time.Second)



	topic := "test"

	opts := mqtt.NewClientOptions().AddBroker("quic://localhost:1883")

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	var wg sync.WaitGroup
	wg.Add(1)

	if token := client.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {

		fmt.Printf("Subscription Topic: %s Mess: %s \n", topic, msg.Payload())
		//if string(msg.Payload()) != "mymessage" {
		//	t.Fatalf("want mymessage, got %s", msg.Payload())
		//}
		//wg.Done()
	}); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	wg.Wait()
}
