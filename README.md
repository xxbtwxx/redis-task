# Redis task implementation

A simple service for consuming and "processing" messages 
received from a Redis Pub/Sub

The number of consumers can be configured using the `config.yaml`
With the current setup each consumer can "process" ~3 messages per second

# Running the service

The service and all of it's dependencies can be run using `make`
This will spin 
`redis-server` on port `6379`
`redis-insight` on port `5540`
`prometheus` on port `9090`
`grafana` on port `3000` 
and the service will expose `/metrics` on port `2112` it will also expose `/debug/pprof/...` endpoints which can be used for profiling and diagnostics

In order to send messages to the service you can execute `python producer.py` 
You may need to execute `pip install redis` in order to run the python script

# Metrics

Some basic panels are already preconfigured in `grafana`
They include total number of "processed" messages, "processed" messages per consumer
and the time needed to "process" a message
The credentials are `user`/`password`

# Testing

The tests can be executed by calling `make test` 
