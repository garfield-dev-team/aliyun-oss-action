# Specify the version of Go to use
FROM golang:1.19

# Copy all the files from the host into the container
WORKDIR /app
COPY . .

# Compile the action
RUN go build -o /bin/action

# Specify the container's entrypoint as the action
ENTRYPOINT [ "/bin/action" ]