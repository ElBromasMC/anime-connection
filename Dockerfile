FROM golang:1.22-alpine AS builder

WORKDIR /app

RUN apk add --no-cache \
    make \
    npm

RUN go install github.com/a-h/templ/cmd/templ@latest

# Use npm to install building tools
COPY package.json package-lock.json ./
RUN npm install --omit=dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN make ./build/server

# Deploy the application binary into a lean image
FROM alpine:latest

WORKDIR /

COPY --from=builder /app/build/server /server

EXPOSE 8080

ENTRYPOINT ["/server"]
