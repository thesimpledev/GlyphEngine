# Makefile for building and deploying the game to Web, PC (Windows & Linux)

# Variables
PROJECT_PATH          = github.com/thesimpledev/GlyphEngine
DIST_DIR             = dist
VERSION              = 1.2.5

GOROOT_WASM          = $(shell go env GOROOT)/misc/wasm/wasm_exec.js

BUTLER_USER          = ApocalypseTheory
BUTLER_PROJECT       = the-social-contract

# Butler channels for each target
BUTLER_CHANNEL_WEB   = html5
BUTLER_CHANNEL_WIN   = windows
BUTLER_CHANNEL_LINUX = linux

# Phony targets
.PHONY: build build-web build-pc-windows build-pc-linux deploy-web deploy-windows deploy-linux clean depend status demo test help

help:
	@echo "Available targets:"
	@echo "  make build              - Clean, deps, build Web/Win/Linux, deploy Win/Linux to itch.io, then clean again"
	@echo "  make build-web          - Build for Web (HTML5)"
	@echo "  make build-pc-windows   - Build for Windows"
	@echo "  make build-pc-linux     - Build for Linux"
	@echo "  make clean              - Remove all build artifacts"
	@echo "  make depend             - Update dependencies"
	@echo "  make demo               - Run the Linux build locally (if it exists)"
	@echo "  make test               - Run tests"
	@echo "  make status             - Check deployment status on itch.io"

# -------------------------------------------------------
# Single Build Target
# -------------------------------------------------------
build: depend build-web build-pc-windows build-pc-linux deploy-web deploy-windows deploy-linux clean
	@echo "=== All builds completed successfully ==="

# -------------------------------------------------------
# Build Web (HTML5)
# -------------------------------------------------------
build-web:
	@echo "=== Building Web (HTML5) ==="
	@mkdir -p build/web
	@cp $(GOROOT_WASM) build/web/

	@env GOOS=js GOARCH=wasm go build -o build/web/$(BUTLER_PROJECT).wasm $(PROJECT_PATH)
	@echo '<!DOCTYPE html><html><head><meta charset="UTF-8"><title>$(BUTLER_PROJECT)</title></head><body><script src="wasm_exec.js"></script><script>if(!WebAssembly.instantiateStreaming){WebAssembly.instantiateStreaming=async(resp,importObject)=>{const source=await(resp).arrayBuffer();return await WebAssembly.instantiate(source,importObject);};}const go=new Go();WebAssembly.instantiateStreaming(fetch("$(BUTLER_PROJECT).wasm"),go.importObject).then(result=>{go.run(result.instance);});</script></body></html>' > build/web/index.html
	@echo "Web (HTML5) build done."

deploy-web:
	@echo "=== Deploying Web (HTML5) Build to itch.io ==="
	@butler push build/web $(BUTLER_USER)/$(BUTLER_PROJECT):$(BUTLER_CHANNEL_WEB) --userversion $(VERSION)
	@echo "Web build deployed."

# -------------------------------------------------------
# Build Windows
# -------------------------------------------------------
build-pc-windows:
	@echo "=== Building for Windows ==="
	@mkdir -p $(DIST_DIR)/pc/windows
	@env GOOS=windows GOARCH=amd64 go build -o $(DIST_DIR)/pc/windows/$(BUTLER_PROJECT).exe $(PROJECT_PATH)
	@echo "Windows build done."

deploy-windows:
	@echo "=== Deploying Windows Build to itch.io ==="
	@butler push $(DIST_DIR)/pc/windows $(BUTLER_USER)/$(BUTLER_PROJECT):$(BUTLER_CHANNEL_WIN) --userversion $(VERSION)
	@echo "Windows build deployed."

# -------------------------------------------------------
# Build Linux
# -------------------------------------------------------
build-pc-linux:
	@echo "=== Building for Linux ==="
	@mkdir -p $(DIST_DIR)/pc/linux
	@env GOOS=linux GOARCH=amd64 go build -o $(DIST_DIR)/pc/linux/$(BUTLER_PROJECT) $(PROJECT_PATH)
	@echo "Linux build done."

deploy-linux:
	@echo "=== Deploying Linux Build to itch.io ==="
	@butler push $(DIST_DIR)/pc/linux $(BUTLER_USER)/$(BUTLER_PROJECT):$(BUTLER_CHANNEL_LINUX) --userversion $(VERSION)
	@echo "Linux build deployed."

# -------------------------------------------------------
# Utility Targets
# -------------------------------------------------------
clean:
	@echo "=== Cleaning up build artifacts ==="
	@rm -rf build
	@rm -rf $(DIST_DIR)
	@echo "All cleaned."

depend:
	@echo "=== Updating Dependencies ==="
	@go mod tidy
	@echo "Dependencies updated."

status:
	@echo "=== Checking Deployment Status on itch.io ==="
	@butler status $(BUTLER_USER)/$(BUTLER_PROJECT):$(BUTLER_CHANNEL_WEB)
	@butler status $(BUTLER_USER)/$(BUTLER_PROJECT):$(BUTLER_CHANNEL_WIN)
	@butler status $(BUTLER_USER)/$(BUTLER_PROJECT):$(BUTLER_CHANNEL_LINUX)

demo:
	@echo "=== Running the Game Locally (Linux Build) ==="
	@if [ -f $(DIST_DIR)/pc/linux/$(BUTLER_PROJECT) ]; then \
	  ./$(DIST_DIR)/pc/linux/$(BUTLER_PROJECT); \
	else \
	  echo 'No Linux build found at $(DIST_DIR)/pc/linux/$(BUTLER_PROJECT)'; \
	fi

test:
	@echo "=== Running Tests ==="
	@go test ./...
