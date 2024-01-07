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

install: release
	install -Dm755 $(BIN) ~/.local/bin/$(BIN)
	cp -n pineapple-update.example.yml ~/.config/pineapple-update.yml || true
	mkdir -p ~/.config/systemd/user
	cp -n systemd/pineapple-update.service ~/.config/systemd/user/pineapple-update.service || true
	cp -n systemd/pineapple-update.timer ~/.config/systemd/user/pineapple-update.timer || true
