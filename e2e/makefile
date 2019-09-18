OFFLINE=""
NOSTOP=""
VERBOSE=""
TEST=crowdfundnet

install: 
	-$(MAKE) -C crowd-funding install
	-$(MAKE) -C tests/jointest install
	-$(MAKE) -C tools install

	./scripts/testdependencies.sh

start:
	./scripts/testdependencies.sh
	./scripts/start.sh --network=$(TEST) --init-ip=172.77.5.10 --heartbeat="20ms"

demo:
	./scripts/crowd-funding-demo.sh $(TEST)

run: 
	./tests/$(TEST)/run-test.sh	

stop:
	./scripts/stop.sh $(TEST)

tests: 
	./scripts/runtests.sh

test:
	./scripts/runtests.sh $(TEST) $(NOSTOP)		

.PHONY: start demo stop run tests test install
