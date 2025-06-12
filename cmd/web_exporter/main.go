package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"log"
	"net/http"
	"os"
	"web_exporter/internal/collector"
	"web_exporter/pkg/logger"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// cmd opts
var (
	API3XUIHOST  = flag.String("api-host", "127.0.0.1", "web api hostname")
	API3XUIPORT  = flag.String("api-port", "8000", "web api port")
	EXPORTERBIND = flag.String("bind", "", "exporter bind address")
	EXPORTERPORT = flag.String("port", "9200", "exporter bind port")

	USERNAME = os.Getenv("THE3XUI_USERNAME")
	PASSWORD = os.Getenv("THE3XUI_PASSWORD")

	UseTLS = flag.Bool("tls", false, "wheter or not use TLS (see -key, -cert and -ca)")

	Key  = flag.String("key", "", "tls key")
	Cert = flag.String("cert", "", "tls certificate")
	CA   = flag.String("ca", "", "tls root ca cert")
)

func main() {
	logger.Logf(logger.MainLogPrefix, "Hello!\n")

	flag.Parse()

	var URL = "http://" + *API3XUIHOST + ":" + *API3XUIPORT
	var ADDR = *EXPORTERBIND + ":" + *EXPORTERPORT

	exporter := collector.NewThe3XUICollector(URL, USERNAME, PASSWORD)

	r := prometheus.NewRegistry()
	r.MustRegister(exporter)

	handler := promhttp.HandlerFor(r, promhttp.HandlerOpts{})

	http.Handle("/metrics", handler)

	if *UseTLS {
		if *Key == "" || *Cert == "" || *CA == "" {
			log.Fatal("must specify private key, certificate and CA cert")
		}

		caCert, err := os.ReadFile(*CA)
		if err != nil {
			log.Fatal(err)
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		// Create the TLS Config with the CA pool and enable Client certificate validation
		tlsConfig := &tls.Config{
			ClientCAs:  caCertPool,
			ClientAuth: tls.RequireAndVerifyClientCert,
		}

		// Create a Server instance to listen on port 8443 with the TLS config
		server := &http.Server{
			Addr:      ADDR,
			TLSConfig: tlsConfig,
			Handler:   http.DefaultServeMux,
		}

		log.Fatal(server.ListenAndServeTLS(*Cert, *Key))
	} else {
		log.Fatal(http.ListenAndServe(ADDR, nil))
	}
}
