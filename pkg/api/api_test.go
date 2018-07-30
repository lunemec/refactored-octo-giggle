package api_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"refactored-octo-giggle/pkg/api"

	"github.com/stretchr/testify/assert"
)

var (
	testBody = `{
		"data": {
			"facet1": {
				"facet3": {
					"facet4": {
						"facet6": {
							"count": 20
						},
						"facet7": {
							"count": 30
						}
					},
					"facet5": {
						"count": 50
					}
				}
			}, 
			"facet2": {
				"count": 0
			}
		}
	}`

	expectedOutput = `{
        "result": [
            {"facet1": 100},
            {"facet2": 0},
            {"facet3": 100},
            {"facet4": 50},
            {"facet5": 50},
            {"facet6": 20},
            {"facet7": 30}
        ]
    }`
)

func TestStreamingChallengeHandler(t *testing.T) {
	req, err := http.NewRequest("POST", "/api/v1/challenge", strings.NewReader(testBody))

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	http.Handler(api.ErrHandler(api.StreamingChallengeHandler)).ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Status code differs. Expected %d .\n Got %d instead", http.StatusOK, status)
	}

	assert.JSONEq(t, expectedOutput, rr.Body.String(), "Response body differs")
}

func TestStreamingChallengeHandlerNoData(t *testing.T) {
	req, err := http.NewRequest("POST", "/api/v1/challenge", strings.NewReader(`{"data": {}}`))

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	http.Handler(api.ErrHandler(api.StreamingChallengeHandler)).ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Status code differs. Expected %d .\n Got %d instead", http.StatusOK, status)
	}

	assert.JSONEq(t, `{"result": []}`, rr.Body.String(), "Response body differs")
}

func TestBufferedChallengeHandler(t *testing.T) {
	req, err := http.NewRequest("POST", "/api/v1/challenge2", strings.NewReader(testBody))

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	http.Handler(api.ErrHandler(api.BufferedChallengeHandler)).ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Status code differs. Expected %d .\n Got %d instead", http.StatusOK, status)
	}

	assert.JSONEq(t, expectedOutput, rr.Body.String(), "Response body differs")
}

func TestBufferedChallengeHandlerNoData(t *testing.T) {
	req, err := http.NewRequest("POST", "/api/v1/challenge2", strings.NewReader(`{"data": {}}`))

	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	http.Handler(api.ErrHandler(api.BufferedChallengeHandler)).ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Status code differs. Expected %d .\n Got %d instead", http.StatusOK, status)
	}

	assert.JSONEq(t, `{"result": []}`, rr.Body.String(), "Response body differs")
}

func BenchmarkStreamingChallengeHandler(b *testing.B) {
	b.StopTimer()

	for i := 0; i < b.N; i++ {
		req, err := http.NewRequest("POST", "/api/v1/challenge", strings.NewReader(testBody))

		if err != nil {
			b.Fatal(err)
		}
		rr := httptest.NewRecorder()
		b.StartTimer()
		http.Handler(api.ErrHandler(api.StreamingChallengeHandler)).ServeHTTP(rr, req)
		b.StopTimer()
		if status := rr.Code; status != http.StatusOK {
			b.Errorf("Status code differs. Expected %d .\n Got %d instead", http.StatusOK, status)
		}

		assert.JSONEq(b, expectedOutput, rr.Body.String(), "Response body differs")
	}
}

func BenchmarkStreamingChallengeHandlerParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			req, err := http.NewRequest("POST", "/api/v1/challenge", strings.NewReader(testBody))

			if err != nil {
				b.Fatal(err)
			}
			rr := httptest.NewRecorder()
			b.StartTimer()
			http.Handler(api.ErrHandler(api.StreamingChallengeHandler)).ServeHTTP(rr, req)
			b.StopTimer()
			if status := rr.Code; status != http.StatusOK {
				b.Errorf("Status code differs. Expected %d .\n Got %d instead", http.StatusOK, status)
			}

			assert.JSONEq(b, expectedOutput, rr.Body.String(), "Response body differs")
		}
	})
}

func BenchmarkBufferedChallengeHandler(b *testing.B) {
	b.StopTimer()

	for i := 0; i < b.N; i++ {
		req, err := http.NewRequest("POST", "/api/v1/challenge", strings.NewReader(testBody))

		if err != nil {
			b.Fatal(err)
		}
		rr := httptest.NewRecorder()
		b.StartTimer()
		http.Handler(api.ErrHandler(api.BufferedChallengeHandler)).ServeHTTP(rr, req)
		b.StopTimer()
		if status := rr.Code; status != http.StatusOK {
			b.Errorf("Status code differs. Expected %d .\n Got %d instead", http.StatusOK, status)
		}

		assert.JSONEq(b, expectedOutput, rr.Body.String(), "Response body differs")
	}
}
