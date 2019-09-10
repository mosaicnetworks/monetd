

The benchmark sub folder here contains a giverny config. Copy it into place:

```bash 
    cp -p -r benchmark ~/.giverny/networks/
```

Start the network:

```bash
    giverny network start benchmark --start-nodes
```    

Install node things:

```bash
    cd e2e/tests/bulktransfers
    npm install
```

Watch the logs:

```bash
    docker logs Jon -f
    docker logs Martin -f
    docker logs Kevin -f
    docker logs Danu -f
```

Launch the tests

```bash
    ./run-test.sh
```    





Stop the network:

```bash
   giverny network stop benchmark
```

Tidy up:
```bash
   docker container prune    
```



Ignore this for now:

```bash
$ make installgiv; rm -rf ~/.giverny/networks/benchmark; giverny network new benchmark   --names sampledata/docker.txt  --initial-ip 172.77.5.10 -v

$ giverny network start benchmark --start-nodes -v

giverny transactions generate -n benchmark --count 100

```