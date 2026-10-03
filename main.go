package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", IndexHandler)
	// can we use BannerHandler in "/"? ask anas
	http.HandleFunc("/Ascii-art", BannerHandler)
	fmt.Printf(" Server is running successfully at port 8080!\n")
	log.Fatal(http.ListenAndServe(":8080", nil))
	//Fatal is equivalent to Print followed by a call to os.Exit(1).
	// port is the enterance
	// nil is the route
	//sdassd
	//asdasdas
}
