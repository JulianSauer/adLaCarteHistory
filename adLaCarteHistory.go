package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"gopkg.in/resty.v1"
)

type Supplier struct {
	id                int
	name              string
	office            int
	reachedOrderValue float32
}

var apiURL = getEnv("ADLACARTE_API_URL", "https://adlacarte.adesso.de/api/")

func main() {
	var suppliers = []Supplier{
		Supplier{0, "Entenhaus", 1, 0.0},
		Supplier{1, "Chili Peppers", 1, 0.0},
		Supplier{2, "PIDÖ", 1, 0.0},
		Supplier{3, "China Imbiss BUI", 1, 0.0},
		Supplier{4, "Pizzeria Mamma Mia", 1, 0.0},
		Supplier{5, "Pinnochio", 1, 0.0},
		Supplier{6, "Pizzeria bei Marco", 1, 0.0},
	}
	client := buildClient()
	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		for _, supplier := range suppliers {
			if e := fetchMetricsForSupplier(&supplier, client); e != nil {
				log.Println(e)
			}
		}
		writeMetrics(suppliers[:], w)
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func fetchMetricsForSupplier(supplier *Supplier, client *resty.Client) error {
	response, e := client.R().Get(fmt.Sprintf("%s/offices/%d/suppliers/%d", apiURL, supplier.office, supplier.id))
	if e != nil {
		return e
	}
	responseAsString := string(response.Body())
	floatVal, e := strconv.ParseFloat(responseAsString, 32)
	if e != nil {
		return e
	}
	supplier.reachedOrderValue = float32(floatVal)
	return nil
}

func writeMetrics(suppliers []Supplier, w io.Writer) {
	for _, supplier := range suppliers {
		fmt.Fprintf(w, "reachedOrderValue{supplier=\"%s\",office=\"%d\"} %.2f\n", supplier.name, supplier.office, supplier.reachedOrderValue)
	}
}

func buildClient() *resty.Client {
	client := resty.New()
	credentials := readCredentialsFromFile("credentials")
	authorizeClientWithCredentials(client, credentials)
	return client
}

func authorizeClientWithCredentials(client *resty.Client, credentials string) {
	_, e := client.R().
		SetHeader("authorization", "Basic "+credentials).
		Get(apiURL + "login")
	if e != nil {
		log.Fatal(e)
	}
}

func readCredentialsFromFile(file string) string {
	var credentials = ""
	if credentialFile, e := ioutil.ReadFile(file); e != nil {
		log.Fatal(e)
	} else {
		credentials = strings.ReplaceAll(string(credentialFile), "\n", "")
		if credentials == "" {
			log.Fatal("Credentials missing")
		}
	}
	return credentials
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
