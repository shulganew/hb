# hb service

## Generate oapi
Use make or bash command build generate
```
make oapi
```
```bash
go run github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen --config=internal/config/oapi.yaml --package oapi api/api.yaml
```
