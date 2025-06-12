package collector

import (
	"fmt"
	"log"
	"net/http"
	"web_exporter/internal/api"
	"web_exporter/internal/metrics"
	"web_exporter/pkg/logger"

	"github.com/prometheus/client_golang/prometheus"
)

type WebCollector struct {
	callUrl  string
	username string
	password string
	cookies  []*http.Cookie
}

func NewThe3XUICollector(callUrl string, username string, password string) *WebCollector {
	cookies, err := api.Auth(callUrl, username, password)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to auth: %w", err))
	}

	return &WebCollector{
		callUrl:  callUrl,
		username: username,
		password: password,
		cookies:  cookies,
	}
}

func (c *WebCollector) reauth() error {
	logger.Logf(logger.CollectorLogPrefix, "trying to reauthentificate")
	cookies, err := api.Auth(c.callUrl, c.username, c.password)

	if err != nil {
		return fmt.Errorf("failed to reauth in WEB panel: %w", err)
	}

	c.cookies = cookies

	return nil
}

/*
var (
	clientsOnline = prometheus.NewDesc(
		prometheus.BuildFQName("", "web", "clients_online"),
		"How many web clients are online currently",
		nil,
		nil,
	)
)
*/

func (e *WebCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(e, ch)
}

func (e *WebCollector) Collect(ch chan<- prometheus.Metric) {
	logger.Logf(logger.CollectorLogPrefix, "collection process started")

	onlines, err := api.GetOnlines(e.callUrl, e.cookies)
	if _, ok := err.(api.SessionExpiredError); ok {
		err = e.reauth()
	}

	if err != nil {
		logger.Logf(logger.CollectorLogPrefix, "failed to get onlines %s", err)
		return
	}

	metrics.ClientsOnlineNum.Set(float64(len(onlines)))
	ch <- metrics.ClientsOnlineNum

	clients, err := api.GetClients(e.callUrl, e.cookies)
	if _, ok := err.(api.SessionExpiredError); ok {
		err = e.reauth()
	}

	if err != nil {
		logger.Logf(logger.CollectorLogPrefix, "failed to get clients %s", err)
		return
	}

	for i := range clients {
		labels := []string{clients[i].Email, fmt.Sprintf("%d", clients[i].InboundID)}

		metrics.ClientUpTraffic.WithLabelValues(labels...).Set(float64(clients[i].UpTraffic))
		ch <- metrics.ClientUpTraffic.WithLabelValues(labels...)

		metrics.ClientDownTraffic.WithLabelValues(labels...).Set(float64(clients[i].DownTraffic))
		ch <- metrics.ClientDownTraffic.WithLabelValues(labels...)

		metrics.ClientTotalTraffic.WithLabelValues(labels...).Set(float64(clients[i].TotalTraffic))
		ch <- metrics.ClientTotalTraffic.WithLabelValues(labels...)

		metrics.ClientTotalTrafficCalculated.WithLabelValues(labels...).Set(float64(clients[i].UpTraffic) + float64(clients[i].DownTraffic))
		ch <- metrics.ClientTotalTrafficCalculated.WithLabelValues(labels...)
	}

	metrics.ClientsTotalNum.Set(float64(len(clients)))
	ch <- metrics.ClientsTotalNum

	for i := range clients {
		isOnline := 0

		for j := range onlines {
			if clients[i].Email == onlines[j] {
				isOnline = 1
				break
			}
			isOnline = 0
		}

		metrics.ClientsOnlineArr.WithLabelValues(clients[i].Email).Set(float64(isOnline))
		ch <- metrics.ClientsOnlineArr.WithLabelValues(clients[i].Email)
	}

	logger.Logf(logger.CollectorLogPrefix, "collection process finished")
}
