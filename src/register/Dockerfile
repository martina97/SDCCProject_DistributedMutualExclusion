FROM golang:1.16-alpine

WORKDIR /app

RUN apk add git
RUN apk add git nano vim


ENV GO111MODULE=on
RUN go mod init github.com/martina97/SDCCProject_DistributedMutualExclusion
COPY utilities/*.go /usr/local/go/src/SDCCProject_DistributedMutualExclusion/src/utilities/

RUN go mod tidy

COPY register/*.go ./register/

RUN go install ./register

#CMD [ "sh", "-c", "${GOPATH}/bin/register" ]

