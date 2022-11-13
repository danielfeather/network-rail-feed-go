package main

import (
	"github.com/go-stomp/stomp/v3"
	"github.com/joho/godotenv"
	"log"
	"net"
	"os"
	"time"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	godotenv.Load()

	log.SetOutput(os.Stdout)
	log.Println("Connecting to Network Rail", os.Getenv("SERVER_ADDRESS"))

	//netConn, err := tls.Dial("tcp", os.Getenv("SERVER_ADDRESS"), &tls.Config{})
	netConn, err := net.DialTimeout("tcp", os.Getenv("SERVER_ADDRESS"), time.Second*10)
	if err != nil {
		failOnError(err, "Failed to connect to Network Rail")
	}
	defer netConn.Close()
	conn, err := stomp.Connect(netConn, stomp.ConnOpt.Login(os.Getenv("NR_USERNAME"), os.Getenv("NR_PASSWORD")))
	if err != nil {
		failOnError(err, "Failed to connect to Network Rail")
	}
	defer conn.Disconnect()

	sub, err := conn.Subscribe("/topic/TD_ALL_SIG_AREA", stomp.AckAuto)
	if err != nil {
		failOnError(err, "Failed to subscribe")
	}

	log.Println(conn.Version().String())

	forever := make(chan bool)

	go func() {
		for msg := range sub.C {
			log.Println(string(msg.Body))
		}
	}()

	log.Println("Waiting for messages...")
	<-forever

}
