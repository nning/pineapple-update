.PHONY: clean run

pineapple-update: main.go
	CGO_ENABLED=0 go build .

run: pineapple-update
	./pineapple-update

clean:
	rm -f pineapple-update *.AppImage