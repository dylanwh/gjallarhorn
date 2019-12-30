package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/dylanwh/gjallarhorn/config"
	"github.com/dylanwh/gjallarhorn/message"
)

type handler struct{ config *config.Server }

/*ErrInvalidSignature the error returned when a signature doesn't match. */
var ErrInvalidSignature error = errors.New("Invalid Signature")

/*ErrInvalidJSON ...*/
var ErrInvalidJSON error = errors.New("Invalid JSON")

func main() {
	// systemd gives us timestamps, so remove this.
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	cfg := config.NewServer()
	cfg.CheckArgs()

	http.Handle("/", &handler{config: cfg})
	http.ListenAndServe(cfg.Listen(), nil)
}

func (h *handler) ServeHTTP(out http.ResponseWriter, req *http.Request) {
	msg, err := h.parseMessage(req)
	if err != nil {
		log.Println(err)
		switch {
		case errors.Is(err, ErrInvalidSignature):
			out.WriteHeader(400)
			fmt.Fprintf(out, "invalid signature")
			return
		case errors.Is(err, ErrInvalidJSON):
			out.WriteHeader(400)
			fmt.Fprintf(out, "invalid json")
			return
		default:
			out.WriteHeader(500)
			fmt.Fprintf(out, "internal error")
			return
		}
	}
	fmt.Fprintf(out, "Good job")

	if !msg.KnownIP.Equal(*msg.PublicIP) {
		log.Printf(
			"%s is misconfigured.\n  %s was found\n  %s was expected\n",
			msg.FullHostname,
			msg.KnownIP,
			msg.PublicIP,
		)
		return
	}
	log.Printf("%s (%s) is [%s]\n", msg.Hostname, msg.FullHostname, msg.PublicIP.String())
}

func (h *handler) parseMessage(req *http.Request) (*message.Message, error) {

	var msg message.Message
	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&msg); err != nil {
		log.Printf("json error: %v\n", err)
		return nil, ErrInvalidJSON
	}
	sig, err := msg.Sign(h.config)
	if err != nil {
		return nil, fmt.Errorf("error calculating signature: %v", err)
	}
	if sig != req.Header.Get("Signature") {
		return nil, ErrInvalidSignature
	}

	return &msg, nil
}
