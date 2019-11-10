#!/bin/bash

# Usage:
#     ./run.sh baseline 300
#     ./run.sh daily-attacks 300

EXP=$1
N=$2

N32="../../../bin/nordic32 -time 10 -file ../../../models/nordic32.json"
JSLOG=../data/$EXP.jslog

if [ -e $JSLOG ]; then rm $JSLOG; fi
for i in `seq 1 $N`
do
    echo "$EXP $i/$N"
    $N32 -network $EXP | grep message | grep "Total load" >> $JSLOG
done
