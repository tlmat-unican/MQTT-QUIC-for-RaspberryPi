package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/lucas-clemente/quic-go"
	"math/big"
	"net"
	"os"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "1883"
	CONN_TYPE = "tcp"
)



// Setup a bare-bones TLS config for the server
func generateTLSConfig() *tls.Config {
	key, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		panic(err)
	}
	template := x509.Certificate{SerialNumber: big.NewInt(1)}
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &key.PublicKey, key)
	if err != nil {
		panic(err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		panic(err)
	}
	return &tls.Config{
		Certificates: []tls.Certificate{tlsCert},
		NextProtos:   []string{"quic-echo-example"},
	}
}


func main() {
	//tcp_bool := true




	// Listen for incoming connections.
	tcp, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}

	listener, err := quic.ListenAddr(CONN_HOST+":"+CONN_PORT, generateTLSConfig(), nil)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}

	conn_chann := make(chan net.Conn)
	sess_chann := make(chan quic.Session)

	go func(){
		for {
			sess, err := listener.Accept(context.Background())
			if err != nil {
				fmt.Println("Error accepting: ", err.Error())
				os.Exit(1)
			}
			sess_chann <- sess
		}
	}()

	go func(){
		for {
			conn, err := tcp.Accept()
			if err != nil {
				fmt.Println("Error accepting: ", err.Error())
				os.Exit(1)
			}
			conn_chann <- conn
		}
	}()

	defer tcp.Close()
	defer listener.Close()
	for {
	select {
	case sess := <- sess_chann:    //QUIC session

		fmt.Println("Established QUIC connection")

		//stream, err := sess.AcceptStream()
		stream,err := sess.AcceptStream(context.Background())
		if err != nil {
			panic(err)
		}

		go handleRequestQuic(stream)

	case conn := <- conn_chann:     //TCP Connection
		fmt.Println("Established TCP connection")
		go handleRequest(conn)

		}
	}

	/*
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
	for {

		if tcp_bool == true{
			// Listen for an incoming connection.
			conn, err := l.Accept()
			if err != nil {
				fmt.Println("Error accepting: ", err.Error())
				os.Exit(1)
			}
			// Handle connections in a new goroutine.
			go handleRequest(conn)
		}else{
			// Listen for an incoming connection.
			conn, err := l.Accept()
			if err != nil {
				fmt.Println("Error accepting: ", err.Error())
				os.Exit(1)
			}
			// Handle connections in a new goroutine.
			go handleRequest(conn)
		}

	}
	 */
}

// Handles incoming requests.
func handleRequest(conn net.Conn) {
	// Make a buffer to hold incoming data.
	buf := make([]byte, 1024)

	for {
		fmt.Printf("Reading")
		// Read the incoming connection into the buffer.
		reqLen, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error reading:", err.Error())
		}
		if reqLen == 0 {
			fmt.Printf("Closing Connection\n")
			conn.Close()
			return
		}
		fmt.Printf("Received message, size: %d %d\n", len(buf), reqLen)
		for i := 0; i <= reqLen; i++{
			fmt.Printf("%x ", buf[i])
		}
		fmt.Printf("\n")

		//Start processing MQTT packets
		if (buf[0] == 0x10) {
			fmt.Printf("MQTT Init Conex\n")
			mesg := make([]byte, 4)
			mesg[0] = 0x20
			mesg[1] = 0x02
			mesg[2] = 0x00
			mesg[3] = 0x00
			conn.Write(mesg)
		} else if (buf[0] == 0x30) {
			fmt.Printf("MQTT Publish \n")
		} else {
			fmt.Printf("Unkown Type: %x", buf[0])
		}
	}


	// Send a response back to person contacting us.
	//conn.Write([]byte("Message received."))
	// Close the connection when you're done with it.
	//conn.Close()
}

// Handles incoming requests.
func handleRequestQuic(stream quic.Stream) {
	// Make a buffer to hold incoming data.
	buf := make([]byte, 1024)

	for {
		// Read the incoming connection into the buffer.
		reqLen, err := stream.Read(buf)
		if err != nil {
			fmt.Println("Error reading:", err.Error())
		}
		if reqLen == 0 {
			fmt.Printf("Closing Connection\n")
			stream.Close()
			return
		}
		fmt.Printf("Received message, size: %d %d\n", len(buf), reqLen)
		for i := 0; i <= reqLen; i++{
			fmt.Printf("%x ", buf[i])
		}
		fmt.Printf("\n")

		//Start processing MQTT packets
		if (buf[0] == 0x10) {
			fmt.Printf("MQTT Init Conex\n")
			mesg := make([]byte, 4)
			mesg[0] = 0x20
			mesg[1] = 0x02
			mesg[2] = 0x00
			mesg[3] = 0x00
			stream.Write(mesg)
		} else if (buf[0] == 0x30) {
			fmt.Printf("MQTT Publish \n")
		} else {
			fmt.Printf("Unkown Type: %x", buf[0])
		}
	}


	// Send a response back to person contacting us.
	//conn.Write([]byte("Message received."))
	// Close the connection when you're done with it.
	//conn.Close()
}