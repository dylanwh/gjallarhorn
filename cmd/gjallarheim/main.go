package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dylanwh/gjallarhorn/message"
	"github.com/pborman/getopt"
)

func main() {
	var listenFlag = getopt.StringLong("listen", 'l', ":8080", "ip and port to listen on")
	getopt.Parse()

	http.HandleFunc("/", Index)
	http.ListenAndServe(*listenFlag, nil)
}

func Index(w http.ResponseWriter, r *http.Request) {
	var msg message.Message
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&msg); err != nil {
		w.WriteHeader(400)
		fmt.Fprintln(w, err)
	}
	fmt.Printf("%+v\n", msg)
	fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
}
