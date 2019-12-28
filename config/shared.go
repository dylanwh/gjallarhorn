package config

import (
	"os"

	"github.com/pborman/getopt"
)

func keyflag() *string {
	return getopt.StringLong("key", 'k', os.Getenv("GJALLARHORN_KEY"), "secret key for signature of monitor message.")
}

/*Keyer are things that can be used to sign messages. */
type Keyer interface {
	Key() string
}
