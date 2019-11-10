#!/bin/bash

export GOPATH=$PWD/../..

go run ../cdf/tools/plot.go charts/baselines.png \
    baseline data/nordic32-baseline.cdf.csv \
    baseline_1c data/nordic32d1-baseline.cdf.csv \
    baseline_2c data/nordic32d2-baseline.cdf.csv

go run ../cdf/tools/plot.go charts/attacks_1c.generators.png \
    baseline data/nordic32d1-baseline.cdf.csv \
    no_insp data/nordic32d1-breaker.attacker-weekly.one-channel.no-inspections.generators.cdf.csv \
    yearly data/nordic32d1-breaker.attacker-weekly.one-channel.yearly-inspections.generators.cdf.csv \
    monthly data/nordic32d1-breaker.attacker-weekly.one-channel.monthly-inspections.generators.cdf.csv

go run ../cdf/tools/plot.go charts/attacks_1c.loads.png \
    baseline data/nordic32d1-baseline.cdf.csv \
    no_insp data/nordic32d1-breaker.attacker-weekly.one-channel.no-inspections.loads.cdf.csv \
    yearly data/nordic32d1-breaker.attacker-weekly.one-channel.yearly-inspections.loads.cdf.csv \
    monthly data/nordic32d1-breaker.attacker-weekly.one-channel.monthly-inspections.loads.cdf.csv

go run ../cdf/tools/plot.go charts/attacks_1c.lines.png \
    baseline data/nordic32d1-baseline.cdf.csv \
    no_insp data/nordic32d1-breaker.attacker-weekly.one-channel.no-inspections.lines.cdf.csv \
    yearly data/nordic32d1-breaker.attacker-weekly.one-channel.yearly-inspections.lines.cdf.csv \
    monthly data/nordic32d1-breaker.attacker-weekly.one-channel.monthly-inspections.lines.cdf.csv

go run ../cdf/tools/plot.go charts/attacks_1c.monthly.png \
    baseline data/nordic32d1-baseline.cdf.csv \
    loads data/nordic32d1-breaker.attacker-weekly.one-channel.monthly-inspections.loads.cdf.csv \
    generators data/nordic32d1-breaker.attacker-weekly.one-channel.monthly-inspections.generators.cdf.csv \
    lines data/nordic32d1-breaker.attacker-weekly.one-channel.monthly-inspections.lines.cdf.csv \
    all data/nordic32d1-breaker.attacker-weekly.one-channel.monthly-inspections.all.cdf.csv

go run ../cdf/tools/plot.go charts/attacks_1c.yearly.png \
    baseline data/nordic32d1-baseline.cdf.csv \
    loads data/nordic32d1-breaker.attacker-weekly.one-channel.yearly-inspections.loads.cdf.csv \
    generators data/nordic32d1-breaker.attacker-weekly.one-channel.yearly-inspections.generators.cdf.csv \
    lines data/nordic32d1-breaker.attacker-weekly.one-channel.yearly-inspections.lines.cdf.csv

go run ../cdf/tools/plot.go charts/attacks_2c.1c_2c.lines.no_monthly.png \
    baseline data/nordic32d2-baseline.cdf.csv \
    1a_no_insp data/nordic32d2-breaker.attacker-weekly.one-channel.no-inspections.lines.cdf.csv \
    1a_monthly data/nordic32d2-breaker.attacker-weekly.one-channel.monthly-inspections.lines.cdf.csv \
    2a_no_insp data/nordic32d2-breaker.attacker-weekly.two-channel.no-inspections.lines.cdf.csv \
    2a_monthly data/nordic32d2-breaker.attacker-weekly.two-channel.monthly-inspections.lines.cdf.csv 

go run ../cdf/tools/plot.go charts/attacks_2c.1c.lines.png \
    baseline data/nordic32d2-baseline.cdf.csv \
    no_insp data/nordic32d2-breaker.attacker-weekly.one-channel.no-inspections.lines.cdf.csv \
    yearly data/nordic32d2-breaker.attacker-weekly.one-channel.yearly-inspections.lines.cdf.csv \
    monthly data/nordic32d2-breaker.attacker-weekly.one-channel.monthly-inspections.lines.cdf.csv

go run ../cdf/tools/plot.go charts/attacks_2c.2c.lines.png \
    baseline data/nordic32d2-baseline.cdf.csv \
    no_insp data/nordic32d2-breaker.attacker-weekly.two-channel.no-inspections.lines.cdf.csv \
    yearly data/nordic32d2-breaker.attacker-weekly.two-channel.yearly-inspections.lines.cdf.csv \
    monthly data/nordic32d2-breaker.attacker-weekly.two-channel.monthly-inspections.lines.cdf.csv

go run ../cdf/tools/plot.go charts/attacks_2c.2c.no.png \
    baseline data/nordic32d2-baseline.cdf.csv \
    loads data/nordic32d2-breaker.attacker-weekly.two-channel.no-inspections.loads.cdf.csv \
    generators data/nordic32d2-breaker.attacker-weekly.two-channel.no-inspections.generators.cdf.csv \
    lines data/nordic32d2-breaker.attacker-weekly.two-channel.no-inspections.lines.cdf.csv \
    all data/nordic32d2-breaker.attacker-weekly.two-channel.no-inspections.all.cdf.csv

go run ../cdf/tools/plot.go charts/attacks_2c_1a.lines.png \
    baseline data/nordic32d2-baseline.cdf.csv \
    no_insp data/nordic32d2-breaker.attacker-weekly.one-channel.no-inspections.lines.cdf.csv \
    yearly data/nordic32d2-breaker.attacker-weekly.one-channel.yearly-inspections.lines.cdf.csv \
    monthly data/nordic32d2-breaker.attacker-weekly.one-channel.monthly-inspections.lines.cdf.csv

go run ../cdf/tools/plot.go charts/attacks_2c_2a.lines.png \
    baseline data/nordic32d2-baseline.cdf.csv \
    no_insp data/nordic32d2-breaker.attacker-weekly.two-channel.no-inspections.lines.cdf.csv \
    yearly data/nordic32d2-breaker.attacker-weekly.two-channel.yearly-inspections.lines.cdf.csv \
    monthly data/nordic32d2-breaker.attacker-weekly.two-channel.monthly-inspections.lines.cdf.csv

