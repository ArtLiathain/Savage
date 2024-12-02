package main

import (
	"flag"
	// "net"
    // "fmt"
	// "strconv"
	// "time"
)

var (
	portFlag   = flag.Int("p", 8080, "Port number for the web API")
	serverFlag = flag.Bool("s", false, "Port number for the web API")
	clientFlag = flag.Bool("c", false, "Port number for the web API")
	autoFlag   = flag.Bool("a", false, "Port number for web API")
)

func main() {
	flag.Parse()
	// if *serverFlag || *autoFlag {
	// 	timeout := time.Second
	// 	conn, err := net.DialTimeout("tcp", net.JoinHostPort("0.0.0.0", strconv.Itoa(*portFlag)), timeout)
	// 	if err != nil {
	// 		initWebServer(*portFlag)
	// 		return
	// 	}
	// 	if conn != nil {
	// 		defer conn.Close()
	// 		fmt.Println("Opened", net.JoinHostPort("0.0.0.0", string(*portFlag)))
	// 	}
	// }


	// if *clientFlag || *autoFlag {
	// 	initClient(*portFlag)
	// }


	initWebApi();
	// initClient();
	// initWebServer();
}
