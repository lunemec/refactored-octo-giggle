package api

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

// StreamingChallengeHandler - implementation of challenge using json.Decoder without
// buffering the data first. This version would work best when the input data
// would be streamed. However it is difficult to run concurently, and might
// be difficult to extend with additional functionality.
// The output, however is not streamed since the "result" array is supposed to be
// ordered.
func StreamingChallengeHandler(rw http.ResponseWriter, req *http.Request) error {
	var (
		out OutputJSON
	)

	facetMap, err := unmarshalWithToken(req.Body)
	defer closer(req.Body)

	if err != nil {
		return errors.Wrap(err, "unable to parse facets json")
	}
	out.Result = mapToSlice(facetMap)

	enc := json.NewEncoder(rw)
	return enc.Encode(&out)
}

// nolint: gocyclo
// unmarshalWithToken name is a bit misleading, but I use it to discern this "token type switch"
// version from the "buffered map[string]interface{}" version.
// This version goes over the JSON tokens and saves all the seen facet names in a buffer
// and upon encountering numeric type, it saves the buffered names a keys in a map
// and increases their values by the number seen.
func unmarshalWithToken(reader io.Reader) (facetValues, error) {
	var (
		seenBuf     []string            // all the facets encountered before "count"
		facetValues = make(facetValues) // map of "facetN": 100
	)

	dec := json.NewDecoder(reader)
	for {
		tok, err := dec.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, errors.Wrap(err, "error decoding input data")
		}

		switch v := tok.(type) {
		case json.Delim:
			switch v {
			case '}':
				// Upon closing of JSON object, we remove 1 item from the end of "seen".
				if len(seenBuf) > 1 {
					seenBuf = seenBuf[:len(seenBuf)-1]
				}
			}
		case string:
			switch v {
			// "data" and "count" values are skipped, everything else is facet.
			case "data", "count":
				continue
			default:
				seenBuf = append(seenBuf, v)
			}
		case float64:
			// Increase all the seen facet's values by v.
			for _, facet := range seenBuf {
				facetValues[facet] += v
			}
		}
	}

	return facetValues, nil
}
