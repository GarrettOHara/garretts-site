# 🧠 garretts-site — a personal portfolio with a side of analytics

Welcome to the nerdy underbelly of my personal website — not just a digital resume, but a full-stack, self-hosted, Dockerized analytics playground.

---

## 🌐 What is this?

This project powers [my personal website](https://garrettohara.live), a portfolio built to show off not just *what* I do, but *how* I do it. Beyond a clean homepage, it also features a live **analytics dashboard** that captures and visualizes site traffic in real-time using:

- 📊 Beautiful charts (thanks to Chart.js)
- 🗺️ Interactive maps (powered by Leaflet.js)
- 📈 Machine learning-based insights (numpy, pandas, scikit-learn, etc.)

It’s a glorified digital business card. But it *thinks* it’s a data scientist.

---

## 🧰 Tech Stack

| Component        | Stack / Tooling |
|------------------|-----------------|
| **Backend**      | Go web server (`net/http`) |
| **Frontend**     | Plain ol’ HTML + CSS |
| **Charts & Maps**| [Chart.js](https://www.chartjs.org/), [Leaflet.js](https://leafletjs.com/) |
| **Database**     | SQLite3 |
| **Deployment**   | Docker Compose on a ThinkPad P52s (yes, seriously) |
| **Access Layer** | Cloudflare Tunnel → self-hosted NGINX proxy container → Go web server |

---

## 🐳 Deployment Architecture

```text
┌────────────┐
│ Internet   │
└─────┬──────┘
      ▼
┌──────────────────────┐
│ Cloudflare Tunnel    │
└────────┬─────────────┘
         ▼
┌──────────────────────┐
│ NGINX Proxy Manager  │
│ (Docker Container)   │
└────────┬─────────────┘
         ▼
┌────────────────────── ┐
│ Go Web Server         │
│ (Dockerized, serves   │
│ portfolio + analytics)│
└────────────────────── ┘
```
## 📉 Analytics Engine

Every time someone visits, the server logs key request data to an SQLite database. That raw data is used to:

- Show traffic over time  
- Cluster visitor patterns  
- Detect anomalies  
- Visualize geolocation and device types  

---

## 🔁 Powered by Python + cron

There's a second Docker container running Python scripts on a schedule via `crontab`. These scripts load traffic data and run deeper ML-powered analysis.

### 🐍 Python Stack

# 📈 Analytics Engine Overview

This analytics pipeline powers the backend data layer of my personal website. It's implemented in Python and runs in a container scheduled via `cron`. It processes traffic logs from SQLite and outputs structured JSON used by the frontend dashboard.

---

## 🔧 Stack & Architecture

- **Database**: SQLite (`requests.db`)
- **Processing**: Python + Pandas
- **Machine Learning**: scikit-learn
- **Deployment**: Docker container w/ crontab scheduler
- **Output**: JSON files served to frontend by Go webserver

---

## 🔍 Techniques Used

### 1. **KMeans Clustering**
- Groups traffic patterns based on hour of access, day of week, and device type (mobile/desktop).
- Helps identify distinct visitor behaviors.
- Output: `clusters.json`

### 2. **Time Series Aggregation**
- Resamples visits hourly and computes rolling 3-hour averages.
- Useful for identifying trends and peak activity periods.
- Output: `time_series.json`

### 3. **Anomaly Detection**
- Uses `IsolationForest` to flag statistically rare access patterns.
- Helps surface bots or unusual user behavior.
- Output: `anomalies.json`

---

## 🌍 IP Geolocation & Visualization

A second script enriches IP data with geolocation via the `ipinfo.io` API. It then generates multiple visualizations:

- **Country Distribution** (`country_distribution.json`)
- **Device Type by Country** (`device_by_country.json`)
- **Visitor Map** (`visitor_map.json`)

Used libraries include:
- `requests` for external API calls
- `pandas` for data joins and filtering
- `json` for writing frontend-compatible artifacts

These scripts generate `.json` data outputs that get served right back to the frontend. For geo insights, we query [ipinfo.io](https://ipinfo.io) to turn IPs into latitude/longitude data.

---

## 🧠 Summary

This system turns raw web traffic into actionable, visual insights with minimal infrastructure. It’s lightweight, containerized, and built entirely with open-source tools — ideal for self-hosting and personal analytics.

---

## 📦 How to Run It

### Pre-reqs:

- Docker  
- Docker Compose

```bash
git clone https://github.com/yourusername/garretts-site.git
cd garretts-site
docker compose up -d --build
```

The site will be up at `http://localhost:8080`. If you’re doing this my way, you’ll be exposing it through a Cloudflare Tunnel from a Debian box in your basement.

---

## 🚧 Under Construction

- [ ] Auth for admin-only views  
- [ ] More ML models (recommendations? anomaly detection v2?)  
- [ ] Archiving old logs  
- [ ] Light/dark theme toggle

---

## 👋 About Me

I'm Garrett, a backend software engineer with extensive devops experience, and a passion for self-hosted infrastructure. This project is equal parts resume, playground, and data lab.