
.PHONY: clean
clean:
	rm -r dist/*

.PHONY: build
build: clean
	go build -ldflags="-H windowsgui -s -w" -o dist/powerplan2go.exe *.go

.PHONY: run
run:
	go run *.go