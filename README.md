# Go Webservices

##### A rest webservice in Go that allows to operate over text files as resources.
- Create a text file with some contents stored in a given path.
- Retrieve the contents of a text file under the given path.
- Replace the contents of a text file.
- Delete the resource that is stored under a given path.

##### It also allows to get some statistics per folder basis and retrieve them through another entry point.
- Total number of files in that folder.
- Average number of alphanumeric characters per text file (and standard deviation) in that folder.
- Average word length (and standard deviation) in that folder.
- Total number of bytes stored in that folder.
- Note: All these computations must be calculated recursively from the provided path to the entry point.

### Quick Start

#### Build 
```
// Get Echo (Web Framework)
go get -u github.com/labstack/echo/...
```

#### Test
```
go test -v ./...
```

#### Run 
```
// Server started on localhost:1323
go run .
```