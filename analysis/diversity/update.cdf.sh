#!/bin/bash

for i in `ls data/*.jslog`
do
    go run ../cdf/tools/load.cdf.go --jslog $i > ${i%.jslog}.cdf.csv
done
