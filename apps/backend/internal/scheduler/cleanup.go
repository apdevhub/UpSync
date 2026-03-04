package scheduler

import (
	"fmt"
	"log/slog"

	"github.com/robfig/cron/v3"
	"github.com/upsync/backend/internal/services"
)

// Scheduler wraps the cron runner.
type Scheduler struct {
	cron *cron.Cron
	svc  *services.FileService
}

// New creates a new Scheduler.
func New(svc *services.FileService) *Scheduler {
	return &Scheduler{
		cron: cron.New(),
		svc:  svc,
	}
}

// Start registers all jobs and starts the scheduler.
// It also runs an initial cleanup immediately on startup.
func (s *Scheduler) Start() {
	// Run cleanup every 15 minutes
	_, err := s.cron.AddFunc("*/15 * * * *", s.runCleanup)
	if err != nil {
		slog.Error("failed to schedule cleanup job", slog.Any("error", err))
		return
	}

	s.cron.Start()
	slog.Info("scheduler started", slog.String("jobs", "cleanup every 15 minutes"))

	// Run once at startup
	go s.runCleanup()
}

// Stop gracefully shuts down the scheduler.
func (s *Scheduler) Stop() {
	ctx := s.cron.Stop()
	<-ctx.Done()
	slog.Info("scheduler stopped")
}

func (s *Scheduler) runCleanup() {
	slog.Info("running expired file cleanup")
	n, err := s.svc.DeleteExpiredFiles()
	if err != nil {
		slog.Error("cleanup error", slog.Any("error", err))
		return
	}
	slog.Info(fmt.Sprintf("cleanup complete: deleted %d expired file(s)", n))
}
