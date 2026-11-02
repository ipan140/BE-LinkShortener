package analytics

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"time"

	appErrors "link-shortener/common/errors"
	"link-shortener/config"
)

type Service interface {
	StartClickCollectorWorker(ctx context.Context)
	GetAnalytics(ctx context.Context, shortCode string) (*AnalyticsSummary, error)
	ExportClicksCSV(ctx context.Context, shortCode string) ([]byte, error)
}

type service struct {
	repo  Repository
	queue RedisQueue
	cfg   *config.AppConfig
}

func NewService(repo Repository, queue RedisQueue, cfg *config.AppConfig) Service {
	return &service{
		repo:  repo,
		queue: queue,
		cfg:   cfg,
	}
}

// StartClickCollectorWorker runs the background worker that drains the Redis queue into the database
func (s *service) StartClickCollectorWorker(ctx context.Context) {
	log.Printf("[Worker] Starting Click Collector background worker (BatchSize: %d, FlushInterval: %v)",
		s.cfg.Worker.BatchSize, s.cfg.Worker.FlushInterval)

	ticker := time.NewTicker(s.cfg.Worker.FlushInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("[Worker] Shutting down Click Collector worker...")
			return
		case <-ticker.C:
			s.flushQueue(ctx)
		}
	}
}

func (s *service) flushQueue(ctx context.Context) {
	items, err := s.queue.PopBatch(ctx, s.cfg.Worker.BatchSize)
	if err != nil {
		log.Printf("[Worker] Error popping batch from Redis: %v", err)
		return
	}

	if len(items) == 0 {
		return
	}

	clicks := make([]Click, len(items))
	for i, it := range items {
		clicks[i] = Click{
			ShortCode:  it.ShortCode,
			IP:         it.IP,
			Country:    it.Country,
			City:       it.City,
			UserAgent:  it.UserAgent,
			DeviceType: it.DeviceType,
			Browser:    it.Browser,
			OS:         it.OS,
			Referrer:   it.Referrer,
			ClickedAt:  it.ClickedAt,
		}
	}

	if err := s.repo.BatchInsertClicks(ctx, clicks); err != nil {
		log.Printf("[Worker] Error persisting %d clicks to database: %v", len(clicks), err)
		return
	}

	log.Printf("[Worker] Successfully processed & flushed %d clicks to database", len(clicks))
}

func (s *service) GetAnalytics(ctx context.Context, shortCode string) (*AnalyticsSummary, error) {
	if shortCode == "" {
		return nil, appErrors.BadRequest("missing short code", nil)
	}

	summary, err := s.repo.GetSummary(ctx, shortCode)
	if err != nil {
		return nil, appErrors.Internal("failed to fetch analytics")
	}

	return summary, nil
}

// ExportClicksCSV generates a CSV file with click details for the requested shortcode
func (s *service) ExportClicksCSV(ctx context.Context, shortCode string) ([]byte, error) {
	if shortCode == "" {
		return nil, appErrors.BadRequest("missing short code", nil)
	}

	clicks, err := s.repo.GetClicksForExport(ctx, shortCode)
	if err != nil {
		return nil, appErrors.Internal("failed to fetch click data for export")
	}

	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)

	// Write header row
	header := []string{
		"ID",
		"ShortCode",
		"ClickedAt (UTC)",
		"IP",
		"Country",
		"City",
		"DeviceType",
		"Browser",
		"OS",
		"Referrer",
	}
	if err := writer.Write(header); err != nil {
		return nil, err
	}

	// Write records
	for _, c := range clicks {
		row := []string{
			fmt.Sprintf("%d", c.ID),
			c.ShortCode,
			c.ClickedAt.UTC().Format(time.RFC3339),
			c.IP,
			c.Country,
			c.City,
			c.DeviceType,
			c.Browser,
			c.OS,
			c.Referrer,
		}
		if err := writer.Write(row); err != nil {
			return nil, err
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
