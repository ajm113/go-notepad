.PHONY: build lint

all: build

build:
	go build -o notepad -tags pango_1_42,gtk_3_22 .

lint:
	golangci-lint run --fast