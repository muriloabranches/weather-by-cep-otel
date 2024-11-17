package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
)

func main() {
	log.Println("Starting server...")

	http.HandleFunc("/cep/", handleCEPRequest)

	port := "8081"
	log.Printf("Server running on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func handleCEPRequest(w http.ResponseWriter, r *http.Request) {
	log.Printf("Request: %s %s", r.Method, r.URL.Path)

	if r.Method != http.MethodGet {
		log.Printf("Method not allowed: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cep := strings.TrimPrefix(r.URL.Path, "/cep/")
	if !isValidCEP(cep) {
		log.Printf("Invalid zipcode: %s", cep)
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	location, err := fetchLocation(cep)
	if err != nil {
		log.Printf("Can not find zipcode: %s", cep)
		http.Error(w, "can not find zipcode", http.StatusNotFound)
		return
	}

	temperature, err := fetchTemperature(location)
	if err != nil {
		log.Printf("Can not find temperature for location: %s", location)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var response struct {
		City  string  `json:"city"`
		TempC float64 `json:"temp_C"`
		TempF float64 `json:"temp_F"`
		TempK float64 `json:"temp_K"`
	}
	response.City = location
	response.TempC = temperature
	response.TempF = round(celsiusToFahrenheit(temperature))
	response.TempK = round(celsiusToKelvin(temperature))

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

func celsiusToFahrenheit(celsius float64) float64 {
	return (celsius * 1.8) + 32
}

func celsiusToKelvin(celsius float64) float64 {
	return celsius + 273
}

func round(value float64) float64 {
	return math.Round(value*10) / 10
}

func fetchLocation(cep string) (string, error) {
	url := fmt.Sprintf("https://viacep.com.br/ws/%s/json/", cep)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("failed to fetch location")
	}

	var data struct {
		Location string `json:"localidade"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}
	if data.Location == "" {
		return "", errors.New("location not found")
	}

	return data.Location, nil
}

func fetchTemperature(location string) (float64, error) {
	apiKey := os.Getenv("WEATHERAPI_KEY")
	if apiKey == "" {
		return 0, errors.New("missing WEATHERAPI_KEY")
	}

	encodedLocation := url.QueryEscape(location)
	url := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s", apiKey, encodedLocation)

	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, errors.New("failed to fetch temperature")
	}

	var data struct {
		Current struct {
			TempC float64 `json:"temp_c"`
		} `json:"current"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, err
	}

	return data.Current.TempC, nil
}
