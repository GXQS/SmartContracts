.PHONY: test fuzz bench security ci

test:
go test ./...

fuzz:
go test ./tests -run Fuzz -fuzz=Fuzz -fuzztime=10s

bench:
go test ./benchmarks -bench=. -run ^$

security:
go vet ./...

ci: test security bench
