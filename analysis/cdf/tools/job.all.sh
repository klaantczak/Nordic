#!/bin/sh

# Usage:
#     ./job.all.sh 300

N=$1

qsub -v EXP=baseline,N=$N job.sh
qsub -v EXP=yearly-attacks,N=$N job.sh
qsub -v EXP=monthly-attacks,N=$N job.sh
qsub -v EXP=weekly-attacks,N=$N job.sh
qsub -v EXP=daily-attacks,N=$N job.sh
qsub -v EXP=modificator.yearly-attacks,N=$N job.sh
qsub -v EXP=modificator.monthly-attacks,N=$N job.sh
qsub -v EXP=modificator.weekly-attacks,N=$N job.sh
qsub -v EXP=modificator.daily-attacks,N=$N job.sh
qsub -v EXP=seven-weekly-attackers,N=$N job.sh
qsub -v EXP=modificator.weekly-attacks.yearly-inspections,N=$N job.sh
qsub -v EXP=modificator.weekly-attacks.monthly-inspections,N=$N job.sh
qsub -v EXP=modificator.weekly-attacks.weekly-inspections,N=$N job.sh
qsub -v EXP=modificator.weekly-attacks.daily-inspections,N=$N job.sh
