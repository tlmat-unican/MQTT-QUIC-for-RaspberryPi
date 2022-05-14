package main

import (
	"flag"
	"fmt"
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

func main(){
	mqtt.DEBUG = log.New(os.Stdout, "", 0)
	mqtt.ERROR = log.New(os.Stdout, "", 0)
	mqtt.WARN = log.New(os.Stdout, "", 0)
	mqtt.CRITICAL = log.New(os.Stdout, "", 0)

	broker := flag.String("broker","tcps://10.0.0.1:1884", "Broker type")
	topic := flag.String("topic", "test", "The topic name to/from which to publish/subscribe")
	fileName := flag.String("file","","Files name")
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


	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf("Msg %s and topic %s\n", string(msg.Payload()), *topic)
	})
	start := time.Now()
	for i := 0 ; i<1000; i++ {

		token := client.Publish(*topic, 0, false, i)
		token.Wait()
		fmt.Printf("Publish: %d\n", i)


		if token := client.Subscribe(*topic, 0, nil); token.Wait() && token.Error() != nil {
			fmt.Println(token.Error())
			os.Exit(1)
		}
		fmt.Printf("Receiving msg %d of the topic %s\n",i,*topic)
	}
	end := time.Now()
	tx_time := float64(end.Sub(start)/1000) //To get microseconds.




	time_tx := fmt.Sprintf("%s_1000_txt.tr",*fileName)

	t.Print(tx_time,time_tx)
	fmt.Printf("Time tx of 1000 msgs: %f us\n", tx_time)

	client.Disconnect(250)
}
