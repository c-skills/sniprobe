// parallel SNI prober
//
// (C) 2022 by c->skills research
//

package main

import (
	"os"
	"fmt"
	"net"
	"time"
	"sync"
	"strings"
	"net/http"
	"io/ioutil"
	"crypto/tls"
)

func sniconnect(node, port, sni string) (int, string) {

	tlsConf := &tls.Config{
		ServerName: sni,
		InsecureSkipVerify: true,
	}

	tConf := &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 5*time.Second,
		}).Dial,
		TLSHandshakeTimeout: 10*time.Second,
		ResponseHeaderTimeout: 15*time.Second,
		TLSClientConfig: tlsConf,
	}

	client := &http.Client{
		//Timeout: 5*time.Second,
		Transport: tConf,
	}

	r, err := client.Get("https://" + node + ":" + port)
	if err != nil {
		return -1, err.Error()
	}
	defer r.Body.Close()

	return 0, r.Header.Get("Date")
}

func snis(file string) []string {

	data, err := ioutil.ReadFile(file)
	if (err != nil) {
		fmt.Printf("%s\n", err.Error())
		os.Exit(1);
	}

	return strings.Split(string(data), "\n")
}

func main() {

	if (len(os.Args) != 4) {
		fmt.Printf("Usage: %s <hostname> <port> <SNI-file>\n", os.Args[0]);
		return
	}

	ch := make(chan string)	// unbuffered

	node := os.Args[1]
	port := os.Args[2]
	snilist := snis(os.Args[3])

	maxgo := 10
	if (len(snilist) < 10) {
		maxgo = len(snilist)
	}

	var wg sync.WaitGroup

	for i := 0; i < maxgo; i++ {
		wg.Add(1)
		go func() {
			var sni string
			for {
				sni = <- ch
				if (len(sni) == 0) {
					break
				}
				r, err := sniconnect(node, port, sni)
				if (r == 0) {
					fmt.Printf("Success: https://%s:%s with SNI %s (%s)\n", node, port, sni, err)
				} else {
					fmt.Printf("Failed: https://%s:%s with SNI %s (%s)\n", node, port, sni, err)
				}
			}
			wg.Done()
		}()
	}

	for _, sni := range snilist {
		if (len(sni) <= 1 || sni[0] == '#') {
			continue;
		}
		ch <- sni
	}
	close(ch)
	wg.Wait()

	return
}

