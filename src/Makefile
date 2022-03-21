PROJECT_NAME = tg-bot
TARGET_DIR = target

.PHONY: all darwin linux windows clean

all: linux

darwin:
	GOOS=darwin GOARCH=amd64 go build -o $(TARGET_DIR)/$(PROJECT_NAME)-darwin main.go

linux:
	GOOS=linux GOARCH=amd64 go build -o $(TARGET_DIR)/$(PROJECT_NAME)-linux main.go

windows:
	GOOS=windows GOARCH=amd64 go build -o $(TARGET_DIR)/$(PROJECT_NAME)-windows.exe main.go

clean:
	@echo "cleaning"
	@rm -rf $(TARGET_DIR)
	@echo "cleaned"
