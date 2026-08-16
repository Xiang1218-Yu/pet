# Healing Pet

Healing Pet is a desktop companion application written in Go and Fyne. It models a pet's hunger, happiness, hygiene, energy, care history, and growth stage, then persists the state locally.

## Layout

- `cmd/pet`: desktop application entry point and UI wiring
- `internal/pet`: pet state model, growth rules, and interaction APIs
- `internal/appearance`: Fyne canvas rendering for pet appearances
- `internal/storage`: local JSON save/load support

## Requirements

- Go 1.26.5 or a compatible Go toolchain
- The platform dependencies required by Fyne for desktop builds

## Build and test

```bash
go build ./...
go test ./...
```

## Run

```bash
go run ./cmd/pet
```

The application stores save data in a `.healing-pet` directory under the user home directory when it is writable, with executable-directory and current-directory fallbacks.
