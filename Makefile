APP_NAME ?= torrent-411
ENTRYPOINT ?= ./inspect_torrent.go
DIST_DIR ?= dist
TARGETS := $(shell go tool dist list)
README_FILE ?= README.md
REQUIRED_WINDOWS_TARGETS := windows/386 windows/amd64
BUILT_TARGETS_FILE := $(DIST_DIR)/built-targets.txt
SKIPPED_TARGETS_FILE := $(DIST_DIR)/skipped-targets.txt

.PHONY: all clean build-all zip-all package-all list-targets validate-targets

all: package-all

package-all: clean validate-targets build-all zip-all

validate-targets:
	@set -eu; \
	for required in $(REQUIRED_WINDOWS_TARGETS); do \
		case " $(TARGETS) " in \
			*" $$required "*) ;; \
			*) \
				echo "Missing required target in TARGETS: $$required"; \
				exit 1; \
				;; \
		esac; \
	done

build-all: validate-targets
	@mkdir -p "$(DIST_DIR)/bin"
	@set -eu; \
	: > "$(BUILT_TARGETS_FILE)"; \
	: > "$(SKIPPED_TARGETS_FILE)"; \
	for target in $(TARGETS); do \
		goos=$${target%/*}; \
		goarch=$${target#*/}; \
		ext=""; \
		if [ "$$goos" = "windows" ]; then ext=".exe"; fi; \
		out="$(DIST_DIR)/bin/$(APP_NAME)-$$goos-$$goarch$$ext"; \
		echo "Building $$target -> $$out"; \
		if CGO_ENABLED=0 GOOS=$$goos GOARCH=$$goarch \
			go build -trimpath -ldflags='-s -w' -o "$$out" "$(ENTRYPOINT)"; then \
			printf '%s\n' "$$target" >> "$(BUILT_TARGETS_FILE)"; \
			continue; \
		fi; \
		echo "WARN: skipping unsupported target $$target (build failed in this environment)"; \
		printf '%s\n' "$$target" >> "$(SKIPPED_TARGETS_FILE)"; \
	done
	@set -eu; \
	for required in $(REQUIRED_WINDOWS_TARGETS); do \
		if ! grep -qx "$$required" "$(BUILT_TARGETS_FILE)"; then \
			echo "Required Windows target did not build: $$required"; \
			exit 1; \
		fi; \
	done; \
	if [ -s "$(SKIPPED_TARGETS_FILE)" ]; then \
		echo "Skipped targets (not buildable in current environment):"; \
		cat "$(SKIPPED_TARGETS_FILE)"; \
	fi

zip-all: validate-targets
	@mkdir -p "$(DIST_DIR)/zip"
	@set -eu; \
	if [ ! -f "$(README_FILE)" ]; then \
		echo "Missing README file: $(README_FILE)"; \
		exit 1; \
	fi; \
	if [ ! -s "$(BUILT_TARGETS_FILE)" ]; then \
		echo "No built targets found. Run make build-all first."; \
		exit 1; \
	fi; \
	while IFS= read -r target; do \
		goos=$${target%/*}; \
		goarch=$${target#*/}; \
		ext=""; \
		if [ "$$goos" = "windows" ]; then ext=".exe"; fi; \
		bin="$(DIST_DIR)/bin/$(APP_NAME)-$$goos-$$goarch$$ext"; \
		archive="$(DIST_DIR)/zip/$(APP_NAME)-$$goos-$$goarch.zip"; \
		if [ ! -f "$$bin" ]; then \
			echo "Missing binary: $$bin"; \
			exit 1; \
		fi; \
		rm -f "$$archive"; \
		zip -j -q "$$archive" "$$bin" "$(README_FILE)"; \
		echo "Zipped $$archive"; \
	done < "$(BUILT_TARGETS_FILE)"

list-targets:
	@go tool dist list

clean:
	@rm -rf "$(DIST_DIR)"
