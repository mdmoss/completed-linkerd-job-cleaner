FROM golang:1.19-alpine AS build
WORKDIR /completed-linkerd-job-cleaner
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -v -o /usr/bin/completed-linkerd-job-cleaner

FROM alpine AS run
COPY --from=build /usr/bin/completed-linkerd-job-cleaner /usr/bin/completed-linkerd-job-cleaner
RUN adduser -D nonroot
USER nonroot
ENTRYPOINT [ "/usr/bin/completed-linkerd-job-cleaner" ]