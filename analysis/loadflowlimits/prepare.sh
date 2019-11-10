#!/bin/bash

EXP=$1
N=$2

N32="../../bin/nordic32 -events 1 -network baseline -file ../../models/nordic32.json"
JSLOG=result.jslog

if [ -e $JSLOG ]; then rm $JSLOG; fi
$N32 | grep message | grep link | grep flow >> $JSLOG
