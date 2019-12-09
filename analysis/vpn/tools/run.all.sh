#!/bin/bash

# Usage:
#     ./run.all.sh 300

N=$1

./run.sh baseline $N
./run.sh yearly-attacks $N
./run.sh monthly-attacks $N
./run.sh weekly-attacks $N
./run.sh daily-attacks $N
./run.sh modificator.yearly-attacks $N
./run.sh modificator.monthly-attacks $N
./run.sh modificator.weekly-attacks $N
./run.sh modificator.daily-attacks $N
./run.sh seven-weekly-attackers $N
./run.sh modificator.weekly-attacks.daily-inspections $N
./run.sh modificator.weekly-attacks.weekly-inspections $N
./run.sh modificator.weekly-attacks.monthly-inspections $N
./run.sh modificator.weekly-attacks.yearly-inspections $N
