#!/bin/sh
#PBS -N n32.analysis.cdf.tools.job
#PBS -d /home/acmp148/logs
#PBS -l walltime=2:00:00
#PBS -l nodes=1:ppn=1
#PBS -m ae
#PBS -M acmp148@city.ac.uk

# Usage:
#   qsub -v EXP=baseline,N=200 job.sh
#   qsub -v EXP=daily-attacks,N=200 job.sh

export GOPATH=/home/acmp148/hps
HPSDIR=/home/acmp148/hps
N32="$GOPATH/bin/nordic32 -time 10 -file $HPSDIR/models/nordic32.json"
JSLOG=$HPSDIR/analysis/cdf/data/$EXP.jslog

if [ -e $JSLOG ]; then rm $JSLOG; fi
for i in `seq 1 $N`
do
	echo "$EXP $i/$N"
	$N32 -network $EXP | grep message | grep "Total load" >> $JSLOG
done
