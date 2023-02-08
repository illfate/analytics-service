package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"go.uber.org/zap"

	"github.com/illfate/analytics-service/internal/analytics"
)

type Server struct {
	http.Handler
	service *analytics.Service
	logger  *zap.Logger
}

func NewServer(service *analytics.Service, logger *zap.Logger) *Server {
	router := chi.NewRouter()
	s := &Server{
		Handler: router,
		service: service,
		logger:  logger,
	}
	router.Post("/v1/events", s.createEvents)
	return s
}

func (s *Server) createEvents(w http.ResponseWriter, req *http.Request) {
	events, err := s.eventsFromReq(req)
	if err != nil {
		s.renderErr(w, req, ErrInvalidRequest(err))
		return
	}
	err = s.service.CreateEvents(req.Context(), events)
	if err != nil {
		s.logger.Error("Failed to create events", zap.Error(err))
		s.renderErr(w, req, ErrRender(err))
		return
	}
	render.Status(req, http.StatusOK)
}

func (s *Server) eventsFromReq(req *http.Request) ([]analytics.Event, error) {
	body, err := io.ReadAll(req.Body)
	if err != nil {

		return nil, err
	}
	jsonObjects := bytes.Split(body, []byte("\n"))
	events := make([]analytics.Event, 0, len(jsonObjects))
	for _, o := range jsonObjects {
		if len(o) == 0 {
			continue
		}
		var event analytics.Event
		err = json.Unmarshal(o, &event)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}

func (s *Server) renderErr(w http.ResponseWriter, req *http.Request, rErr render.Renderer) {
	err := render.Render(w, req, rErr)
	if err != nil {
		s.logger.Error("Failed to render err", zap.Error(err))
	}
}
