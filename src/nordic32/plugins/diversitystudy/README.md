Diversity Study
===

* Breaker has two diverse components.
* Each component has four states: Ok, Compromised, Fail.
* Ok breaker component fails once in five years.
* Compromied breaker component fails once a month.
* Failed breaker component disconnects the line and recovers in one day.
* Failed breaker component does not stop the line from being connected.
* Inspection restores the breaker component from compromised to ok.
* Recovery restores the breaker component from fail to ok.
* Breaker attacker is a simple event generator.
* Breaker inspector is a simple event generator.
* Study scenarious:
  - baseline (no attacks)
  - attacker attacks one breaker component once in a year
  - attacker attacks two breaker components once in a year

