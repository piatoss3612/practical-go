build:
	@echo build docker image for test
	docker build -t dining-philosophers-test -f Dockerfile .

test:
	@echo run test
	docker run dining-philosophers-test