package main

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "gopkg.in/resty.v1"
    "io/ioutil"
    "log"
    "net/http"
    "strconv"
    "strings"
    "time"
)

var CREDENTIALS string

const URL_API = "https://adlacarte.adesso.de/api/"
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

func readCredentialsFromFile(file string) string {
	if credentialFile, e := ioutil.ReadFile(file); e != nil {
		log.Fatal(e)
		panic("Error reading credentials file")
    } else {
        credentials := strings.ReplaceAll(string(credentialFile), "\n", "")
        if credentials == "" {
			log.Fatal("Credentials missing")
		}
		return credentials
    }
}