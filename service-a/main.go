package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

const (
	SERVICE_NAME    = "service-a"
	ZIPKIN_URL      = "http://zipkin:9411/api/v2/spans"
	CEP_SERVICE_URL = "http://service-b:8081"
)

func main() {
	log.Println("Starting server...")

	initZipkin()

	http.HandleFunc("/", handleRequest)

	port := "8080"
	log.Printf("Server running on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func initZipkin() {
	exporter, err := zipkin.New(ZIPKIN_URL)
	if err != nil {
		log.Fatalf("Could not create zipkin exporter: %s", err.Error())
	}
	tp := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewWithAttributes(semconv.SchemaURL, semconv.ServiceNameKey.String(SERVICE_NAME))),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	carrier := propagation.HeaderCarrier(r.Header)
	ctx := r.Context()
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)

	ctx, span := otel.Tracer(SERVICE_NAME).Start(ctx, "handleCEPRequest")
	defer span.End()

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

	response, err := fetchTemperatureByCEP(ctx, cep)
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

func fetchTemperatureByCEP(ctx context.Context, cep string) (*Response, error) {
	_, span := otel.Tracer(SERVICE_NAME).Start(ctx, "fetchTemperatureByCEP")
	defer span.End()

	url := CEP_SERVICE_URL + "/cep/" + cep
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	resp, err := http.DefaultClient.Do(req)
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
