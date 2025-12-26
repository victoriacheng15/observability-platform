#!/bin/bash
cd /home/server/software/observability-platform/system-metrics

./metrics-collector >> /var/log/system-metrics.log 2>&1