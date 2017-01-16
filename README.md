# SPUD

SPUD is a user-space CCN stack written in Go.

# Athena configuration

    - (set up the key pair)
    - athenactl add link tcp://localhost:9696/listener/name=prod
    - athenactl add route "tcp://localhost:9696<->localhost:49236" ccnx:/
