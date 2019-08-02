OFFLINE=""

install: 
	$(MAKE) -C crowd-funding install
	./scripts/testdependencies.sh


start:
	./scripts/testdependencies.sh
	./scripts/start.sh

demo:
	./scripts/crowd-funding-demo.sh crowdfundnet 8080 $(OFFLINE)

stop:
	./scripts/stop.sh

.PHONY: start demo stop