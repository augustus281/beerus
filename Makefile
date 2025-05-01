generate:
	cd proto && buf generate
update:
	cd proto && buf mod update
run:
	go run cmd/beerus/main.go
