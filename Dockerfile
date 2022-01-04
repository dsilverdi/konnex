# # Start from golang base image
# FROM golang:1.13

# # Set the current working directory inside the container
# WORKDIR /usr/src/app

# # Copy go.mod, go.sum files and download deps
# COPY go.mod go.sum ./
# RUN go mod download
# # RUN go get github.com/rakyll/gotest
# # RUN go get github.com/codegangsta/gin

# # Copy sources to the working directory
# COPY . .

# # Set the Go environment
# ENV GOOS linux
# ENV CGO_ENABLED 1
# ENV GOARCH amd64

# # Run the app
# ARG project
# ENV PROJECT $project
# CMD /go/bin/gin -d ${PROJECT} run main.go
FROM golang:1.14-alpine AS builder

WORKDIR /konnex/app

COPY go.mod go.sum ./
RUN go mod download

ARG revisionID=unknown
ARG buildTimestamp=unknown
ARG SERVICE
# Now copy all the source...
COPY . .

# ...and build it.
RUN CGO_ENABLED=0 go build -o ./konnex/bin/${SERVICE} \
  -ldflags="-s -w -X main.revisionID=${revisionID} -X main.buildTimestamp=${buildTimestamp}" \
  ./cmd/${SERVICE}

# Build the runtime image
FROM alpine:3.11
WORKDIR /root
# set date time to local (Asia/Jakarta)
ENV TZ=Asia/Jakarta

ARG SERVICE

RUN apk add --no-cache --update  tzdata \
  && cp /usr/share/zoneinfo/${TZ} /etc/localtime \
  && echo ${TZ} > /etc/timezone

# install wkhtmltopdf
COPY --from=builder /konnex/app/konnex/bin/${SERVICE} ./service

# HTTP
# EXPOSE 8080

ENTRYPOINT ["./service"]