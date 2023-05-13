.PHONY: clean run

BIN = pineapple-update
GOFLAGS += -ldflags "$(GOLDFLAGS)"

$(BIN): main.go
	CGO_ENABLED=0 go build $(GOFLAGS)

run: $(BIN)
	./$(BIN)

clean:
	rm -f $(BIN) *.AppImage

release: GOLDFLAGS += -s -w
release: GOFLAGS += -trimpath -buildmode=pie -mod=readonly -modcacherw
release: $(BIN)