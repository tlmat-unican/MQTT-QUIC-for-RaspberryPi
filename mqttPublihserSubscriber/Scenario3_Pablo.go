package main

import (
	"flag"
	"fmt"
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


	//mqtt.DEBUG = log.New(os.Stdout, "", 0)
	//mqtt.ERROR = log.New(os.Stdout, "", 0)
	//mqtt.WARN = log.New(os.Stdout, "", 0)
	//mqtt.CRITICAL = log.New(os.Stdout, "", 0)

	broker := flag.String("broker","quic://127.0.0.1:1883", "Broker type")
	//broker := flag.String("broker","quic://127.0.0.1:1884", "Broker type")
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



	start := time.Now()
	if token := client.Connect(); token.Wait() && token.Error() != nil{
		panic(token.Error())
	}
	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf("Msg %s and topic %s\n", string(msg.Payload()), *topic)
	})
	token := client.Publish(*topic, 0, false, "hola que tal.. oasopansdasduinasmdfdfsgssdfggsdfgsdfgsdfgsdfsdafasdfsadfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdadczasdfasdcasedfasdcasdfsdgsdfgsdfgsdfgsdfgsdfgsdfgsdfgsdfgsdfgsdfgsdfgsdfgsdfgsdfgsdfgsdfgkdnbasuidnasdmsaukidbnasmdbakjshbdnasbdhjasgmdjna bdsjhasjdbashgdjjhabsdkjajsbdhkasykdubhasgdmjy,ashjdbagvsdjgua,hjsdbvgasdg,kuasjbdhjagsydlukhasbjvygakstugdukhasbjhvdtanbhjhjksbsdkfjhbsdhbfsdf sdahfjnbsdhjfsdfghdfgnbdfjshgnkjdnhgmsdbfkjgybnshdfgbhskfughnsdmjfbhsdf,hjgbnsdkjfhmghsjmfbdgnjmfdhgkjsmdbfhgnskufhsj fmnghjksmyfnhbgbmbdfmsnjfgkjbsdhgfnuamsdbfgnmaj,dbhgmasjbkgjaydfhgb jadfhngjadgfnads,kgh,bhjddsajkbhjgb dfchjsdjfkcxbdjhfgeyjgnfcqjefgbhjygenberfhejbngceurgfgfbqhejcvujhqbdghkfvbsbdhlksd.jbalksdfghbuasdfhghajbslñgihsabdulkisfhudghgasuñkfgilsagisbdklugfuaklsbghlifs")
	token.Wait()
	end := time.Now()
	tx_time := float64(end.Sub(start)/1000) //To get microseconds.
	time_tx := fmt.Sprintf("%s_txt.tr",*fileName)
	t.Print(tx_time,time_tx)
	fmt.Printf("Time of, connect-publish: %f us\n", tx_time)

	time.Sleep(1 * time.Second)

	client.Disconnect(0)

	time.Sleep(1 * time.Second)


	start = time.Now()
	if token := client.Connect(); token.Wait() && token.Error() != nil{
		panic(token.Error())
	}

	token = client.Publish(*topic, 0, false, "hola que tal.. oasopansdasduinasmdfdfsgssdfggsdfgsdfgsdfgsdfsdafasdfsadfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdadczasdfasdcasedfasdcasdfsdgsdfgsdfgsdfgsdfgsdfgsdfgsdfgsdfgsdfgsdfgsdfgsdfgsdfgsdfgsdfgsdfgkdnbasuidnasdmsaukidbnasmdbakjshbdnasbdhjasgmdjna bdsjhasjdbashgdjjhabsdkjajsbdhkasykdubhasgdmjy,ashjdbagvsdjgua,hjsdbvgasdg,kuasjbdhjagsydlukhasbjvygakstugdukhasbjhvdtanbhjhjksbsdkfjhbsdhbfsdf sdahfjnbsdhjfsdfghdfgnbdfjshgnkjdnhgmsdbfkjgybnshdfgbhskfughnsdmjfbhsdf,hjgbnsdkjfhmghsjmfbdgnjmfdhgkjsmdbfhgnskufhsj fmnghjksmyfnhbgbmbdfmsnjfgkjbsdhgfnuamsdbfgnmaj,dbhgmasjbkgjaydfhgb jadfhngjadgfnads,kgh,bhjddsajkbhjgb dfchjsdjfkcxbdjhfgeyjgnfcqjefgbhjygenberfhejbngceurgfgfbqhejcvujhqbdghkfvbsbdhlksd.jbalksdfghbuasdfhghajbslñgihsabdulkisfhudghgasuñkfgilsagisbdklugfuaklsbghlifs")
	token.Wait()
	end = time.Now()

	//time.Sleep(10 * time.Second)

	client.Disconnect(0)


	tx_time = float64(end.Sub(start)/1000) //To get microseconds.
	time_tx = fmt.Sprintf("%s_txt.tr",*fileName)
	t.Print(tx_time,time_tx)
	fmt.Printf("Time of, connect-publish: %f us\n", tx_time)
}
