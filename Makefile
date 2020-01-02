build:
	go build

run: build
	./main

clean:
	@rm main