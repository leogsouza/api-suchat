# Base build image
FROM golang:1.14.2-alpine as build

# Install some dependencies needed to build the project
RUN apk add --update --no-cache bash ca-certificates git gcc g++ libc-dev

WORKDIR /suchat

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

#This is the ‘magic’ step that will download all the dependencies that are specified in 
# the go.mod and go.sum file.
# Because of how the layer caching system works in Docker, the  go mod download 
# command will _ only_ be re-run when the go.mod or go.sum file change 
# (or when we add another docker instruction this line)
RUN go mod download

COPY . .

# And compile the project
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/suchat
#RUN go install ./...

FROM alpine:3.11
# Finally we copy the statically compiled Go binary.
WORKDIR /usr/bin
COPY --from=build /go/bin .
EXPOSE 8000 9999

CMD [ "./suchat" ]