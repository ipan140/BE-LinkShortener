package redirect

import (
	"context"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type mockRepo struct {
	link *RedirectLink
}

func (m *mockRepo) FindActiveLink(ctx context.Context, shortCode string) (*RedirectLink, error) {
	return m.link, nil
}

type mockCache struct {
	cachedURL string
}

func (m *mockCache) GetRedirectURL(ctx context.Context, shortCode string) (string, error) {
	return m.cachedURL, nil
}

func (m *mockCache) SetRedirectCache(ctx context.Context, shortCode, originalURL string, ttl time.Duration) error {
	return nil
}

func (m *mockCache) EnqueueClick(ctx context.Context, click *ClickData) error {
	return nil
}

func TestMergeQueryParams(t *testing.T) {
	dest := "https://example.com/checkout?promo=DISCOUNT10"
	incoming := "utm_source=google&utm_medium=cpc"

	result := mergeQueryParams(dest, incoming)

	expectedContains := []string{
		"promo=DISCOUNT10",
		"utm_source=google",
		"utm_medium=cpc",
	}

	for _, exp := range expectedContains {
		if !contains(result, exp) {
			t.Fatalf("expected result %s to contain %s", result, exp)
		}
	}
}

func TestResolve_SmartDeviceRouting(t *testing.T) {
	repo := &mockRepo{
		link: &RedirectLink{
			ID:          1,
			OriginalURL: "https://example.com/web",
			IOSURL:      "https://apps.apple.com/app/id12345",
			AndroidURL:  "https://play.google.com/store/apps/details?id=com.app",
			ShortCode:   "app-link",
			IsActive:    true,
		},
	}
	cache := &mockCache{}
	service := NewService(repo, cache)

	// 1. iOS visitor
	iosClick := &ClickData{ShortCode: "app-link", OS: "iOS"}
	targetIOS, _, err := service.Resolve(context.Background(), "app-link", "", "", iosClick)
	if err != nil || targetIOS != "https://apps.apple.com/app/id12345" {
		t.Fatalf("expected iOS target, got %s, err: %v", targetIOS, err)
	}

	// 2. Android visitor
	androidClick := &ClickData{ShortCode: "app-link", OS: "Android"}
	targetAndroid, _, err := service.Resolve(context.Background(), "app-link", "", "", androidClick)
	if err != nil || targetAndroid != "https://play.google.com/store/apps/details?id=com.app" {
		t.Fatalf("expected Android target, got %s, err: %v", targetAndroid, err)
	}

	// 3. Desktop visitor
	desktopClick := &ClickData{ShortCode: "app-link", OS: "Windows"}
	targetDesktop, _, err := service.Resolve(context.Background(), "app-link", "", "", desktopClick)
	if err != nil || targetDesktop != "https://example.com/web" {
		t.Fatalf("expected Default target, got %s, err: %v", targetDesktop, err)
	}
}

func TestResolve_PasswordProtected(t *testing.T) {
	password := "secret123"
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	repo := &mockRepo{
		link: &RedirectLink{
			ID:           1,
			OriginalURL:  "https://target.com",
			ShortCode:    "secure1",
			IsActive:     true,
			PasswordHash: string(hashed),
			IsProtected:  true,
		},
	}
	cache := &mockCache{}
	service := NewService(repo, cache)

	click := &ClickData{ShortCode: "secure1"}

	// 1. Without password -> should require password
	_, requiresPassword, err := service.Resolve(context.Background(), "secure1", "", "", click)
	if err != nil || !requiresPassword {
		t.Fatalf("expected requiresPassword=true, got %v, err: %v", requiresPassword, err)
	}

	// 2. With wrong password -> should return unauthorized error
	_, _, err = service.Resolve(context.Background(), "secure1", "", "wrong", click)
	if err == nil {
		t.Fatalf("expected error on wrong password, got nil")
	}

	// 3. With correct password -> should succeed
	target, reqPass, err := service.Resolve(context.Background(), "secure1", "", "secret123", click)
	if err != nil || reqPass || target != "https://target.com" {
		t.Fatalf("expected success, got target=%s, reqPass=%v, err=%v", target, reqPass, err)
	}
}

func TestResolve_MaxClicks(t *testing.T) {
	maxClicks := int64(10)
	repo := &mockRepo{
		link: &RedirectLink{
			ID:          1,
			OriginalURL: "https://target.com",
			ShortCode:   "limit1",
			IsActive:    true,
			MaxClicks:   &maxClicks,
			ClickCount:  10, // already reached max
		},
	}
	cache := &mockCache{}
	service := NewService(repo, cache)

	click := &ClickData{ShortCode: "limit1"}

	_, _, err := service.Resolve(context.Background(), "limit1", "", "", click)
	if err == nil {
		t.Fatalf("expected error when click limit reached, got nil")
	}
}

func contains(str, substr string) bool {
	return len(str) >= len(substr) && (str == substr || findSubstr(str, substr))
}

func findSubstr(str, substr string) bool {
	for i := 0; i+len(substr) <= len(str); i++ {
		if str[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
