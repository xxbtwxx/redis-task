.PHONY: run
run:
	docker compose up

.PHONY: test
test:
	docker compose up redis-server -d
	@go test ./... -cover -race -count=1 && touch _testok; \
	docker compose down redis-server; \
	if [ -f "_testok" ]; then \
		rm -f _testok; \
		exit 0; \
	else \
		rm -f _testok; \
		exit 1; \
	fi
	