package analytics

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type Repository interface {
	InsertEvents(ctx context.Context, events []Event) error
}

type Service struct {
	repo     Repository
	pipeLine *pipeLine
	wg       sync.WaitGroup
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo, pipeLine: newPipeLine()}
}

func (s *Service) CreateEvents(_ context.Context, events []Event) error {
	return s.pipeLine.write(events)
}

func (s *Service) StartCreatingEvents(ctx context.Context, handleErr func(err error)) {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		s.pipeLine.read(func(events []Event) {
			err := s.repo.InsertEvents(ctx, events)
			if err != nil {
				handleErr(fmt.Errorf("failed to insert events: %w", err))
			}
		})
	}()
}

func (s *Service) Close() {
	s.pipeLine.close()
	s.wg.Wait()
}

type pipeLine struct {
	isClosed    *atomic.Bool
	eventsCh    chan []Event
	tickerDur   time.Duration
	eventsLimit int
}

func newPipeLine() *pipeLine {
	isClosed := atomic.Bool{}
	isClosed.Store(false)
	// TODO: configure
	return &pipeLine{
		isClosed:    &isClosed,
		eventsCh:    make(chan []Event, 100),
		tickerDur:   time.Second,
		eventsLimit: 1000,
	}
}

func (p *pipeLine) write(events []Event) error {
	isClosed := p.isClosed.Load()
	if isClosed {
		return errors.New("pipeline closed")
	}
	p.eventsCh <- events
	return nil
}

func (p *pipeLine) read(f func(events []Event)) {
	buf := make([]Event, 0, p.eventsLimit)
	done := false
	ticker := time.NewTicker(p.tickerDur)
	defer ticker.Stop()
	flushF := func() {
		if len(buf) != 0 {
			f(buf)
			buf = buf[:0]
		}
	}
	for !done {
		select {
		case events, ok := <-p.eventsCh:
			if ok == false {
				done = true
			}
			buf = append(buf, events...)
			if len(buf) > p.eventsLimit {
				flushF()
			}
		case <-ticker.C:
			flushF()
		}
	}
	flushF()
}

func (p *pipeLine) close() {
	p.isClosed.Store(true)
	close(p.eventsCh)
}
