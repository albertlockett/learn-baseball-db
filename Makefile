default: build

clean:
	rm learn-baseball-db

build:
	go build
	docker build . -t albertlockett2/learnbaeball-db:latest

push:
	docker push albertlockett2/learnbaeball-db:latest
