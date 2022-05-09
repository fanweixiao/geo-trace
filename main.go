package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/http/httptrace"
	"time"
)

func main() {
	// timeGet("https://ws.vi-server.org")
	timeGet("https://prsc.yomo.dev")
	// timeGet("")
	// req, _ := http.NewRequest("GET", "http://ws.vi-server.org", nil)
	// trace := &httptrace.ClientTrace{
	// 	GotConn: func(connInfo httptrace.GotConnInfo) {
	// 		fmt.Printf("Got Conn: %+v\n", connInfo)
	// 	},
	// 	DNSDone: func(dnsInfo httptrace.DNSDoneInfo) {
	// 		fmt.Printf("DNS Info: %+v\n", dnsInfo)
	// 	},
	// }
	// req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	// _, err := http.DefaultTransport.RoundTrip(req)
	// if err != nil {
	// 	log.Fatal(err)
	// }
}

func timeGet(url string) {
	req, _ := http.NewRequest("GET", url, nil)

	var start, connect, dns, tlsHandshake time.Time

	trace := &httptrace.ClientTrace{
		DNSStart: func(dsi httptrace.DNSStartInfo) { dns = time.Now() },
		DNSDone: func(ddi httptrace.DNSDoneInfo) {
			fmt.Printf("DNS Done: %v\n\t-> IP: %v\n\t->%+v\n", time.Since(dns), ddi.Addrs, ddi)
		},

		TLSHandshakeStart: func() { tlsHandshake = time.Now() },
		TLSHandshakeDone: func(cs tls.ConnectionState, err error) {
			fmt.Printf("TLS Handshake: %v\n", time.Since(tlsHandshake))
			fmt.Printf("\t-> %+v\n", cs)
		},

		ConnectStart: func(network, addr string) { connect = time.Now() },
		ConnectDone: func(network, addr string, err error) {
			fmt.Printf("Connect time: %v\n", time.Since(connect))
		},

		GotFirstResponseByte: func() {
			fmt.Printf("Got FB: Time from start to first byte: %v\n", time.Since(start))
		},

		GotConn: func(connInfo httptrace.GotConnInfo) {
			fmt.Printf("Got Conn: %+v\n", connInfo)
		},

		Got100Continue: func() {
			fmt.Printf("Got 100 Continue\n")
		},
	}

	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	start = time.Now()
	if _, err := http.DefaultTransport.RoundTrip(req); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Total time: %v\n", time.Since(start))
}
