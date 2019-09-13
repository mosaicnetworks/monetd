OFFLINE=""
NOSTOP=""
VERBOSE=""
TEST=crowdfundnet

install: 
	-$(MAKE) -C crowd-funding install
	-$(MAKE) -C tests/jointest install
	-$(MAKE) -C tests/transfers install
	-$(MAKE) -C tests/bulktransfers install

	./scripts/testdependencies.sh

start:
	./scripts/testdependencies.sh
	./scripts/start.sh $(TEST) 172.77.5.10 $(VERBOSE)

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
