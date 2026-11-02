package utils

import (
	"strings"

	"github.com/mileusna/useragent"
)

type DeviceInfo struct {
	DeviceType string `json:"device_type"` // desktop, mobile, tablet, bot
	Browser    string `json:"browser"`
	OS         string `json:"os"`
}

// ParseUserAgent extracts device, browser, and OS information from the User-Agent header
func ParseUserAgent(uaString string) DeviceInfo {
	if strings.TrimSpace(uaString) == "" {
		return DeviceInfo{
			DeviceType: "unknown",
			Browser:    "unknown",
			OS:         "unknown",
		}
	}

	ua := useragent.Parse(uaString)

	deviceType := "desktop"
	if ua.Bot {
		deviceType = "bot"
	} else if ua.Mobile {
		deviceType = "mobile"
	} else if ua.Tablet {
		deviceType = "tablet"
	}

	browser := ua.Name
	if browser == "" {
		browser = "Other"
	}

	os := ua.OS
	if os == "" {
		os = "Other"
	}

	return DeviceInfo{
		DeviceType: deviceType,
		Browser:    browser,
		OS:         os,
	}
}
