package config

import (
	"os"

	"github.com/pborman/getopt"
)

/*Server is the server-side configuration. */
type Server struct {
	listen *string
	key    *string
}

/*NewServer parses the command line arguments and returns the configuration struct. */
func NewServer() *Server {
	s := &Server{
		listen: getopt.StringLong("listen", 'l', ":8080", "ip and port to listen on"),
		key:    keyflag(),
	}
	getopt.Parse()
	return s
}

/*CheckArgs ensures all required flags are present. If not, it prints a usage
 * message and exits.
 */
func (s *Server) CheckArgs() {
	if *s.key == "" {
		getopt.Usage()
		os.Exit(1)
	}
}

/*Listen returns the host:port string for the http server. */
func (s *Server) Listen() string {
	return *s.listen
}

/*Key returns the shared key used to generate hmac signatures. */
func (s *Server) Key() string {
	return *s.key
}
