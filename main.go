package main

import (
	"encoding/json"
	"fmt"
	"github.com/metskem/tibberexporter/conf"
	"github.com/metskem/tibberexporter/model"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	priceTotal     = promauto.NewGauge(prometheus.GaugeOpts{Name: "tibber_price_total", Help: "Total price, combined energy and tax"})
	priceEnergy    = promauto.NewGauge(prometheus.GaugeOpts{Name: "tibber_price_energy", Help: "Net energy price, without tax"})
	priceTax       = promauto.NewGauge(prometheus.GaugeOpts{Name: "tibber_price_tax", Help: "Energy tax price"})
	tibberResponse model.TibberResponse
)

/* getDataFromTibber - Get data from the api endpoint from tibber, and store in TibberDataMap, repeat this every 24 hours. */
func getDataFromTibber() {
	go func() {
		firstCall := true
		ticker := time.NewTicker(8 * time.Hour) // every 8 hours should be enough to keep 2 days of data all time
		for ; true; <-ticker.C {
			// api call to tibber
			transport := http.Transport{IdleConnTimeout: time.Second}
			if request, err := http.NewRequest(http.MethodPost, conf.TibberGraphQLUrl, strings.NewReader(conf.Query)); err != nil {
				log.Printf("failed to create an http request: %s", err)
			} else {
				request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", conf.ApiToken))
				request.Header.Add("Content-Type", "application/json")
				client := http.Client{Timeout: time.Duration(conf.HttpTimeout) * time.Second, Transport: &transport}
				if resp, err := client.Do(request); err != nil {
					log.Printf("failed to POST request to %s: %s", conf.TibberGraphQLUrl, err)
				} else {
					if resp == nil {
						log.Printf("got a nil response from %s", conf.TibberGraphQLUrl)
					} else {
						var body []byte
						if body, err = io.ReadAll(resp.Body); err != nil {
							log.Printf("failed to read response body: %s", err)
						} else {
							_ = resp.Body.Close()
							if resp.StatusCode != http.StatusOK {
								log.Printf("%d response from %s: %s", resp.StatusCode, conf.TibberGraphQLUrl, string(body))
							} else {
								// now the response should be good
								if err = json.Unmarshal(body, &tibberResponse); err != nil {
									log.Printf("failed to unmarshal response: %s", err)
								} else {
									if firstCall {
										startExportData()
										firstCall = false
									}
								}
							}
						}
					}
				}
			}
		}
	}()
}

/* startExportData - Copy data from the TibberResponse for the current hour to the prometheus metrics, repeat this every hour. */
func startExportData() {
	go func() {
		ticker := time.NewTicker(time.Hour) // every 8 hours should be enough to keep 2 days of data all time
		for ; true; <-ticker.C {
			currentTime := time.Now()
			var tibberDataAll []model.TibberData
			tibberDataAll = append(tibberDataAll, tibberResponse.Data.Viewer.Homes[0].CurrentSubscription.PriceInfo.Today...)
			if tibberResponse.Data.Viewer.Homes[0].CurrentSubscription.PriceInfo.Tomorrow != nil {
				tibberDataAll = append(tibberDataAll, tibberResponse.Data.Viewer.Homes[0].CurrentSubscription.PriceInfo.Tomorrow...)
			}
			for _, tibberData := range tibberDataAll {
				if tibberData.StartsAt.Before(currentTime) && tibberData.StartsAt.Add(time.Hour).After(currentTime) {
					log.Printf("%v", tibberData)
					priceTotal.Set(tibberData.Total)
					priceEnergy.Set(tibberData.Energy)
					priceTax.Set(tibberData.Tax)
				}
			}
		}
	}()
}

func main() {
	log.SetOutput(os.Stdout)
	conf.EnvironmentComplete()
	prometheus.Unregister(prometheus.NewGoCollector())                                       // switch off default golang metrics
	prometheus.Unregister(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{})) // switch off default process metrics

	getDataFromTibber()

	// startup the listener that handles the scrape requests
	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(":2112", nil); err != nil {
		fmt.Println(err)
	}
}
