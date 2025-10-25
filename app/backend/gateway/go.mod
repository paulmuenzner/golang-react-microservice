module github.com/app/gateway

go 1.23.0

toolchain go1.24.1

replace github.com/app/shared/go => ../../../shared/go

require (
	github.com/app/shared/go v0.0.0-00010101000000-000000000000
	github.com/google/uuid v1.6.0
)

require (
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/rs/zerolog v1.34.0 // indirect
	golang.org/x/sys v0.35.0 // indirect
)
