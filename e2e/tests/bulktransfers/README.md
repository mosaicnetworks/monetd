

```bash
$ make installgiv; rm -rf ~/.giverny/networks/benchmark; giverny network new benchmark   --names sampledata/docker.txt  --initial-ip 172.77.5.10 -v

$ giverny network start benchmark --start-nodes -v

giverny transactions generate -n benchmark --count 100

```