// Package utils contains generic and auxiliary methods, such as type converters, result checks and etc.
package utils

import (
	"log"
	"os"
	"strconv"
	"time"
)

// GetArgs will check and capture the argument passed by user in the command line or
// present an error if no argument was passed.
func GetArgs() (cmd, port, timeout string) {
	args := os.Args
	if len(os.Args[1:]) != 6 {
		UsageError()
	}

	if args[1] != "-script" && args[3] != "-port" && args[5] != "-timeout" {
		UsageError()
	}
	return args[2], args[4], args[6]
}

// CheckCmdOutput will check if the output of the custom script has the
// required number of arguments to proceed
func CheckCmdOutput(s []string) {
	for i := range s {
		if len(s) != 4 {
			log.Fatal(`ERROR: Custom script ouput does not follow the required format and cannot be exported. Exiting!!

For the output of custom script to be exported, output should be in the following format!

Ex: ./your-custom-script.sh 
    hostname, instance, metric_name, metric_value
    hostname, instance, metric_name, metric_value
    hostname, instance, metric_name, metric_value
			
			`)
		}
		i++
	}
}

// StringToInteger will receive a string and convert it an integer
func StringToInteger(s string) (r int) {
	r, err := strconv.Atoi(s)
	if err != nil {
		log.Fatal("ERROR: Invalid data type! Please use an integer instead of a string. Exiting!\n\n")
	}
	return r
}

// StringToSeconds will receive a string and convert it to seconds (time.Duration)
func StringToSeconds(s string) (r time.Duration) {
	_sec := StringToInteger(s)
	sec := time.Duration(_sec * 1000 * 1000 * 1000)
	return sec
}

// CheckEmptyMap verifies map content and and returns an error message if empty
func CheckEmptyMap(m map[int]map[string]string) {
	mapLength := len(m)
	if mapLength == 0 {
		log.Fatal("ERROR: Custom script provided an empty output and cannot be exported. Exiting!\n\n")
	}
}

// UsageError will display usage intructions when users types incorrect arguments.
func UsageError() {
	log.Fatal(`ERROR!

You did not specify a valid command or failed to pass the proper options. Exiting!

Ex.: prom_exporter -script <script_path> -port <port> -timeout <seconds>

`)
}

// StringToFloat64 converts a string to a float64
func StringToFloat64(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		os.Exit(1)
	}
	return f

}
