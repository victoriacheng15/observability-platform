# Self-Hosted Observability Hub

A personal telemetry system that collects **system metrics** and **application events**, stores them in **PostgreSQL** (supercharged with **TimescaleDB** and **PostGIS**), and visualizes everything in **Grafana**, enabling high-performance time-series analysis and flexible data correlation.

---

## 🔍 What It Does

| Capability | Details |
| :--- | :--- |
| **System Metrics** | Collects CPU, memory, disk, and network stats via a lightweight Go collector (`gopsutil`). |
| **App Events** | Tracks application events from [personal reading analytics dashboard](https://github.com/victoriacheng15/personal-reading-analytics-dashboard). |
| **Storage** | Stores data in **TimescaleDB (PostgreSQL)** using a flexible JSONB schema for cross-type querying. |
| **Visualization** | **Grafana** dashboards for **Infrastructure Health** (metrics) and **Application Telemetry** (trends). |
| **Reliability** | Ensures durability via automated volume backups and idempotent collectors. |

---

## 🏗 Architecture Overview

The system is designed for simplicity and reliability:

| Component | Approach |
| :--- | :--- |
| **Collectors** | Custom Go services—cron-driven, stateless, and idempotent. |
| **Storage** | **TimescaleDB** with JSONB for unified metric/event storage. |
| **Visualization** | **Grafana** dashboards separated by concern (infra vs. app). |
| **Observability** | **Loki** (logs) and a **Go proxy** (ETL) for full telemetry coverage. |

For full architecture details, data model, and Mermaid diagrams:  
→ [docs/architecture/README.md](./docs/architecture/README.md)

---

## 💡 Why I Built This

I built this to understand how observability can reveal system behavior over time—not just through logs or one-off checks, but through continuous, contextual telemetry.

It started as a personal experiment.  
Now, it’s how I monitor host health, debug pipeline failures, and connect infrastructure metrics to application events—all without cloud dependencies.
