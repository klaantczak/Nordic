#!/bin/bash

# Usage:
#     ./run.all.sh 25

N=$1

./run.sh nordic32 baseline $N
./run.sh nordic32d1 baseline $N
./run.sh nordic32d1 breaker.attacker-weekly.one-channel.no-inspections.all $N
./run.sh nordic32d1 breaker.attacker-weekly.one-channel.no-inspections.lines $N
./run.sh nordic32d1 breaker.attacker-weekly.one-channel.no-inspections.loads $N
./run.sh nordic32d1 breaker.attacker-weekly.one-channel.no-inspections.generators $N
./run.sh nordic32d1 breaker.attacker-weekly.one-channel.yearly-inspections.all $N
./run.sh nordic32d1 breaker.attacker-weekly.one-channel.yearly-inspections.lines $N
./run.sh nordic32d1 breaker.attacker-weekly.one-channel.yearly-inspections.loads $N
./run.sh nordic32d1 breaker.attacker-weekly.one-channel.yearly-inspections.generators $N
./run.sh nordic32d1 breaker.attacker-weekly.one-channel.monthly-inspections.all $N
./run.sh nordic32d1 breaker.attacker-weekly.one-channel.monthly-inspections.lines $N
./run.sh nordic32d1 breaker.attacker-weekly.one-channel.monthly-inspections.loads $N
./run.sh nordic32d1 breaker.attacker-weekly.one-channel.monthly-inspections.generators $N
./run.sh nordic32d2 baseline $N
./run.sh nordic32d2 breaker.attacker-weekly.one-channel.no-inspections.all $N
./run.sh nordic32d2 breaker.attacker-weekly.one-channel.no-inspections.lines $N
./run.sh nordic32d2 breaker.attacker-weekly.one-channel.no-inspections.loads $N
./run.sh nordic32d2 breaker.attacker-weekly.one-channel.no-inspections.generators $N
./run.sh nordic32d2 breaker.attacker-weekly.one-channel.yearly-inspections.all $N
./run.sh nordic32d2 breaker.attacker-weekly.one-channel.yearly-inspections.lines $N
./run.sh nordic32d2 breaker.attacker-weekly.one-channel.yearly-inspections.loads $N
./run.sh nordic32d2 breaker.attacker-weekly.one-channel.yearly-inspections.generators $N
./run.sh nordic32d2 breaker.attacker-weekly.one-channel.monthly-inspections.all $N
./run.sh nordic32d2 breaker.attacker-weekly.one-channel.monthly-inspections.lines $N
./run.sh nordic32d2 breaker.attacker-weekly.one-channel.monthly-inspections.loads $N
./run.sh nordic32d2 breaker.attacker-weekly.one-channel.monthly-inspections.generators $N
./run.sh nordic32d2 breaker.attacker-weekly.two-channel.no-inspections.all $N
./run.sh nordic32d2 breaker.attacker-weekly.two-channel.no-inspections.lines $N
./run.sh nordic32d2 breaker.attacker-weekly.two-channel.no-inspections.loads $N
./run.sh nordic32d2 breaker.attacker-weekly.two-channel.no-inspections.generators $N
./run.sh nordic32d2 breaker.attacker-weekly.two-channel.yearly-inspections.all $N
./run.sh nordic32d2 breaker.attacker-weekly.two-channel.yearly-inspections.lines $N
./run.sh nordic32d2 breaker.attacker-weekly.two-channel.yearly-inspections.loads $N
./run.sh nordic32d2 breaker.attacker-weekly.two-channel.yearly-inspections.generators $N
./run.sh nordic32d2 breaker.attacker-weekly.two-channel.monthly-inspections.all $N
./run.sh nordic32d2 breaker.attacker-weekly.two-channel.monthly-inspections.lines $N
./run.sh nordic32d2 breaker.attacker-weekly.two-channel.monthly-inspections.loads $N
./run.sh nordic32d2 breaker.attacker-weekly.two-channel.monthly-inspections.generators $N
