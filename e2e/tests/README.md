# Tests

## End to End Tests

+ `crowdfundnet` --- This is the demo non-interactive
+ `jointest` --- This is a test of the poa joining


To run all tests

```bash
[...e2e] $ make tests
```


To run one test:

```bash
[...e2e] $ make test TEST=jointest
```

To run one test, but not destroy the network:

```bash
[...e2e] $ make test TEST=jointest  NOSTOP=nostop
```

Running one test is equivalent to:
```bash
[...e2e] $ make start TEST=jointest 
[...e2e] $ make run TEST=jointest
[...e2e] $ make stop TEST=jointest
```


## Creating a test

Each test need a network definition. For a network named `newnet` it would be in `.../e2e/networks/newnet.txt`.

A folder called `.../e2e/tests/newnet` will need to be created. In that folder would be a file `run-test.sh`. That file should exit with a non-zero exit code if the test fails. 




