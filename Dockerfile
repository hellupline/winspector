# syntax=docker/dockerfile:1

FROM golang:1.17-alpine AS build
WORKDIR /build
COPY go.mod go.sum ./
COPY vendor/ vendor/
COPY *.go ./
ENV CGO_ENABLED=0
RUN go build -mod vendor -o /winspector


# FROM gcr.io/distroless/base-debian11
FROM gcr.io/distroless/static
COPY --from=build /winspector /winspector
EXPOSE 8000
USER nonroot:nonroot
ENTRYPOINT ["/winspector"]
