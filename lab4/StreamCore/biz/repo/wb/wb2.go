package wb

import (
	"context"
	"sync"
	"time"
)

type cache interface {
	GetCachedValue(ctx context.Context) any
}

type DedupStrategy struct {
	repo         wbRepo
	pendingSync  sync.Map
	wg           sync.WaitGroup
	cancelCtx    context.Context
	cancelFunc   context.CancelFunc
	interval     time.Duration
	batchLimit   int
	dbWriteMutex sync.Mutex
	stopOnce     sync.Once
}

type DedupConfig struct {
	Config
	BatchLimit int
}

func NewDedupStrategy(cfg *DedupConfig) *DedupStrategy {
	ctx, cancel := context.WithCancel(context.Background())
	s := &DedupStrategy{
		repo:       cfg.Repo,
		cancelCtx:  ctx,
		cancelFunc: cancel,
		interval:   cfg.Interval,
		batchLimit: cfg.BatchLimit,
	}
	s.wg.Add(1)
	go s.bgWorker()
	return s
}

func (s *DedupStrategy) SetTask(key uint, c cache) {
	_, dup := s.pendingSync.LoadOrStore(key, c)
	if dup { // given key exists, should not dup
		// log dup.
	}
}

func (s *DedupStrategy) Stop() {
	s.stopOnce.Do(func() {
		s.cancelFunc() // signal cancellation
		s.wg.Wait()    // wait for worker goroutine to finish
		// successfully done
	})
}

func (s *DedupStrategy) bgWorker() {
	defer s.wg.Done()
	batch := make([]interface{}, 0, s.batchLimit)
	keys := make([]uint, 0, s.batchLimit)
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			empty := true
			s.pendingSync.Range(func(k, v any) bool {
				empty = false
				u, c := k.(uint), v.(cache)
				if len(batch) < s.batchLimit {
					batch = append(batch, c.GetCachedValue(s.cancelCtx))
					keys = append(keys, u)
				} else {
					return false // reach limit of this round, abort batching
				}
				return true
			})
			if empty { // no task coming, skip
				continue
			}
			ok := s.dbBatchUpdate(s.cancelCtx, batch)
			batch = batch[:0]
			if ok { // delete completed tasks
				for _, k := range keys {
					s.pendingSync.Delete(k)
				}
			}
			keys = keys[:0]

		case <-s.cancelCtx.Done():
			for {
				empty := true
				s.pendingSync.Range(func(k, v any) bool {
					empty = false
					u, c := k.(uint), v.(cache)
					if len(batch) < s.batchLimit {
						batch = append(batch, c.GetCachedValue(context.Background()))
						keys = append(keys, u)
					} else {
						return false // reach limit of this round, abort batching
					}
					return true
				})
				if empty {
					break
				}
				// ignore db failure because we're cancelling
				_ = s.dbBatchUpdate(context.Background(), batch)
				batch = batch[:0]
				for _, k := range keys {
					s.pendingSync.Delete(k)
				}
				keys = keys[:0]
			}
		}
	}
}

func (s *DedupStrategy) dbBatchUpdate(ctx context.Context, batch []interface{}) bool {
	s.dbWriteMutex.Lock()
	defer s.dbWriteMutex.Unlock()

	err := s.repo.BatchUpdate(ctx, batch)
	if err != nil {
		// TODO: handle db write failure
		return false
	}
	return true
}
