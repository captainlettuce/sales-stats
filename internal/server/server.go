package server

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/captainlettuce/sales-stats/pkg/types"
	"log/slog"
	"net/http"
)

type OrderAggregationService interface {
	AggregateSalesPrice(ctx context.Context, c types.AggregationRequest) ([]types.AggregationResult, error)
}

type Server struct {
	service OrderAggregationService
}

func New(service OrderAggregationService) *Server {
	return &Server{service: service}
}

func (s *Server) AggregateOrders(w http.ResponseWriter, r *http.Request) {

	var req types.AggregationRequest
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&req)
	if err != nil {
		slog.With(slog.Any("err", err)).Error("error decoding request")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	res, err := s.service.AggregateSalesPrice(r.Context(), req)
	if err != nil {
		switch {

		case errors.Is(err, types.ErrInvalidArgument):
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(err.Error()))
			return

		case errors.Is(err, types.ErrUnexpectedError):
			fallthrough
		default:
			slog.With(slog.Any("err", err), slog.Any("request", req)).Error("got unexpected error aggregating orders")
			// Do not write the error on unknown errors since we might not want to leak them to end user
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	err = enc.Encode(res)
	if err != nil {
		slog.With(slog.Any("err", err)).Error("unexpected error encoding response")
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *Server) CreateHTTPServer(bind string) *http.Server {
	m := http.NewServeMux()

	m.HandleFunc("POST /aggregate", s.AggregateOrders)

	return &http.Server{Addr: bind, Handler: m}
}
