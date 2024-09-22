# Specify the version of Go to use
FROM golang:1.22 AS builder

ENV GO111MODULE=on CGO_ENABLED=0

# Copy all the files from the host into the container
WORKDIR /app
COPY . .

# Compile the action
RUN go build -o /bin/action

FROM scratch

# Copy over the compiled action from the first step
COPY --from=builder /bin/action /bin/action

# Specify the container's entrypoint as the action
ENTRYPOINT [ "/bin/action" ]
