package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)
type Result struct{
	Status string
}

func main(){

	http.HandleFunc("/", home)
	http.ListenAndServe(":9093", nil)
}

func home(w http.ResponseWriter, r *http.Request){
	result := Result{Status: "invalid coupon code"}

	coupon := r.PostFormValue("coupon")

	if coupon == "A123" {
		result.Status = "approved"
	}

	jsonData, err := json.Marshal(result)
	if err != nil {
		log.Fatal("Error processing json")
	}

	fmt.Fprintf(w, string(jsonData))
}

