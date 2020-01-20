package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"gopkg.in/resty.v1"
)

type Supplier struct {
	id                int
	name              string
	office            int
	reachedOrderValue float64
}

const URL_API = "https://adlacarte.adesso.de/api/"

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
			fetchMetricsForSupplier(&supplier, client)
		}
		writeMetrics(suppliers[:], w)
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func fetchMetricsForSupplier(supplier *Supplier, client *resty.Client) {
	var e error
	response, e := client.R().Get(fmt.Sprintf("%s/offices/%d/suppliers/%d", URL_API, supplier.office, supplier.id))
	if e != nil {
		log.Println(e)
	}
	var responseAsString string = string(response.Body())
	supplier.reachedOrderValue, e = strconv.ParseFloat(responseAsString, len(responseAsString))
	if e != nil {
		log.Println(e)
	}
}

func writeMetrics(suppliers []Supplier, w http.ResponseWriter) {
	for _, supplier := range suppliers {
		fmt.Fprintf(w, "reachedOrderValue{supplier=\"%s\",office=\"%d\"} %f \n", supplier.name, supplier.office, supplier.reachedOrderValue)
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
		Get(URL_API + "login")
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

/*
var CREDENTIALS string

const URL_LOGIN = URL_API + "login"
const URL_SUPPLIERS = URL_API + "offices/1/suppliers/"
const URL_REACHED_ORDER_VALUE = "/reachedOrderValue"

var suppliersToName = [7]string{"Entenhaus", "ChiliPeppers", "PiDoe", "ChinaImbissBUI", "PizzeriaMammaMia", "Pinnochio", "PizzariabeiMarco"}

var suppliersToGauge [7]prometheus.Gauge

func main() {
	CREDENTIALS = readCredentialsFromFile("credentials")
	http.Handle("/metrics", promhttp.Handler())

	for i := 0; i < len(suppliersToGauge); i++ {
		suppliersToGauge[i] = prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: "adLaCarte",
				Name:      suppliersToName[i],
			})
		prometheus.MustRegister(suppliersToGauge[i])
	}
	client := resty.New()

	go func() {
		// Console output
		line := ""
		for i := 0; i < len(suppliersToName); i++ {
			line += "| " + suppliersToName[i] + " "
		}
		log.Println(line + "|")

		for true {
			updateValues(client)
			time.Sleep(15 * time.Second)
		}
	}()
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func updateValues(client *resty.Client) {
	line := ""
	for i := 0; i < len(suppliersToGauge); i++ {
		value := getValueOf(i+1, client)
		suppliersToGauge[i].Set(value)

		// Console output
		line += "| "
		valueAsString := strconv.FormatFloat(value, 'f', 2, 64)
		if valueAsString == "0" {
			valueAsString = "0.00"
		}
		length := len(suppliersToName[i]) - len(valueAsString)
		for j := 0; j < length; j++ {
			line += " "
		}
		line += valueAsString + " "
	}
	log.Println(line + "|")
}

func authorize(client *resty.Client) {
	_, e := client.R().
		SetHeader("authorization", "Basic "+CREDENTIALS).
		Get(URL_LOGIN)
	if e != nil {
		log.Fatal(e)
	}
}

func getValueOf(supplier int, client *resty.Client) float64 {
	url := URL_SUPPLIERS + strconv.Itoa(supplier) + URL_REACHED_ORDER_VALUE
	response, e := client.R().
		Get(url)
	if e != nil {
		log.Fatal(e)
	}
	if response.StatusCode() == 401 {
		authorize(client)
		response, e = client.R().
			Get(url)
		if e != nil {
			log.Fatal(e)
		}
	}
	valueAsString := string(response.Body())
	for len(valueAsString) < 3 {
		valueAsString = "0" + valueAsString
	}
	length := len(valueAsString)
	valueAsString = valueAsString[:length-2] + "." + valueAsString[length-2:]
	value, e := strconv.ParseFloat(valueAsString, 64)
	if e != nil {
		log.Fatal(e)
	}
	return value
}
*/
