#!/bin/bash
# Friendly print /proc/loadavg values

function printMetric {
    SYSTEM=$1
    SUBSYSTEM=$2
    METRIC_NAME=$3
    METRIC_VALUE=$4
    echo "$SYSTEM $SUBSYSTEM $METRIC_NAME $METRIC_VALUE"
}

while read -r load1min load5min load15min jobs lastpid; do
    printMetric "demo01," "instance_01," "loadscript_load1," "$load1min"
    printMetric "demo01," "instance_01," "loadscript_load5," "$load5min"
    printMetric "demo01," "instance_01," "loadscript_load15," "$load15min"
    while IFS='/' read -r running background; do
        printMetric "demo01," "instance_01," "loadscript_jobs_running," "$running"
        printMetric "demo01," "instance_01," "loadscript_jobs_background," "$background"
    done <<< "$jobs"
    printMetric "demo01," "instance_01," "loadscript_pid_last," "$lastpid"
done < /proc/loadavg


