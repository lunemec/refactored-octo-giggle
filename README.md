# API challenge

I decided to create 2 implementations of facet parsing. One is the `BufferedChallengeHandler` version,
which would be my favorite. It is much clearer what the code does, and easier to test, especially
the `Node` struct.

The second implementation uses simple adding of numbers while going over the JSON tokens. It is very
rigid and would be pain to extend and change. But it is faster and supports input JSON as a stream.

To run the API, tests, lints and packaging, you can use the prepared Makefile:
```
 λ make
  Build                          
make build                            Build production binary.                           
make docker                           Build docker.                                      
  Dev                            
make run                              Run API in dev mode, all logging and race detector ON. 
make test                             Run tests.                                         
make vet                              Run go vet.                                        
make lint                             Run gometalinter (you have to install it).   
```

The API endpoints are following (by default the server runs on `0.0.0.0:8888` because of docker):
```
/api/v1/buffered
/api/v1/streaming
```

## Performance comparison
```
 λ benchstat buffered_bench.txt
name                        time/op
BufferedChallengeHandler-4  26.3µs ± 3%

name                        alloc/op
BufferedChallengeHandler-4  13.2kB ± 0%

name                        allocs/op
BufferedChallengeHandler-4     153 ± 0%
```

```
 λ benchstat streaming_bench.txt
name                                 time/op
StreamingChallengeHandler-4          22.9µs ± 4%
StreamingChallengeHandlerParallel-4  23.2µs ± 1%

name                                 alloc/op
StreamingChallengeHandler-4          7.81kB ± 0%
StreamingChallengeHandlerParallel-4  9.51kB ± 0%

name                                 allocs/op
StreamingChallengeHandler-4             197 ± 0%
StreamingChallengeHandlerParallel-4     212 ± 0%
```

There is also `wrk` lua script that can be used to simulate load on the API
and measure performance.
