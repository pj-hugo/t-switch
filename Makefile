BINARY_NAME=t-switch
INSTALL_PATH=/usr/local/bin
CONFIG_DIR=~/.config/t-switch

all: build

build:
	go build -o $(BINARY_NAME) .

install: build
	sudo cp $(BINARY_NAME) $(INSTALL_PATH)/$(BINARY_NAME)
	sudo chmod +x $(INSTALL_PATH)/$(BINARY_NAME)

uninstall:
	sudo rm -f $(INSTALL_PATH)/$(BINARY_NAME)
	@echo "Configuration files in $(CONFIG_DIR) were not removed."

clean:
	rm -f $(BINARY_NAME)

.PHONY: build install uninstall clean