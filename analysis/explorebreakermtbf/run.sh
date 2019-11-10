#!/bin/bash

# Usage:
#     ./run.sh 1 25

MTBF=$1
N=$2

N32="../../bin/nordic32 -time 10 -file ../../models/nordic32d.json -network baseline"
NAME=data/mtbf-$MTBF
JSLOG=$NAME.jslog

if [ -e $JSLOG ]; then rm $JSLOG; fi
for i in `seq 1 $N`
do
    echo "mtbf=$MTBF $i/$N"
    BREAKER_COMPONENT_MTBF=$MTBF $N32 | grep message >> $JSLOG
done
go run ../cdf/tools/load.avg.go --jslog $JSLOG
go run ../cdf/tools/load.cdf.go --jslog $JSLOG > $NAME.cdf.csv


