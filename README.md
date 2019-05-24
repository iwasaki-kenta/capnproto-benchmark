# Cap'n Proto RPC Benchmark

Quick benchmark utility against a local/remote server for Cap'n Proto. The testing setup is against a DigitalOcean droplet with 4vCPUs and 8GB of RAM in Frankfurt with a ping latency of ~220ms.

It appears on a single-core case, only about 3 or 4 RPCs may be made per second. A concurrent case with 16 workers allows for 64 RPCs at most to be made per second. The average is around ~50 RPCs per second.

Each request/response individually holds a randomly-generated 600 byte payload. It appears under this case, gRPC with its HTTP/2 framing protocol is a clear winner in the case of networking over high-latency environments.

A simple Makefile is provided to build and upload either the single/concurrent case to some remote server via `rsync`, with the endpoint of the server being specified as an environment variable `ENDPOINT`.

Logs of number of received responses and sent requests are periodically emitted to stdout every second.