FROM docker.io/golang:1.16.3 AS build_base
LABEL maintainer="jtbonhomme@gmail.com"

# use go module by default
ENV GO111MODULE=on

WORKDIR /asteboids/

# We want to populate the module cache based on the go.{mod,sum} files. 
COPY go.mod .
COPY go.sum .

#This is the ‘magic’ step that will download all the dependencies that are specified in 
# the go.mod and go.sum file.

# Because of how the layer caching system works in Docker, the go mod download 
# command will _ only_ be re-run when the go.mod or go.sum file change 
# (or when we add another docker instruction this line) 
RUN go mod download

FROM build_base AS service_builder

# Copy code source dirs in workdir
COPY . .

# Build go executable
RUN CGO_ENABLED=0 go build -o main \
  -ldflags "-X github.com/jtbonhomme/asteboids/internal/version.Tag=$(git describe --tags) \
  -X github.com/jtbonhomme/asteboids/internal/version.GitCommit=$(git rev-parse --short HEAD) \
  -X github.com/jtbonhomme/asteboids/internal/version.BuildTime=$(date -u +%FT%T%z)" \
  cmd/asteboids/main.go

# Expose server listening port
FROM alpine:3.8

COPY --from=server_builder /asteboids/main /
RUN chmod +x /main
CMD ["/main"]
