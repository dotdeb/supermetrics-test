#
#   Application dockerfile
#   !Notice that even API_SECRET is marked as argument it is expected to be passed via kubernetes secret 
#   in the test/production environment.
#
FROM golang

# Args should be resolved during the run
ARG PORT
ARG JWT_ISSUER
ARG JWT_AUDIENCE
ARG TIMEOUT_SEC
ARG API_SECRET
ENV GOPATH=/go

# Run as non-root user
RUN useradd appuser
USER appuser
WORKDIR /home/appuser/go/src/github.com/dotdeb/supermetrics-test

COPY main.go go.mod go.sum ./
COPY internals/ ./internals

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o supermetrics-test

EXPOSE 8000
CMD [ "./supermetrics-test" ]