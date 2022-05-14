package main

import (
	"flag"
	"fmt"
	"math/rand"
	"log"
	"os"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"time"
)
var t trace
type trace struct{
	fileName string
	file *os.File
	timeStart time.Time
}

func (t *trace) Print(tx_time float64, fileName string) {
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		t.file, _ = os.OpenFile(fileName, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
		fmt.Fprintf(t.file, " TX TIME \n")
	} else {
		t.file, _ = os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	}

	fmt.Fprintf(t.file,"%f\n",tx_time)
	t.file.Close()
}


var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func main() {
	mqtt.DEBUG = log.New(os.Stdout, "", 0)
	mqtt.ERROR = log.New(os.Stdout, "", 0)
	broker := flag.String("broker","tcps://127.0.0.1:1884", "Broker type")
	topic := flag.String("topic", "test", "The topic name to/from which to publish/subscribe")
	fileName := flag.String("file","","Files name")
	num := flag.Int("num",10,"number of pkts to transmit")
	flag.Parse()

	//Client options//
	opts := mqtt.NewClientOptions()
	opts.AddBroker(*broker)
	opts.SetClientID("ClientID")
	opts.SetUsername("")
	opts.SetPassword("")

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil{
		panic(token.Error())
	}
	/////////////////
	rand.Seed(time.Now().UnixNano())

	/////////////////
	opts.SetDefaultPublishHandler(f)

	start := time.Now()
	if token := client.Subscribe(*topic, 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	for i := 0; i < *num; i++ {
		text := fmt.Sprintf("%d", i)
		fmt.Printf("MSG: %d\n", i)
		token := client.Publish(*topic, 0, false, text)
		/*rand.Seed(time.Now().UnixNano())
		n := 50 + rand.Float64() * (100 - 50)
		fmt.Printf("Sleeping %f seconds...\n", n)*/
		n:=24
		time.Sleep(time.Duration(n)*time.Millisecond)
		token.Wait()
	}
	time.Sleep(5*time.Second)
	client.Disconnect(250)
	end := time.Now()
	tx_time := float64(end.Sub(start)/1000)
	time_tx := fmt.Sprintf("%s_sat_txt.tr",*fileName)
	t.Print(tx_time,time_tx)
	fmt.Printf("Time tx of 1000 msgs: %f us\n",tx_time)
}
