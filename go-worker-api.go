package main

import (
	"log"
	"net/http"

	"github.com/streadway/amqp"

	"github.com/tkanos/gonfig"

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

// create a new item
func CreateEvent(w http.ResponseWriter, r *http.Request) {
	type Body string
	//var body Body
	body := "test me" // Body(r.Body)
	err := ch.Publish(
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
	configuration := Configuration{}
	err := gonfig.GetConf("/home/earl/go/src/go-worker-api/config.json", &configuration)
	if err != nil {
		panic(err)
	}

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

	router := mux.NewRouter()
	//router.HandleFunc("/people", GetPeople).Methods("GET")
	//router.HandleFunc("/people/{id}", GetPerson).Methods("GET")
	//router.HandleFunc("/people/{id}", CreatePerson).Methods("POST")
	//router.HandleFunc("/people/{id}", DeletePerson).Methods("DELETE")

	router.HandleFunc("/events", CreateEvent).Methods("POST")
	log.Fatal(http.ListenAndServe(":8000", router))
}
