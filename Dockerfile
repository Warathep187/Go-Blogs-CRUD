FROM registry.thinknet.co.th/sredev/golang:1.22.1-custom-netrc as build

# Ignore host check on SSH
RUN mkdir -p /root/.ssh && \
  echo "Host *\n\tStrictHostKeyChecking no" > /root/.ssh/config

COPY ./go.mod ./go.sum ./
RUN --mount=type=ssh go mod tidy && go mod download
COPY . .

RUN go install github.com/air-verse/air@latest
CMD ["air", "-c", ".air.toml"]
