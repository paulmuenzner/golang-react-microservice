module github.com/app/service-b

go 1.23.0

toolchain go1.24.1

require github.com/app/shared/go v0.0.0

replace github.com/app/shared/go => ../../../shared/go
