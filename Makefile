build:
	@docker build -t go-json .

start:
	@docker run -it --rm --name go-json go-json

sample:
	@go run . /sample ./sample.json
