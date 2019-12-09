#!/bin/bash

go run load.cdf.go --jslog ../data/baseline.jslog > ../data/baseline.cdf.csv
go run load.cdf.go --jslog ../data/yearly-vpn-attacks.jslog > ../data/yearly-vpn-attacks.cdf.csv
go run load.cdf.go --jslog ../data/monthly-vpn-attacks.jslog > ../data/monthly-vpn-attacks.cdf.csv
go run load.cdf.go --jslog ../data/weekly-vpn-attacks.jslog > ../data/weekly-vpn-attacks.cdf.csv
go run load.cdf.go --jslog ../data/daily-vpn-attacks.jslog > ../data/daily-vpn-attacks.cdf.csv
go run load.cdf.go --jslog ../data/seven-weekly-attackers.jslog > ../data/seven-weekly-attackers.cdf.csv
go run load.cdf.go --jslog ../data/modificator.yearly-attacks.jslog > ../data/modificator.yearly-attacks.cdf.csv
go run load.cdf.go --jslog ../data/modificator.monthly-attacks.jslog > ../data/modificator.monthly-attacks.cdf.csv
go run load.cdf.go --jslog ../data/modificator.weekly-attacks.jslog > ../data/modificator.weekly-attacks.cdf.csv
go run load.cdf.go --jslog ../data/modificator.daily-attacks.jslog > ../data/modificator.daily-attacks.cdf.csv
go run load.cdf.go --jslog ../data/modificator.weekly-attacks.yearly-inspections.jslog > ../data/modificator.weekly-attacks.yearly-inspections.cdf.csv
go run load.cdf.go --jslog ../data/modificator.weekly-attacks.monthly-inspections.jslog > ../data/modificator.weekly-attacks.monthly-inspections.cdf.csv
go run load.cdf.go --jslog ../data/modificator.weekly-attacks.weekly-inspections.jslog > ../data/modificator.weekly-attacks.weekly-inspections.cdf.csv
go run load.cdf.go --jslog ../data/modificator.weekly-attacks.daily-inspections.jslog > ../data/modificator.weekly-attacks.daily-inspections.cdf.csv
go run plot.go ../charts/attacks.png \
    baseline ../data/baseline.cdf.csv \
    yearly ../data/yearly-vpn-attacks.cdf.csv \
    monthly ../data/monthly-vpn-attacks.cdf.csv \
    weekly ../data/weekly-vpn-attacks.cdf.csv \
    daily ../data/daily-vpn-attacks.cdf.csv
go run plot.go ../charts/modificators.baseline.yearly.png \
    baseline ../data/baseline.cdf.csv \
    yearly ../data/modificator.yearly-attacks.cdf.csv
go run plot.go ../charts/modificators.yearly.monthly.png \
    yearly ../data/modificator.yearly-attacks.cdf.csv \
    monthly ../data/modificator.monthly-attacks.cdf.csv
go run plot.go ../charts/modificators.monthly.weekly.png \
    monthly ../data/modificator.monthly-attacks.cdf.csv \
    weekly ../data/modificator.weekly-attacks.cdf.csv
go run plot.go ../charts/modificators.weekly.daily.png \
    weekly ../data/modificator.weekly-attacks.cdf.csv \
    daily ../data/modificator.daily-attacks.cdf.csv
go run plot.go ../charts/attakers.png \
    baseline ../data/baseline.cdf.csv \
    weekly ../data/weekly-attacks.cdf.csv \
    weekly7 ../data/seven-weekly-attackers.cdf.csv \
    daily ../data/daily-attacks.cdf.csv
go run plot.go ../charts/inspections.png \
    baseline ../data/baseline.cdf.csv \
    monthly ../data/modificator.weekly-attacks.monthly-inspections.cdf.csv \
    weekly ../data/modificator.weekly-attacks.weekly-inspections.cdf.csv \
    daily ../data/modificator.weekly-attacks.daily-inspections.cdf.csv
