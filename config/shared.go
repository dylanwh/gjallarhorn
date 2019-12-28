package config

import (
	"os"

	"github.com/pborman/getopt"
)

func keyflag() *string {
	return getopt.StringLong("key", 'k', os.Getenv("GJALLARHORN_KEY"), "secret key for signature of monitor message.")
}

type Keyer interface {
	Key() string
}
