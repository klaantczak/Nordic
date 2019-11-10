HPS
===

HPS is a high-performance simulation engine that runs the networks of
probabilistic state machines. Such networks also known as stochastic activity
networks and can be applied while analysing behaviour and assessing savety
and security properties of the systems.

The engine is created in Go. Engine can be embedded into the application
or service or can be used as command-line tool.

The base component of the model is state machine, which resides within the
environment. The state machine can be implemented in various ways: as
Markov state machine, as hierarchical composition of the Markov state machines,
or as engine plugin.
 