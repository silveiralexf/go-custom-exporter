// Package go-custom-exporter provides an easy way to export metrics from custom scripts or
// commands to Prometheus without having to worry about doing much code.
//
// The goal is that any routine which can return data in a expected output can be easily
// exported as a Gauge metric with a single command.
package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/fsilveir/go-custom-exporter/pkg/utils"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	script, _port, timeout := utils.GetArgs()
	port := fmt.Sprintf(":%v", utils.StringToInteger(_port))

	path := "/metrics"

	// Initialize the first level of the map of strings to store command output
	m := map[int]map[string]string{}

	// Execute command once out of a loop in order to get the required fields
	cmd := exec.Command(script)
	ch := make(chan struct{})

	go updateMap(cmd, ch, m)

	ch <- struct{}{}
	cmd.Start()
	<-ch

	if err := cmd.Wait(); err != nil {
		fmt.Println(err)
	}

	utils.CheckEmptyMap(m)

	gauge := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "prom",
			Subsystem: "custom",
			Name:      "exporter",
			Help:      "Prometheus Gauge Metrics from Custom script/command exporter",
		},
		[]string{
			// Name of the metric
			"system",
			"subsystem",
			"metric",
		},
	)
	// Start to serve the data on the target port
	http.Handle(path, promhttp.Handler())
	prometheus.MustRegister(gauge)

	for range m {

		fmt.Println(fmt.Sprintf("INFO: Successfully exposing metrics on port %s", (_port)))

		go func() {
			for {
				i := 0
				go executeCmd(script, m)
				for range m {
					go updateMap(cmd, ch, m)
					gauge.With(prometheus.Labels{
						"system":    m[i]["system"],
						"subsystem": m[i]["subsystem"],
						"metric":    m[i]["metric"]}).
						Set(utils.StringToFloat64(m[i]["value"]))
					i++
				}
				time.Sleep(utils.StringToSeconds(timeout))
			}
		}()
		log.Fatal(http.ListenAndServe(port, nil))
	}

}

// executeCmd will execute a custom scrip or command as a goroutine and return the output to a map
func executeCmd(script string, m map[int]map[string]string) {
	cmd := exec.Command(script)
	ch := make(chan struct{})

	go updateMap(cmd, ch, m)

	ch <- struct{}{}
	cmd.Start()
	<-ch

	if err := cmd.Wait(); err != nil {
		fmt.Println(err)
	}
	utils.CheckEmptyMap(m)
}

// updateMap will get the values from the custom script or command executed by executeCmd method
// and update a map that will be used to set the values to be exported my main funcion
func updateMap(cmd *exec.Cmd, ch chan struct{}, m map[int]map[string]string) (cmdOutput map[int]map[string]string) {

	defer func() { ch <- struct{}{} }()
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	<-ch

	scanner := bufio.NewScanner(stdout)

	i := 0
	for scanner.Scan() {
		s := strings.Split(scanner.Text(), ",")
		utils.CheckCmdOutput(s)

		// Initialize the second level of the map of strings to store the command output
		m[i] = map[string]string{}
		m[i]["system"] = strings.TrimSpace(s[0])
		m[i]["subsystem"] = strings.TrimSpace(s[1])
		m[i]["metric"] = strings.TrimSpace(s[2])
		m[i]["value"] = strings.TrimSpace(s[3])

		i++
	}
	return m
}
