# cmanage

## Developpers

- Installing golang : https://golang.org/doc/install
- Running application : 

```
go run main.go
```

Cmanage uses go-bindata https://github.com/jteeuwen/go-bindata to bundle static files (docker-compose.yml ...)

Rebinding data folder : 
go-bindata data/

Change the namespace of the generated file (bindata.go) to "package bin" and replace the old one in the bin folder

### Compiling
```
go complile
go install
```
