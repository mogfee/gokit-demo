package sum

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"net/http"
	"os"
)

func MakeSumTransport() {
	logger := log.NewLogfmtLogger(os.Stdout)

	filedKeys := []string{"method", "error"}
	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "echina",
		Subsystem: "sum_server",
		Name:      "sum_counter",
		Help:      "sum 请求次数",
	}, filedKeys)
	requestLetency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "echina",
		Subsystem: "sum_server",
		Name:      "sum_letency",
		Help:      "sum 请求次数",
	}, filedKeys)

	var srv SumService
	srv = sumService{}
	srv = ProxyMiddleware(context.Background(), "localhost:8080", logger)(srv)
	srv = loggingMiddleware{logger: logger, next: srv}
	srv = instrumentingMiddleware{requestLatency: requestLetency, requestCount: requestCount, next: srv}

	var sum endpoint.Endpoint
	sum = makeSumEndpoint(srv)
	//sum = logingMiddleware(log.With(logger, "method", "sum"))(sum)

	sumHandler := httptransport.NewServer(sum, encode, decode, )
	router := mux.NewRouter()
	router.Path("/sum").Methods("POST").Handler(sumHandler)
	http.Handle("/sum", sumHandler)
	http.Handle("/metrics", promhttp.Handler())
	logger.Log("stop", http.ListenAndServe(":8080", nil))

}
func encode(ctx context.Context, request2 *http.Request) (request interface{}, err error) {
	var resp SumRequest
	if err := json.NewDecoder(request2.Body).Decode(&resp); err != nil {
		return nil, err
	}
	return resp, nil
}
func decode(ctx context.Context, writer http.ResponseWriter, i interface{}) error {
	return json.NewEncoder(writer).Encode(i)
}
