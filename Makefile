build: 
	go build ./main.go

test: build
	./test.sh

clean:
	rm -f main tmp.s tmp

.PHONY: test clean