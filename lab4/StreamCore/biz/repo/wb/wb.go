package wb

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type wbRepo interface {
	BatchUpdate(ctx context.Context, batch []interface{}) error
}

type domain interface {
	WbId() string                              // an id for wb logging
	ToWbModel(ctx context.Context) interface{} // convert to model for db
}

// Write-behind Caching strategy
type Strategy struct {
	repo         wbRepo
	updateQueue  chan domain
	batchSize    int
	wg           sync.WaitGroup
	cancelCtx    context.Context
	cancelFunc   context.CancelFunc
	interval     time.Duration
	dbWriteMutex sync.Mutex
	stopOnce     sync.Once
}

type Config struct {
	Repo      wbRepo
	QueueSize int
	BatchSize int
	Interval  time.Duration
}

func NewStrategy(cfg *Config) *Strategy {
	ctx, cancel := context.WithCancel(context.Background())

	s := &Strategy{
		repo:        cfg.Repo,
		updateQueue: make(chan domain, cfg.BatchSize),
		cancelCtx:   ctx,
		cancelFunc:  cancel,
		batchSize:   cfg.BatchSize,
		interval:    cfg.Interval,
	}

	s.wg.Add(1)
	go s.dbBgWorker()
	return s
}

func (s *Strategy) Enqueue(ctx context.Context, d domain) error {
	// add to async queue.
	select {
	case s.updateQueue <- d:
		return nil // success
	default:
		return fmt.Errorf("Update queue overflow for model:%s", d.WbId())
	}
}

func (s *Strategy) Stop() {
	s.stopOnce.Do(func() {
		s.cancelFunc() // signal cancellation
		s.wg.Wait()    // wait for worker goroutine to finish
		// successfully stopped
	})
}

func (s *Strategy) dbBgWorker() {
	defer s.wg.Done()
	batch := make([]interface{}, 0, s.batchSize)
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case d, ok := <-s.updateQueue:
			if !ok { // should not happen if Stop is called correctly
				s.flushBatchToDB(context.Background(), batch)
				return
			}
			batch = append(batch, d.ToWbModel(s.cancelCtx))
			if len(batch) >= s.batchSize { // flush if batch is full
				s.flushBatchToDB(s.cancelCtx, batch)
				batch = batch[:0]
			}

		case <-ticker.C: // write periodically
			if len(batch) > 0 {
				s.flushBatchToDB(s.cancelCtx, batch)
				batch = batch[:0]
			}
		case <-s.cancelCtx.Done():
			s.flushBatchToDB(context.Background(), batch) // clean batch
			draining := true
			for draining {
				select {
				case d, ok := <-s.updateQueue:
					if !ok {
						draining = false
						break
					}
					batch = append(batch, d.ToWbModel(s.cancelCtx))
					if len(batch) > s.batchSize {
						s.flushBatchToDB(context.Background(), batch)
						batch = batch[:0]
					}
				default:
					draining = false // queue empty
				}
				s.flushBatchToDB(context.Background(), batch) // final batch
				return
			}
		}

	}
}

func (s *Strategy) flushBatchToDB(ctx context.Context, batch []interface{}) {
	if len(batch) == 0 {
		return
	}
	s.dbWriteMutex.Lock()
	defer s.dbWriteMutex.Unlock()

	err := s.repo.BatchUpdate(ctx, batch)
	if err != nil {
		// TODO: handle db write failure
	}
}
