package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type predictions struct {
	Predictions [][]float64 `json:"predictions"`
}

type request struct {
	MeteoST1 string `json:"meteost1"`
	MeteoST2 string `json:"meteost2"`
	Kitchen  string `json:"kitchen"`
}

type response struct {
	Boiler string `json:"boiler"`
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/goldfish", goldfish)
	log.Printf("server started")
	log.Fatal(http.ListenAndServe(":3000", mux))
}

func goldfish(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		log.Println("Method is:\t" + r.Method)
		log.Println("Request URL is:\t" + r.RequestURI)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

	case "POST":
		log.Println("Method is:\t" + r.Method)
		log.Println("Request URL is:\t" + r.RequestURI)

		body, err := ioutil.ReadAll(r.Body)
		var r request
		err = json.Unmarshal(body, &r)
		if err != nil {
			log.Println("goldfish: error unmarshaling request")
			fmt.Println(err)
			return
		}
		var resp response
		resp.Boiler = predict(r)

		fmt.Println("New task:\n", r)
		fmt.Println("Predition is:", resp.Boiler)

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			log.Println(err)
		}

	default:
		fmt.Fprintf(w, "Sorry, only GET and POST methods are supported.")
	}
}

func predict(r request) string {
	body := getBody(r)

	URL := "http://46.4.240.40:8501/v1/models/goldfish:predict"
	req, err := http.NewRequest("POST", URL, bytes.NewBuffer(body))
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	var p predictions
	err = json.Unmarshal(bodyBytes, &p)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(p.Predictions)
	var max float64
	var maxIndex int
	for i, p := range p.Predictions[0] {
		if p > max {
			max = p
			maxIndex = i
		}
	}
	fmt.Println("Max is:", max)
	fmt.Println("Max index is:", maxIndex)
	var res string
	switch maxIndex {
	case 0:
		res = "40"
	case 1:
		res = "41"
	case 2:
		res = "42"
	case 3:
		res = "43"
	case 4:
		res = "44"
	case 5:
		res = "45"
	case 6:
		res = "46"
	case 7:
		res = "47"
	case 8:
		res = "48"
	case 9:
		res = "49"
	}
	return res
}

func getBody(r request) []byte {
	meteoST1, err := strconv.ParseFloat(r.MeteoST1, 64)
	checkError("getBody: meteoST1 error parse:\n", err)

	meteoST2, err := strconv.ParseFloat(r.MeteoST2, 64)
	checkError("getBody: meteoST2 error parse:\n", err)

	kitchen, err := strconv.ParseFloat(r.Kitchen, 64)
	checkError("getBody: kitchen error parse:\n", err)

	var sequence []float64
	sequence = append(sequence, meteoST1)
	sequence = append(sequence, meteoST2)
	sequence = append(sequence, kitchen)

	var body [][]float64
	body = append(body, sequence)
	sJSON, err := json.Marshal(body)
	if err != nil {
		fmt.Println(err)
	}
	reqBody := `{"instances" : ` + string(sJSON) + ` }`
	return []byte(reqBody)
}

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}
