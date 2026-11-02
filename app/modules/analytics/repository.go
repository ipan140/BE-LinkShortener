package analytics

import (
	"context"
	"time"

	"gorm.io/gorm"
	"link-shortener/app/modules/link"
)

type Repository interface {
	BatchInsertClicks(ctx context.Context, clicks []Click) error
	GetSummary(ctx context.Context, shortCode string) (*AnalyticsSummary, error)
	GetClicksForExport(ctx context.Context, shortCode string) ([]Click, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// BatchInsertClicks bulk inserts click records and atomically updates click counters on links table
func (r *repository) BatchInsertClicks(ctx context.Context, clicks []Click) error {
	if len(clicks) == 0 {
		return nil
	}

	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Bulk insert clicks
		if err := tx.CreateInBatches(clicks, 100).Error; err != nil {
			return err
		}

		// Group delta click count per short_code
		counts := make(map[string]int64)
		for _, c := range clicks {
			counts[c.ShortCode]++
		}

		for code, delta := range counts {
			if err := tx.Model(&link.Link{}).
				Where("short_code = ?", code).
				UpdateColumn("click_count", gorm.Expr("click_count + ?", delta)).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// GetSummary retrieves aggregated analytics for a given short code
func (r *repository) GetSummary(ctx context.Context, shortCode string) (*AnalyticsSummary, error) {
	summary := &AnalyticsSummary{
		ShortCode: shortCode,
	}

	db := r.db.WithContext(ctx).Model(&Click{}).Where("short_code = ?", shortCode)

	// Total Clicks
	if err := db.Count(&summary.TotalClicks).Error; err != nil {
		return nil, err
	}

	// Unique Visitors (distinct IP)
	if err := db.Distinct("ip").Count(&summary.UniqueVisitors).Error; err != nil {
		return nil, err
	}

	// Device breakdown
	var devices []StatItem
	db.Select("device_type as name, count(*) as count").
		Group("device_type").
		Order("count DESC").
		Limit(10).
		Scan(&devices)
	summary.Devices = devices

	// Browser breakdown
	var browsers []StatItem
	db.Select("browser as name, count(*) as count").
		Group("browser").
		Order("count DESC").
		Limit(10).
		Scan(&browsers)
	summary.Browsers = browsers

	// OS breakdown
	var osList []StatItem
	db.Select("os as name, count(*) as count").
		Group("os").
		Order("count DESC").
		Limit(10).
		Scan(&osList)
	summary.OperatingSys = osList

	// Top Referrers
	var referrers []StatItem
	db.Select("referrer as name, count(*) as count").
		Where("referrer != '' AND referrer IS NOT NULL").
		Group("referrer").
		Order("count DESC").
		Limit(10).
		Scan(&referrers)
	summary.TopReferrers = referrers

	// Top Countries
	var countries []StatItem
	db.Select("country as name, count(*) as count").
		Where("country != '' AND country != 'Unknown'").
		Group("country").
		Order("count DESC").
		Limit(10).
		Scan(&countries)
	summary.TopCountries = countries

	// Recent 7 days timeline
	sevenDaysAgo := time.Now().AddDate(0, 0, -7)
	var timeSeries []TimeSeriesItem
	r.db.WithContext(ctx).Raw(`
		SELECT DATE(clicked_at) as date, count(*) as count 
		FROM clicks 
		WHERE short_code = ? AND clicked_at >= ?
		GROUP BY DATE(clicked_at)
		ORDER BY date ASC
	`, shortCode, sevenDaysAgo).Scan(&timeSeries)
	summary.RecentClicks = timeSeries

	return summary, nil
}

// GetClicksForExport retrieves raw clicks for CSV export
func (r *repository) GetClicksForExport(ctx context.Context, shortCode string) ([]Click, error) {
	var clicks []Click
	err := r.db.WithContext(ctx).
		Where("short_code = ?", shortCode).
		Order("clicked_at DESC").
		Limit(10000). // Limit export to top 10,000 for safety
		Find(&clicks).Error
	if err != nil {
		return nil, err
	}
	return clicks, nil
}
