package link

import (
	"context"
	"log"
	"time"

	"gorm.io/gorm"
)

// ArchiverWorker handles periodic cleanup of expired shortlinks
type ArchiverWorker struct {
	db       *gorm.DB
	cache    RedisCache
	interval time.Duration
}

func NewArchiverWorker(db *gorm.DB, cache RedisCache, interval time.Duration) *ArchiverWorker {
	if interval <= 0 {
		interval = 1 * time.Hour
	}
	return &ArchiverWorker{
		db:       db,
		cache:    cache,
		interval: interval,
	}
}

// Start begins periodic background cleanup until context is cancelled
func (w *ArchiverWorker) Start(ctx context.Context) {
	log.Printf("[Archiver] Starting background link archiver & cleanup worker (Interval: %v)", w.interval)
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	// Run an initial sweep on startup
	w.archiveExpiredLinks(ctx)

	for {
		select {
		case <-ctx.Done():
			log.Println("[Archiver] Shutting down link archiver worker...")
			return
		case <-ticker.C:
			w.archiveExpiredLinks(ctx)
		}
	}
}

func (w *ArchiverWorker) archiveExpiredLinks(ctx context.Context) {
	// Find active links that expired more than 30 days ago
	threshold := time.Now().AddDate(0, 0, -30)

	var expiredLinks []Link
	err := w.db.WithContext(ctx).
		Select("id, short_code, expires_at").
		Where("is_active = ? AND expires_at IS NOT NULL AND expires_at < ?", true, threshold).
		Limit(500).
		Find(&expiredLinks).Error

	if err != nil {
		log.Printf("[Archiver] Error querying expired links for archiving: %v", err)
		return
	}

	if len(expiredLinks) == 0 {
		return
	}

	ids := make([]uint, len(expiredLinks))
	for i, l := range expiredLinks {
		ids[i] = l.ID
		// Evict from Redis cache
		if w.cache != nil {
			_ = w.cache.DeleteLink(ctx, l.ShortCode)
		}
	}

	// Deactivate batch
	res := w.db.WithContext(ctx).
		Model(&Link{}).
		Where("id IN ?", ids).
		Update("is_active", false)

	if res.Error != nil {
		log.Printf("[Archiver] Error batch deactivating %d expired links: %v", len(ids), res.Error)
		return
	}

	log.Printf("[Archiver] Successfully archived and deactivated %d stale expired links", res.RowsAffected)
}
