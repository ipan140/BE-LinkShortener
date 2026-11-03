package web

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// RegisterDashboardRoutes mounts the interactive Web Dashboard at "/" and "/dashboard"
func RegisterDashboardRoutes(e *echo.Echo, baseURL string) {
	handler := func(c echo.Context) error {
		c.Response().Header().Set("Content-Type", "text/html; charset=utf-8")
		return c.String(http.StatusOK, dashboardHTML)
	}

	e.GET("/", handler)
	e.GET("/dashboard", handler)
}

const dashboardHTML = `<!DOCTYPE html>
<html lang="en" class="dark">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>UltraFast Redirector & Telemetry Engine</title>
  <script src="https://cdn.tailwindcss.com"></script>
  <link rel="preconnect" href="https://fonts.googleapis.com">
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
  <link href="https://fonts.googleapis.com/css2?family=Plus+Jakarta+Sans:wght@400;500;600;700;800&family=JetBrains+Mono:wght@400;500&display=swap" rel="stylesheet">
  <script>
    tailwind.config = {
      darkMode: 'class',
      theme: {
        extend: {
          fontFamily: {
            sans: ['"Plus Jakarta Sans"', 'sans-serif'],
            mono: ['"JetBrains Mono"', 'monospace'],
          },
          colors: {
            brand: {
              50: '#f0f9ff',
              100: '#e0f2fe',
              400: '#38bdf8',
              500: '#0ea5e9',
              600: '#0284c7',
              700: '#0369a1',
            }
          }
        }
      }
    }
  </script>
  <style>
    body {
      background-color: #090d16;
      background-image: 
        radial-gradient(at 0% 0%, rgba(14, 165, 233, 0.12) 0px, transparent 50%),
        radial-gradient(at 100% 100%, rgba(56, 189, 248, 0.08) 0px, transparent 50%);
    }
    .glass-card {
      background: rgba(15, 23, 42, 0.75);
      backdrop-filter: blur(12px);
      border: 1px solid rgba(255, 255, 255, 0.08);
    }
    .glass-input {
      background: rgba(2, 6, 23, 0.6);
      border: 1px solid rgba(255, 255, 255, 0.1);
    }
    .glass-input:focus {
      border-color: #38bdf8;
      outline: none;
      box-shadow: 0 0 0 2px rgba(56, 189, 248, 0.2);
    }
  </style>
</head>
<body class="text-slate-100 font-sans min-h-screen flex flex-col antialiased selection:bg-brand-500 selection:text-white">

  <!-- Top Navigation Header -->
  <header class="border-b border-slate-800/80 sticky top-0 z-50 bg-[#090d16]/80 backdrop-blur-md">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 h-16 flex items-center justify-between">
      <div class="flex items-center space-x-3">
        <div class="w-10 h-10 rounded-xl bg-gradient-to-tr from-brand-600 to-sky-400 flex items-center justify-center shadow-lg shadow-sky-500/20">
          <svg class="w-6 h-6 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M13 10V3L4 14h7v7l9-11h-7z" />
          </svg>
        </div>
        <div>
          <span class="text-lg font-extrabold tracking-tight bg-gradient-to-r from-white via-slate-100 to-sky-300 bg-clip-text text-transparent">UltraFast</span>
          <span class="text-xs px-2 py-0.5 ml-1.5 rounded-full bg-sky-500/10 text-sky-400 font-medium border border-sky-500/20">&lt;10ms SLA</span>
        </div>
      </div>

      <div class="flex items-center space-x-3">
        <div class="hidden sm:flex items-center space-x-2 px-3 py-1.5 rounded-full bg-emerald-500/10 border border-emerald-500/20 text-emerald-400 text-xs font-medium">
          <span class="w-2 h-2 rounded-full bg-emerald-400 animate-pulse"></span>
          <span>PostgreSQL & Redis Active</span>
        </div>

        <a href="/swagger/index.html" target="_blank" class="px-3 py-1.5 rounded-lg text-xs font-medium text-slate-300 bg-slate-800/80 hover:bg-slate-700 hover:text-white border border-slate-700/60 transition flex items-center space-x-1.5">
          <svg class="w-4 h-4 text-sky-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
          </svg>
          <span>Swagger UI</span>
        </a>

        <a href="/metrics" target="_blank" class="px-3 py-1.5 rounded-lg text-xs font-medium text-slate-300 bg-slate-800/80 hover:bg-slate-700 hover:text-white border border-slate-700/60 transition flex items-center space-x-1.5">
          <svg class="w-4 h-4 text-amber-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 8v8m-4-5v5m-4-2v2m-2 4h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
          </svg>
          <span>Metrics</span>
        </a>
      </div>
    </div>
  </header>

  <!-- Hero Content -->
  <main class="flex-1 max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-10 w-full">
    <div class="text-center max-w-3xl mx-auto mb-10">
      <h1 class="text-4xl sm:text-5xl font-extrabold tracking-tight text-white mb-4">
        High-Performance <span class="bg-gradient-to-r from-sky-400 to-indigo-400 bg-clip-text text-transparent">Redirector</span> & Click Collector
      </h1>
      <p class="text-slate-400 text-sm sm:text-base leading-relaxed">
        Production-ready Go modular monolith. Delivers sub-10ms redirects via Redis, collects real-time visitor telemetry to PostgreSQL, and supports device & geo routing.
      </p>
    </div>

    <!-- Main Grid: Shortener on Left, Analytics on Right -->
    <div class="grid grid-cols-1 lg:grid-cols-12 gap-8 items-start">
      
      <!-- LEFT COLUMN: Generator & Controls (7 cols) -->
      <div class="lg:col-span-7 space-y-6">
        <div class="glass-card rounded-2xl p-6 sm:p-8 shadow-2xl relative overflow-hidden">
          <div class="absolute -right-16 -top-16 w-48 h-48 bg-sky-500/10 rounded-full blur-3xl pointer-events-none"></div>

          <div class="flex items-center justify-between mb-6">
            <div class="flex items-center space-x-2.5">
              <span class="w-3 h-3 rounded-full bg-sky-400"></span>
              <h2 class="text-lg font-bold text-white tracking-tight">Create Ultra-Fast Link</h2>
            </div>
            <span class="text-xs text-slate-400">REST API v1</span>
          </div>

          <form id="shortenForm" class="space-y-4" onsubmit="handleCreateLink(event)">
            <!-- Destination URL -->
            <div>
              <label class="block text-xs font-semibold text-slate-300 uppercase tracking-wider mb-2">Destination URL <span class="text-sky-400">*</span></label>
              <div class="relative">
                <input type="url" id="originalUrl" required placeholder="https://example.com/long-promotional-link"
                  class="glass-input w-full px-4 py-3 rounded-xl text-sm text-white placeholder-slate-500 transition font-mono pr-20" />
                <button type="button" onclick="pasteFromClipboard()" class="absolute right-2 top-2 px-3 py-1.5 rounded-lg text-xs font-medium text-slate-400 hover:text-white bg-slate-800/80 hover:bg-slate-700 transition">
                  Paste
                </button>
              </div>
            </div>

            <!-- Custom Alias & Title -->
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <div>
                <label class="block text-xs font-semibold text-slate-300 uppercase tracking-wider mb-2">Custom Alias (Optional)</label>
                <div class="flex rounded-xl glass-input overflow-hidden">
                  <span class="px-3 py-2.5 text-xs text-slate-500 bg-slate-900/60 font-mono select-none flex items-center">/</span>
                  <input type="text" id="customCode" placeholder="flash-sale" pattern="[a-zA-Z0-9_-]{3,30}"
                    class="bg-transparent w-full px-3 py-2 text-sm text-white placeholder-slate-600 focus:outline-none font-mono" />
                </div>
              </div>
              <div>
                <label class="block text-xs font-semibold text-slate-300 uppercase tracking-wider mb-2">Campaign Title</label>
                <input type="text" id="linkTitle" placeholder="Black Friday Campaign"
                  class="glass-input w-full px-4 py-2.5 rounded-xl text-sm text-white placeholder-slate-600 transition" />
              </div>
            </div>

            <!-- Accordion 1: Smart Device & Geo Routing -->
            <div class="pt-2">
              <button type="button" onclick="toggleSection('routingSection', 'routingIcon')" class="w-full flex items-center justify-between py-2 text-xs font-semibold text-slate-400 hover:text-sky-400 transition">
                <span class="flex items-center space-x-2">
                  <svg class="w-4 h-4 text-sky-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3.055 11H5a2 2 0 012 2v1a2 2 0 002 2 2 2 0 012 2v2.945M8 3.935V5.5A2.5 2.5 0 0010.5 8h.5a2 2 0 012 2 2 2 0 104 0 2 2 0 012-2h1.064M15 20.488V18a2 2 0 012-2h3.064" />
                  </svg>
                  <span>Smart Device & Geo-Targeting</span>
                </span>
                <svg id="routingIcon" class="w-4 h-4 transition-transform duration-200" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                </svg>
              </button>

              <div id="routingSection" class="hidden space-y-3 pt-3 border-t border-slate-800">
                <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
                  <div>
                    <label class="block text-xs text-slate-400 mb-1">iOS App Store URL</label>
                    <input type="url" id="iosUrl" placeholder="https://apps.apple.com/app/id..."
                      class="glass-input w-full px-3 py-2 rounded-lg text-xs text-white placeholder-slate-600 font-mono" />
                  </div>
                  <div>
                    <label class="block text-xs text-slate-400 mb-1">Android Play Store URL</label>
                    <input type="url" id="androidUrl" placeholder="https://play.google.com/store/apps/..."
                      class="glass-input w-full px-3 py-2 rounded-lg text-xs text-white placeholder-slate-600 font-mono" />
                  </div>
                </div>

                <div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
                  <div>
                    <label class="block text-xs text-slate-400 mb-1">Country Code (ISO-2)</label>
                    <input type="text" id="countryCode" placeholder="ID, US, SG" maxlength="2"
                      class="glass-input w-full px-3 py-2 rounded-lg text-xs text-white placeholder-slate-600 font-mono uppercase text-center" />
                  </div>
                  <div class="sm:col-span-2">
                    <label class="block text-xs text-slate-400 mb-1">Country Localized URL</label>
                    <input type="url" id="countryUrl" placeholder="https://example.com/id/promo"
                      class="glass-input w-full px-3 py-2 rounded-lg text-xs text-white placeholder-slate-600 font-mono" />
                  </div>
                </div>
              </div>
            </div>

            <!-- Accordion 2: Advanced Controls & Security -->
            <div>
              <button type="button" onclick="toggleSection('advancedSection', 'advancedIcon')" class="w-full flex items-center justify-between py-2 text-xs font-semibold text-slate-400 hover:text-sky-400 transition">
                <span class="flex items-center space-x-2">
                  <svg class="w-4 h-4 text-emerald-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
                  </svg>
                  <span>Advanced Controls, Webhooks & Fallback URL</span>
                </span>
                <svg id="advancedIcon" class="w-4 h-4 transition-transform duration-200" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                </svg>
              </button>

              <div id="advancedSection" class="hidden space-y-3 pt-3 border-t border-slate-800">
                <div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
                  <div>
                    <label class="block text-xs text-slate-400 mb-1">Access Password</label>
                    <input type="password" id="linkPassword" placeholder="Optional password"
                      class="glass-input w-full px-3 py-2 rounded-lg text-xs text-white placeholder-slate-600" />
                  </div>
                  <div>
                    <label class="block text-xs text-slate-400 mb-1">Max Clicks Limit</label>
                    <input type="number" id="maxClicks" min="1" placeholder="e.g. 500"
                      class="glass-input w-full px-3 py-2 rounded-lg text-xs text-white placeholder-slate-600 font-mono" />
                  </div>
                  <div>
                    <label class="block text-xs text-slate-400 mb-1">Expiry Date</label>
                    <input type="datetime-local" id="expiresAt"
                      class="glass-input w-full px-3 py-2 rounded-lg text-xs text-white placeholder-slate-600" />
                  </div>
                </div>

                <div>
                  <label class="block text-xs text-slate-400 mb-1">Custom Fallback URL (if expired or limit reached)</label>
                  <input type="url" id="fallbackUrl" placeholder="https://example.com/expired-fallback"
                    class="glass-input w-full px-3 py-2 rounded-lg text-xs text-white placeholder-slate-600 font-mono" />
                </div>

                <div>
                  <label class="block text-xs text-slate-400 mb-1">Webhook URL (Discord / Slack / Telegram)</label>
                  <input type="url" id="webhookUrl" placeholder="https://discord.com/api/webhooks/..."
                    class="glass-input w-full px-3 py-2 rounded-lg text-xs text-white placeholder-slate-600 font-mono" />
                </div>
              </div>
            </div>

            <!-- Submit Button -->
            <button type="submit" id="submitBtn" class="w-full py-3.5 px-6 rounded-xl font-bold text-sm text-white bg-gradient-to-r from-sky-500 to-brand-600 hover:from-sky-400 hover:to-brand-500 transition shadow-lg shadow-sky-500/25 flex items-center justify-center space-x-2">
              <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1" />
              </svg>
              <span>Shorten URL Now</span>
            </button>
          </form>

          <!-- Result Box -->
          <div id="resultBox" class="hidden mt-6 p-5 rounded-xl bg-slate-900/90 border border-sky-500/40">
            <div class="flex items-center justify-between mb-3">
              <span class="text-xs font-semibold text-sky-400 uppercase tracking-wider flex items-center space-x-1.5">
                <span class="w-2 h-2 rounded-full bg-sky-400 animate-ping"></span>
                <span>Shortlink Generated!</span>
              </span>
              <span id="createdCodeBadge" class="text-xs font-mono px-2 py-0.5 rounded bg-slate-800 text-slate-300"></span>
            </div>

            <div class="flex flex-col sm:flex-row items-center gap-4">
              <div class="flex-1 w-full">
                <input type="text" id="shortUrlDisplay" readonly
                  class="w-full px-3.5 py-2.5 rounded-lg bg-slate-950 border border-slate-700 text-sky-300 font-mono text-sm select-all" />
                
                <div class="flex items-center space-x-2 mt-3">
                  <button onclick="copyToClipboard()" id="copyBtn" class="px-4 py-2 rounded-lg text-xs font-semibold bg-sky-500 hover:bg-sky-400 text-white transition flex items-center space-x-1.5">
                    <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
                    </svg>
                    <span id="copyBtnText">Copy Link</span>
                  </button>
                  <a id="testLinkBtn" href="#" target="_blank" class="px-4 py-2 rounded-lg text-xs font-semibold bg-slate-800 hover:bg-slate-700 text-slate-200 transition">
                    Visit Link &rarr;
                  </a>
                  <button onclick="loadAnalyticsForCreated()" class="px-3 py-2 rounded-lg text-xs font-medium text-slate-400 hover:text-white bg-slate-800 transition">
                    View Stats
                  </button>
                </div>
              </div>

              <!-- QR Code Thumbnail -->
              <div class="w-24 h-24 sm:w-28 sm:h-28 rounded-xl bg-white p-1.5 flex-shrink-0 shadow-md">
                <img id="qrCodeImg" src="" alt="QR Code" class="w-full h-full object-contain" />
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- RIGHT COLUMN: Real-Time Analytics Inspector (5 cols) -->
      <div class="lg:col-span-5 space-y-6">
        <div class="glass-card rounded-2xl p-6 sm:p-8 shadow-2xl">
          <div class="flex items-center justify-between mb-4">
            <div class="flex items-center space-x-2.5">
              <span class="w-3 h-3 rounded-full bg-emerald-400"></span>
              <h2 class="text-lg font-bold text-white tracking-tight">Analytics Explorer</h2>
            </div>
            <span class="text-xs text-slate-400">Live Telemetry</span>
          </div>

          <div class="flex space-x-2 mb-6">
            <input type="text" id="inspectCode" placeholder="Enter short code (e.g. flash-sale)"
              class="glass-input flex-1 px-3.5 py-2 rounded-xl text-xs sm:text-sm text-white placeholder-slate-500 font-mono" />
            <button type="button" onclick="inspectAnalytics()" class="px-4 py-2 rounded-xl text-xs font-bold text-white bg-slate-800 hover:bg-slate-700 border border-slate-700 transition">
              Inspect
            </button>
          </div>

          <!-- Empty State -->
          <div id="analyticsEmptyState" class="text-center py-12 px-4 border border-dashed border-slate-800 rounded-xl">
            <svg class="w-12 h-12 text-slate-600 mx-auto mb-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
            </svg>
            <p class="text-xs text-slate-400">Shorten a URL or input an existing shortcode to inspect click counts, device metrics, and geo distributions.</p>
          </div>

          <!-- Active Analytics Display -->
          <div id="analyticsContent" class="hidden space-y-4">
            <div class="grid grid-cols-2 gap-3">
              <div class="p-3.5 rounded-xl bg-slate-900/80 border border-slate-800">
                <span class="text-[10px] font-semibold text-slate-400 uppercase tracking-wider block">Total Clicks</span>
                <span id="statTotalClicks" class="text-2xl font-extrabold text-white mt-1 block">0</span>
              </div>
              <div class="p-3.5 rounded-xl bg-slate-900/80 border border-slate-800">
                <span class="text-[10px] font-semibold text-slate-400 uppercase tracking-wider block">Unique Visitors</span>
                <span id="statUniqueVisitors" class="text-2xl font-extrabold text-sky-400 mt-1 block">0</span>
              </div>
            </div>

            <!-- Top Devices -->
            <div class="p-4 rounded-xl bg-slate-900/80 border border-slate-800">
              <h3 class="text-xs font-bold text-slate-300 uppercase tracking-wider mb-2.5">Devices & Browsers</h3>
              <div id="deviceList" class="space-y-2 text-xs"></div>
            </div>

            <!-- Top Countries -->
            <div class="p-4 rounded-xl bg-slate-900/80 border border-slate-800">
              <h3 class="text-xs font-bold text-slate-300 uppercase tracking-wider mb-2.5">Top Geographic Locations</h3>
              <div id="countryList" class="space-y-2 text-xs"></div>
            </div>

            <!-- CSV Export button -->
            <a id="exportCsvBtn" href="#" target="_blank" class="w-full py-2.5 px-4 rounded-xl text-xs font-semibold text-slate-300 bg-slate-800/80 hover:bg-slate-700 border border-slate-700 transition flex items-center justify-center space-x-2">
              <svg class="w-4 h-4 text-emerald-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
              </svg>
              <span>Download Click Telemetry (CSV)</span>
            </a>
          </div>

        </div>
      </div>

    </div>
  </main>

  <footer class="border-t border-slate-800/60 py-6 text-center text-xs text-slate-500">
    <div class="max-w-7xl mx-auto px-4 flex flex-col sm:flex-row items-center justify-between gap-2">
      <p>UltraFast Redirector Engine &bull; Go Clean Architecture &bull; PostgreSQL &bull; Redis</p>
      <div class="flex items-center space-x-4">
        <a href="/swagger/index.html" class="hover:text-slate-300 transition">API Docs</a>
        <a href="/metrics" class="hover:text-slate-300 transition">Prometheus</a>
        <a href="/health" class="hover:text-slate-300 transition">Health</a>
      </div>
    </div>
  </footer>

  <script>
    let latestCreatedCode = '';

    function toggleSection(sectionId, iconId) {
      const el = document.getElementById(sectionId);
      const icon = document.getElementById(iconId);
      if (el.classList.contains('hidden')) {
        el.classList.remove('hidden');
        icon.classList.add('rotate-180');
      } else {
        el.classList.add('hidden');
        icon.classList.remove('rotate-180');
      }
    }

    async function pasteFromClipboard() {
      try {
        const text = await navigator.clipboard.readText();
        if (text) {
          document.getElementById('originalUrl').value = text;
        }
      } catch (e) {
        // clipboard access denied or not supported
      }
    }

    async function handleCreateLink(event) {
      event.preventDefault();
      const submitBtn = document.getElementById('submitBtn');
      submitBtn.disabled = true;
      submitBtn.innerText = 'Creating Shortlink...';

      const payload = {
        original_url: document.getElementById('originalUrl').value.trim(),
        custom_code: document.getElementById('customCode').value.trim() || undefined,
        title: document.getElementById('linkTitle').value.trim() || undefined,
        ios_url: document.getElementById('iosUrl').value.trim() || undefined,
        android_url: document.getElementById('androidUrl').value.trim() || undefined,
        country_code: document.getElementById('countryCode').value.trim() || undefined,
        country_url: document.getElementById('countryUrl').value.trim() || undefined,
        password: document.getElementById('linkPassword').value || undefined,
        fallback_url: document.getElementById('fallbackUrl').value.trim() || undefined,
        webhook_url: document.getElementById('webhookUrl').value.trim() || undefined,
      };

      const maxClicksVal = document.getElementById('maxClicks').value;
      if (maxClicksVal) payload.max_clicks = parseInt(maxClicksVal, 10);

      const expiresVal = document.getElementById('expiresAt').value;
      if (expiresVal) payload.expires_at = new Date(expiresVal).toISOString();

      try {
        const res = await fetch('/api/v1/links/public', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(payload)
        });

        const data = await res.json();
        if (!res.ok) {
          alert('Error: ' + (data.message || 'Failed to create shortlink'));
          return;
        }

        const link = data.data;
        latestCreatedCode = link.short_code;

        document.getElementById('createdCodeBadge').innerText = '/' + link.short_code;
        document.getElementById('shortUrlDisplay').value = link.short_url;
        document.getElementById('testLinkBtn').href = link.short_url;
        document.getElementById('qrCodeImg').src = '/qr/' + link.short_code;
        document.getElementById('resultBox').classList.remove('hidden');

        // Automatically prefill inspector
        document.getElementById('inspectCode').value = link.short_code;
      } catch (err) {
        alert('Request failed: ' + err.message);
      } finally {
        submitBtn.disabled = false;
        submitBtn.innerHTML = '<svg class="w-5 h-5 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1" /></svg>Shorten URL Now';
      }
    }

    function copyToClipboard() {
      const input = document.getElementById('shortUrlDisplay');
      input.select();
      navigator.clipboard.writeText(input.value);
      const textSpan = document.getElementById('copyBtnText');
      textSpan.innerText = 'Copied!';
      setTimeout(() => textSpan.innerText = 'Copy Link', 2000);
    }

    function loadAnalyticsForCreated() {
      if (latestCreatedCode) {
        document.getElementById('inspectCode').value = latestCreatedCode;
        inspectAnalytics();
      }
    }

    async function inspectAnalytics() {
      const code = document.getElementById('inspectCode').value.trim();
      if (!code) {
        alert('Please enter a short code');
        return;
      }

      try {
        const res = await fetch('/api/v1/analytics/public/' + encodeURIComponent(code));
        const data = await res.json();
        if (!res.ok) {
          alert('Analytics Error: ' + (data.message || 'Not found'));
          return;
        }

        const summary = data.data;
        document.getElementById('analyticsEmptyState').classList.add('hidden');
        document.getElementById('analyticsContent').classList.remove('hidden');

        document.getElementById('statTotalClicks').innerText = (summary.total_clicks || 0).toLocaleString();
        document.getElementById('statUniqueVisitors').innerText = (summary.unique_clicks || 0).toLocaleString();

        // Render Devices
        const deviceList = document.getElementById('deviceList');
        deviceList.innerHTML = '';
        if (summary.device_breakdown && summary.device_breakdown.length > 0) {
          summary.device_breakdown.forEach(item => {
            deviceList.innerHTML += '<div class="flex justify-between items-center py-1 border-b border-slate-800/40"><span class="text-slate-400 capitalize">' + (item.device_type || 'Unknown') + '</span><span class="font-mono font-bold text-slate-200">' + item.count + '</span></div>';
          });
        } else {
          deviceList.innerHTML = '<span class="text-slate-500 text-xs">No click logs recorded yet</span>';
        }

        // Render Countries
        const countryList = document.getElementById('countryList');
        countryList.innerHTML = '';
        if (summary.country_breakdown && summary.country_breakdown.length > 0) {
          summary.country_breakdown.forEach(item => {
            countryList.innerHTML += '<div class="flex justify-between items-center py-1 border-b border-slate-800/40"><span class="text-slate-400">' + (item.country || 'Unknown') + '</span><span class="font-mono font-bold text-slate-200">' + item.count + '</span></div>';
          });
        } else {
          countryList.innerHTML = '<span class="text-slate-500 text-xs">No geographical data yet</span>';
        }

        document.getElementById('exportCsvBtn').href = '/api/v1/analytics/' + encodeURIComponent(code) + '/export.csv';
      } catch (err) {
        alert('Failed to fetch analytics: ' + err.message);
      }
    }
  </script>
</body>
</html>
`
