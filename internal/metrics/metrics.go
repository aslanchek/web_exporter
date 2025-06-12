package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	// Node-related metrics
	ClientsOnlineNum = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Subsystem: "web",
			Name:      "clients_online_num",
			Help:      "WEB number of clients are currently online",
		},
	)

	ClientsTotalNum = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Subsystem: "web",
			Name:      "clients_total_num",
			Help:      "WEB total number of clients",
		},
	)

	ClientsOnlineArr = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Subsystem: "web",
			Name:      "client_online",
			Help:      "WEB client are currently online",
		},
		[]string{"email"},
	)

	// Client-related metrics
	ClientUpTraffic = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Subsystem: "web",
			Name:      "client_up_traffic",
			Help:      "WEB up client traffic",
		},
		[]string{"email", "inboundID"},
	)

	ClientDownTraffic = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Subsystem: "web",
			Name:      "client_down_traffic",
			Help:      "WEB down client traffic",
		},
		[]string{"email", "inboundID"},
	)

	ClientTotalTraffic = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Subsystem: "web",
			Name:      "client_total_traffic",
			Help:      "WEB total client traffic",
		},
		[]string{"email", "inboundID"},
	)

	ClientTotalTrafficCalculated = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Subsystem: "web",
			Name:      "client_total_traffic_calculated",
			Help:      "WEB total client traffic",
		},
		[]string{"email", "inboundID"},
	)
)
