#!/bin/bash

if [ -e data/report.csv ]; then rm data/report.csv; fi

go run averages.go -label 1 -jslog data/mtbf-1.jslog >> data/report.csv
go run averages.go -label 0.9 -jslog data/mtbf-0.9.jslog >> data/report.csv
go run averages.go -label 0.8 -jslog data/mtbf-0.8.jslog >> data/report.csv
go run averages.go -label 0.7 -jslog data/mtbf-0.7.jslog >> data/report.csv
go run averages.go -label 0.6 -jslog data/mtbf-0.6.jslog >> data/report.csv
go run averages.go -label 0.5 -jslog data/mtbf-0.5.jslog >> data/report.csv
go run averages.go -label 0.4 -jslog data/mtbf-0.4.jslog >> data/report.csv
go run averages.go -label 0.3 -jslog data/mtbf-0.3.jslog >> data/report.csv
go run averages.go -label 0.2 -jslog data/mtbf-0.2.jslog >> data/report.csv
go run averages.go -label 0.1 -jslog data/mtbf-0.1.jslog >> data/report.csv
go run averages.go -label 0.09 -jslog data/mtbf-0.09.jslog >> data/report.csv
go run averages.go -label 0.08 -jslog data/mtbf-0.08.jslog >> data/report.csv
go run averages.go -label 0.07 -jslog data/mtbf-0.07.jslog >> data/report.csv
go run averages.go -label 0.06 -jslog data/mtbf-0.06.jslog >> data/report.csv
go run averages.go -label 0.05 -jslog data/mtbf-0.05.jslog >> data/report.csv
go run averages.go -label 0.04 -jslog data/mtbf-0.04.jslog >> data/report.csv
go run averages.go -label 0.03 -jslog data/mtbf-0.03.jslog >> data/report.csv
go run averages.go -label 0.02 -jslog data/mtbf-0.02.jslog >> data/report.csv
go run averages.go -label 0.01 -jslog data/mtbf-0.01.jslog >> data/report.csv
go run averages.go -label 0.009 -jslog data/mtbf-0.009.jslog >> data/report.csv
go run averages.go -label 0.008 -jslog data/mtbf-0.008.jslog >> data/report.csv
go run averages.go -label 0.007 -jslog data/mtbf-0.007.jslog >> data/report.csv
go run averages.go -label 0.006 -jslog data/mtbf-0.006.jslog >> data/report.csv
go run averages.go -label 0.005 -jslog data/mtbf-0.005.jslog >> data/report.csv
go run averages.go -label 0.004 -jslog data/mtbf-0.004.jslog >> data/report.csv
go run averages.go -label 0.003 -jslog data/mtbf-0.003.jslog >> data/report.csv
