#!/bin/bash

# Usage:
#     ./run.sh nordic32 baseline 25
#     ./run.sh nordic32d2 daily-attacks 25

MODEL=$1
EXP=$2
N=$3

N32="../../bin/nordic32 -time 10"
NAME=data/$MODEL-$EXP
JSLOG=$NAME.jslog

if [ -e $JSLOG ]; then rm $JSLOG; fi
for i in `seq 1 $N`
do
    echo "$EXP $i/$N"
    $N32 -file ../../models/$MODEL.json -network $EXP | grep message | grep "Total load" >> $JSLOG
done
go run ../cdf/tools/load.avg.go --jslog $JSLOG
go run ../cdf/tools/load.cdf.go --jslog $JSLOG > $NAME.cdf.csv


