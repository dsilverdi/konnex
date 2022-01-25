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