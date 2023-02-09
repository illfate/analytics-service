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
	errCh := s.pipeLine.write(events)
	err := <-errCh
	return err
}

func (s *Service) StartCreatingEvents(ctx context.Context) {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		s.pipeLine.read(func(events []Event) error {
			err := s.repo.InsertEvents(ctx, events)
			if err != nil {
				return fmt.Errorf("failed to insert events: %w", err)
			}
			return nil
		})
	}()
}

func (s *Service) Close() {
	s.pipeLine.close()
	s.wg.Wait()
}

type eventsWithChan struct {
	events []Event
	ch     chan<- error
}

type pipeLine struct {
	isClosed    *atomic.Bool
	eventsCh    chan eventsWithChan
	tickerDur   time.Duration
	eventsLimit int
}

func newPipeLine() *pipeLine {
	isClosed := atomic.Bool{}
	isClosed.Store(false)
	// TODO: configure
	return &pipeLine{
		isClosed:    &isClosed,
		eventsCh:    make(chan eventsWithChan, 100),
		tickerDur:   time.Second,
		eventsLimit: 1000,
	}
}

func (p *pipeLine) write(events []Event) <-chan error {
	ch := make(chan error, 1)
	isClosed := p.isClosed.Load()
	if isClosed {
		err := errors.New("pipeline closed")
		ch <- err
		return ch
	}
	p.eventsCh <- eventsWithChan{
		events: events,
		ch:     ch,
	}
	return ch
}

func (p *pipeLine) read(f func(events []Event) error) {
	buf := make([]Event, 0, p.eventsLimit)
	chanBuf := make([]chan<- error, 0, p.eventsLimit)
	done := false
	ticker := time.NewTicker(p.tickerDur)
	defer ticker.Stop()
	flushF := func() {
		if len(buf) != 0 {
			err := f(buf)
			for _, ch := range chanBuf {
				ch <- err
			}
			buf = buf[:0]
			chanBuf = chanBuf[:0]
		}
	}
	for !done {
		select {
		case events, ok := <-p.eventsCh:
			if ok == false {
				done = true
			}
			buf = append(buf, events.events...)
			chanBuf = append(chanBuf, events.ch)
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
