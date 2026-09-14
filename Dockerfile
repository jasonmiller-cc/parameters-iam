FROM golang:1.22-alpine AS builder

ARG VERSION=dev
ARG COMMIT=unknown

WORKDIR /build

# Copy go.mod files so the replace directive resolves.
COPY go.mod go.sum* ./
COPY ../parameters-core /parameters-core

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
      -ldflags "-s -w \
        -X github.com/jasonmiller-cc/parameters-core/pkg/version.Version=${VERSION} \
        -X github.com/jasonmiller-cc/parameters-core/pkg/version.Commit=${COMMIT}" \
      -o /parameters-iam ./cmd/server

# ---

FROM gcr.io/distroless/static:nonroot

COPY --from=builder /parameters-iam /parameters-iam

EXPOSE 8080

ENTRYPOINT ["/parameters-iam"]
