package api

import (
	"fmt"
	"io"
	"net/http"
	"sort"
	"time"

	"github.com/gorilla/mux"
	log "github.com/mgutz/logxi/v1"
)

// Config is API server configuration, may contain configuration options
// like timeouts, TLS configuration or other.
type Config struct {
	Address string
	Port    int

	// Timeouts
	ReadTimeout       time.Duration `mapstructure:"read_timeout"`
	ReadHeaderTimeout time.Duration `mapstructure:"read_header_timeout"`
	WriteTimeout      time.Duration `mapstructure:"write_timeout"`
	IdleTimeout       time.Duration `mapstructure:"idle_timeout"`
}

// Addr returns the API listen address (address:port).
func (a *Config) Addr() string {
	return fmt.Sprintf("%s:%d", a.Address, a.Port)
}

// InputJSON represents incomming facets.
type InputJSON struct {
	Data map[string]interface{} `json:"data"`
}

// OutputJSON represents outgoing computed facets.
type OutputJSON struct {
	Result []facetValues `json:"result"`
}

type facetValues map[string]float64

// RunServer runs net/http based API server.
func RunServer(conf Config) error {
	router := mux.NewRouter()
	// Clarify this is API.
	apiRouter := router.PathPrefix("/api").Subrouter()
	// API should be versioned. Period.
	v1Router := apiRouter.PathPrefix("/v1").Subrouter()
	v1Router.Handle("/buffered", panicHandler(ErrHandler(BufferedChallengeHandler))).Methods("POST")
	v1Router.Handle("/streaming", panicHandler(ErrHandler(StreamingChallengeHandler))).Methods("POST")

	server := http.Server{
		Addr:              conf.Addr(),
		Handler:           router,
		ReadTimeout:       conf.ReadTimeout,
		ReadHeaderTimeout: conf.ReadHeaderTimeout,
		WriteTimeout:      conf.WriteTimeout,
		IdleTimeout:       conf.IdleTimeout,
	}
	return server.ListenAndServe()
}

// mapToSlice sorts keys of facets and produces correctly sorted slice of individual
// {"facetN": 100} objects.
func mapToSlice(facets facetValues) (out []facetValues) {
	// Take length of facets map, and sort the keys by name.
	lenFacets := len(facets)
	var (
		i         int
		facetKeys = make([]string, lenFacets)
	)
	for key := range facets {
		facetKeys[i] = key
		i++
	}
	sort.Strings(facetKeys)

	// Iterate over the sorted keys and insert facets into output slice in the
	// sorted order.
	out = make([]facetValues, lenFacets)
	for i, facet := range facetKeys {
		out[i] = facetValues{facet: facets[facet]}
	}

	return
}

// closer serves as utility function to handle errors while closing any closer,
// but namely it is used with req.Body.Close():
// defer closer(req.body)
func closer(c io.Closer) {
	if err := c.Close(); err != nil {
		log.Error("Unable to close io.Closer", "err", err)
	}
}
