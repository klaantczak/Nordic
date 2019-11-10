CDF Analysis
===

Cumulative Distribution Function (CDF) for the Nordic32 simulation results is
defined as f(x) := | { l | l from L and l < x } | / | L | where | | denotes
cardinality of set, L is a set of total loads delivered for each simulation.
The function f is defined on the set [0; Lm] where Lm is the maximum possible
delivered load and has values from [0; 1]. The value of the function represents
probability for the randomly selected simulation being less than the x.

The CDF is calculated for baseline and different attack frequency of single and
multiple attackers.

The model is nordic32 model with corrected reward function and max for links
on 120% from normal flow.

Tools:

    tools/
        jobs.all.sh     - run the simulations as HPC jobs
        run.all.sh      - simulate all networks
        load.cdf.go     - generate CDF raw data from simulation results
        load.avg.go     - calculate average over the simulation results
        plot.go         - draw CDF raw data
