package main

import (
	"caozhipan/nsq-prometheus-exporter/controllers"
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
)

var (
	nsqLookupdAddress = flag.String("nsq.lookupd.address", "127.0.0.1:4161", "nsqllookupd address list with comma")
	k8s       = flag.String("nsq.k8s.mode", "", "k8s mode")
	port = flag.String("port", ":9527", "port usage")
	k8sMode = false
)


func main() {
	flag.Parse()
	if len(*k8s) > 0 {
		k8sMode = true
	}

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		for {
			controllers.SyncNodeList(*nsqLookupdAddress, k8sMode)
			<-ticker.C
		}
	}()

	prometheus.MustRegister(controllers.Collector)

	http.Handle("/metrics", promhttp.Handler())
	//pingHandler := ping()
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request){
		fmt.Fprint(w, "pong")
	})

	fmt.Printf("Listen %s", *port)
	log.Fatal(http.ListenAndServe(*port, nil))

}
