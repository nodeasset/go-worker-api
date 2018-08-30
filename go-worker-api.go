package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"

	"github.com/streadway/amqp"

	"github.com/gorilla/mux"
)

var ch *amqp.Channel

type Configuration struct {
	//Rabbitmq config
	User      string
	Password  string
	Host      string
	Port      string
	Queue     string
	Consumer  string
	Autoack   bool
	Exclusive bool
	Nolocal   bool
	Nowait    bool
	Args      string
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func Health(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("I am healthy"))
}

func ExecCommand(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)

	cmdPrep := string(body)
	cmdOutput := exec.Command("bash", "-c", cmdPrep)
	stdout, err := cmdOutput.CombinedOutput()
	if err != nil {
		fmt.Println(err)
		//w.Write(string(err})
	}
	//fmt.Println(stdout)
	fmt.Printf("%s\n", stdout)
	w.Write([]byte(stdout))
}

// create a new item
func CreateEvent(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)

	err = ch.Publish(
		"exchange.events", // exchange
		"events",          // routing key
		false,             // mandatory
		false,             // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})

	log.Printf(" [x] Sent %s", body)
	failOnError(err, "Failed to publish a message")

}

// main function to boot up everything
func main() {

	//conf := "config.json"

	/*	if os.Args[1] == "development" {
			//"/home/earl/go/src/go-worker-api/config.json"
			conf = os.Args[2]
		} else {
			conf = "config.json"
		}

		configuration := Configuration{}
		err := gonfig.GetConf(conf, &configuration)
		if err != nil {
			panic(err)
		}*/

	/*
		Uncomment this section to use rabbitmq
		user := configuration.User
		password := configuration.Password
		host := configuration.Host
		port := configuration.Port

		connstring := "amqp://" + user + ":" + password + "@" + host + ":" + port + "/"

		conn, err := amqp.Dial(connstring)
		failOnError(err, "Failed to connect to RabbitMQ")
		defer conn.Close()

		ch, err = conn.Channel()
		failOnError(err, "Failed to open a channel")
		defer ch.Close()
	*/

	router := mux.NewRouter()

	router.HandleFunc("/events", CreateEvent).Methods("POST")
	router.HandleFunc("/exec", ExecCommand).Methods("POST")
	router.HandleFunc("/health", Health).Methods("GET")
	log.Fatal(http.ListenAndServe(":8000", router))
}
