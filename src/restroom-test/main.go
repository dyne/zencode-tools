package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type Args struct {
	scheme string
	host   string
	port   int
	url    string
	data   string
}

type ZenroomResult struct {
	Result string `json:"result"`
	Logs   string `json:"logs"`
}

// Warning: names of properties must start with upper caser
// otherwise JSON Unmarshal doesn't work without errors
type RestroomError struct {
	ZenroomError ZenroomResult `json:"zenroom_errors"`
	Result       string        `json:"result"`
	Exception    string        `json:"exception"`
}

type RestroomResult struct {
	ZenroomResult ZenroomResult `json:"zenroom_errors"`
	Result        string        `json:"result"`
	Exception     string        `json:"exception"`
}

func (args *Args) loadCli() {
	// flags declaration using flag package
	flag.StringVar(&args.url, "u", "", "Path to contract.")
	flag.StringVar(&args.scheme, "s", "http", "Restroom scheme (http/https)")
	flag.StringVar(&args.host, "h", "127.0.0.1", "Restroom host.")
	flag.IntVar(&args.port, "p", 5000, "Restroom port.")
	flag.StringVar(&args.data, "a", "", "Path to data file.")

	flag.Parse() // after declaring flags we need to call it
}

func (args Args) requestUrl() string {
	if (args.scheme == "http" && args.port == 80) ||
		(args.scheme == "https" && args.port == 443) {
		return fmt.Sprintf("%s://%s/api/%s", args.scheme, args.host, args.url)
	}
	return fmt.Sprintf("%s://%s:%d/api/%s", args.scheme, args.host, args.port, args.url)
}

// Example usage: ./restroom-test -p 12001 -u sandbox/did-document-create -a data.json
func main() {
	var args Args
	var err error
	var resp *http.Response
	args.loadCli()

	log.Printf("Testing endpong: %s\n", args.requestUrl())
	if args.data == "" {
		resp, err = http.Get(args.requestUrl())
	} else {
		data, err := os.Open(args.data)
		keys := make(map[string]interface{})
		if err != nil {
			log.Printf("Error while opening %s: %s\n", args.data, err.Error())
			os.Exit(-1)
		}
		defer data.Close()
		dataBytes, _ := io.ReadAll(data)
		datakeys := make(map[string]interface{})
		dataMap := make(map[string]interface{})
		err = json.Unmarshal(dataBytes, &dataMap)
		if err != nil {
			log.Printf("Could not read json in file: %s\n", err.Error())
			os.Exit(-1)
		}
		datakeys["data"] = dataMap
		datakeys["keys"] = keys
		datakeysBytes, err := json.Marshal(datakeys)

		log.Printf("Data in file: %s\n", args.data)
		resp, err = http.Post(args.requestUrl(), "application/json",
			bytes.NewReader(datakeysBytes))
	}
	if err != nil {
		log.Printf("Error in the request: %s\n", err)
		os.Exit(-1)
	}
	if resp.StatusCode == 500 {
		rr := RestroomError{}
		body, _ := io.ReadAll(resp.Body)
		err = json.Unmarshal([]byte(body), &rr)
		if err != nil {
			log.Println(err)
			os.Exit(-1)
		}
		fmt.Println("==== Zenroom result ====")
		fmt.Println(rr.ZenroomError.Result)
		fmt.Println("==== Zenroom logs ====")
		fmt.Println(rr.ZenroomError.Logs)
		fmt.Println("==== Restrooom exceptions  ====")
		fmt.Println(rr.Exception)
		fmt.Println("==== Restrooom result  ====")
		fmt.Println(rr.Result)
		os.Exit(-1)
	}
	if resp.StatusCode != 200 {
		log.Printf("Received status code %d in response\n", resp.StatusCode)
		log.Println(resp)
		os.Exit(-1)
	}
	//var rr map[string]string
	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))
	//_ = json.Unmarshal(body, &rr)
	/*fmt.Println("Result")
	fmt.Println(rr["result"])

	fmt.Println("Logs")
	fmt.Println(rr["result"])*/

	os.Exit(0)
}
