FROM golang:1.19-alpine AS build
WORKDIR /linkerd-completed-job-cleaner
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -v -o /usr/bin/linkerd-completed-job-cleaner

FROM alpine AS run
COPY --from=build /usr/bin/linkerd-completed-job-cleaner /usr/bin/linkerd-completed-job-cleaner
RUN adduser -D nonroot
USER nonroot
ENTRYPOINT [ "/usr/bin/linkerd-completed-job-cleaner" ]