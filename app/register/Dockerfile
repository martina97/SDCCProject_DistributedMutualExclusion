FROM golang:1.16-alpine

WORKDIR /app

RUN apk add git
RUN apk add git nano vim


ENV GO111MODULE=on
RUN go mod init github.com/martina97/SDCCProject_DistributedMutualExclusion
COPY utilities/*.go /usr/local/go/src/SDCCProject_DistributedMutualExclusion/app/utilities/

RUN go mod tidy

COPY register/*.go ./register/
#COPY peer1/*.go ./peer1/

RUN go install ./register
#RUN go install ./peer1

#CMD [ "sh", "-c", "${GOPATH}/bin/register" ]

