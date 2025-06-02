FROM golang:1.24 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /student-api

# Use a distroless base image for the final build
FROM gcr.io/distroless/base-debian12  AS build-release-stage

WORKDIR /

COPY --from=build-stage /student-api /student-api

EXPOSE 8090

USER nonroot:nonroot

ENTRYPOINT ["/student-api"]

