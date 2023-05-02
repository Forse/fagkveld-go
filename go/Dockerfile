FROM golang:1.20 AS build

# Add non-root user as a good security practice
RUN useradd -u 1001 -m todo-api

WORKDIR /app

# Fetch dependencies
COPY go.mod go.sum ./
RUN go mod download
RUN go mod verify

# Build API
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /api cmd/todo-api/main.go

FROM scratch

# Copy our non-root user
COPY --from=build /etc/passwd /etc/passwd
USER 1001

COPY --from=build /api /api

EXPOSE 3000

ENTRYPOINT ["/api"]
