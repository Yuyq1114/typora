FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /mini-operator .

FROM scratch
COPY --from=build /mini-operator /mini-operator
ENTRYPOINT ["/mini-operator"]
