package main

import (
	"encoding/json"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/joho/godotenv"
	uuid "github.com/satori/go.uuid"
	"github.com/streadway/amqp"
	"io/ioutil"
	"log"
	"net/url"
	"github.com/wesleywillians/go-rabbitmq/queue"
)

type Result struct {
	Status string
}

type Order struct{
	ID 			uuid.UUID
	Coupon 		string
	CcNumber 	string
}

func NewOrder() Order{
	return Order{ID: uuid.NewV4()}
}

const (
	InvalidCoupom="invalid"
	ValidCoupon="valid"
	ConnectionError ="connection error"
)

func init(){
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env")
	}
}

func main() {
	messageChannel := make(chan amqp.Delivery)

	rabbitMQ := queue.NewRabbitMQ()
	ch := rabbitMQ.Connect()
	defer ch.Close()

	rabbitMQ.Consume(messageChannel)

	for msg := range messageChannel{
		process(msg)
	}
}

func process(msg amqp.Delivery) {
	order:= NewOrder()
	json.Unmarshal(msg.Body,&order)

	/*result := Result{Status: "declined"}

	if ccNumber == "1" {
		result.Status = "approved"
	}else{
		result.Status = "declined"
	}
	*/
	/*if order.Coupon != "" {

		//Aqui seria um outro microserviço para validação do coupon code
		resultCoupon := makeHttpCall("http://localhost:9092", order.Coupon)
		//result.Status = resultCoupon.Status
	}
	*/
	resultCoupon := makeHttpCall("http://localhost:9092", order.Coupon)

	switch resultCoupon.Status{
	case InvalidCoupom:
		log.Println("Order: ",order.ID,": invalid coupon!")

	case ConnectionError:
		msg.Reject(false)
		log.Println("Order: ",order.ID,": could not process!")

	case ValidCoupon:
		log.Println("Order: ",order.ID,": Processed")
	}


	/*
		jsonData, err := json.Marshal(result)
		if err != nil {
			log.Fatal("Error processing json")
		}

		fmt.Fprintf(w, string(jsonData))
	*/
}


func makeHttpCall(urlMicroservice string, coupon string) Result {
	values := url.Values{}
	values.Add("coupon", coupon)

	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 5

	res, err := retryClient.PostForm(urlMicroservice, values)
	if err != nil {
		result := Result{Status: ConnectionError}
		return result
	}

	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal("Error processing result")
	}

	result := Result{}

	json.Unmarshal(data, &result)

	return result

}
