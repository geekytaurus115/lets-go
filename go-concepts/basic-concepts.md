#### Golang main features:

- Light Weight
- Concurrency
- Easy to Learn - chhota language

#### Disadvanges

- libraries limited but mainly which are required
  - so apart from ML related work other can do

##### Few Contepts

- Workspace: where we do actual coding
- Root: where go is installed
- After installation of Go, the structure looks like below:

```
go
|-- bin
|-- pkg
|-- src
```

- The `src` directory, is the workspace basically
- In a folder all files' package names will same (naming convention)
- 2 types of packages:
  - core (e.g. 'fmt', 'log', 'net/http')
  - third party (that we import, like from github.com/xyz/abc)
- `main`: package is the package which is Required to run
  - inside `main` package, `main()` is the entry point
- `non-main` packages can't be run, but those can be imported and called
- Inside a folder, multiple packages NOT posible
- Module: `go mod init` - intialize module
  - it creates `go.mod` which manage all dependencies
  - `go mod tidy`: add missing and remove unnecessary libraries
    - it creates `go.sum` file
  - `go mod vendor`: add vendor folder with all dependecies to localize these libraries

#### Important Commands

- `go build file.go`: generate the executable file ( binary ) of the source code, inside the CURRENT directory
- `go install file.go`: do the same but inside the **bin** directory, this can RUN from ANY Location (but for that gopath needs setup correctly in environment)
  - Both of the above commands will generate binary executable, Only for the OS you're using
  - `GOOS=windows GOARCH=amd64 go build file.go`: generate for windows

##### Few another points:

- Go routines take min 4 KB
- Whereas Java's `Thread` takes 1 MB minimum

### Confirm Go Installation

- `go version` : will print the version of Go
- `go env` : will list down all environment variables related to Go
  - GOPATH : is the path where Go is installed
- `echo $PATH` : to check gopath setup or not in environment
  - if path not present, add using below command:
  - export PATH=$PATH:$(go env GOPATH)/bin
  - this export command add the path to environment Temporarily (that means, if terminal close and open --> not available this path)
  - To add permanently:
    - open ~/.zshrc
    - export PATH=$PATH:$(go env GOPATH)/bin
    - save the file
    - source ~/.zshrc
