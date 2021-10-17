build: 
	go build -o main .

test: build
	./test.sh

clean:
	rm -f main tmp.s tmp

.PHONY: test clean