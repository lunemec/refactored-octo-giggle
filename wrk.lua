wrk.scheme = "http"
wrk.host = "localhost"
wrk.port = 8888
wrk.method = "POST"
wrk.path = "/api/v1/buffered"
wrk.headers["Content-Type"] = "application/json"

wrk.body = [[
{
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
	}
]]

