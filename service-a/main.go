package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
)

func main() {
	log.Println("Starting server...")

	http.HandleFunc("/", handleRequest)

	port := "8080"
	log.Printf("Server running on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	log.Printf("Request: %s %s", r.Method, r.URL.Path)

	if r.Method != http.MethodPost {
		log.Printf("Method not allowed: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		CEP string `json:"cep"`
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Printf("Invalid request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if request.CEP == "" {
		log.Printf("CEP is required")
		http.Error(w, "CEP is required", http.StatusBadRequest)
		return
	}

	cep := request.CEP
	if !isValidCEP(cep) {
		log.Printf("Invalid zipcode: %s", cep)
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	response, err := fetchTemperatureByCEP(cep)
	if err != nil {
		log.Print(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func isValidCEP(cep string) bool {
	if len(cep) != 8 {
		return false
	}
	re := regexp.MustCompile(`^\d+$`)
	return re.MatchString(cep)
}

type Response struct {
	City  string  `json:"city"`
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

func fetchTemperatureByCEP(cep string) (*Response, error) {
	url := fmt.Sprintf("http://service-b:8081/cep/%s", cep)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Error fetching temperature by CEP, status code: %d", resp.StatusCode)
		return nil, errors.New("can not fetch temperature by CEP: " + cep)
	}

	var data Response
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return &data, nil
}
