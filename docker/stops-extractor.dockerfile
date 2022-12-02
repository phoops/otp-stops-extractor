FROM golang:1.18 as builder

WORKDIR /app

RUN mkdir -p -m 0600 ~/.ssh && ssh-keyscan bitbucket.org >> ~/.ssh/known_hosts
RUN git config --global url."git@bitbucket.org:phoops".insteadOf "https://bitbucket.org/phoops"

COPY . /app
RUN go mod download
RUN go build ./cmd/stops-extractor
RUN ls -lah . && chmod +x stops-extractor && pwd

FROM alpine

LABEL maintainer="Phoops info@phoops.it"
LABEL environment="production"
LABEL project="otp-stops-extractor"

RUN apk update && apk add --no-cache tzdata


WORKDIR /app
COPY --from=builder /app/stops-extractor /app

CMD ["./stops-extractor"]