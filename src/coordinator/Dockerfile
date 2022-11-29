FROM golang:1.16-alpine

WORKDIR /app

RUN apk add git
RUN apk add git nano vim


ENV GO111MODULE=on
RUN go mod init github.com/martina97/SDCCProject_DistributedMutualExclusion
COPY coordinator/*.go /usr/local/go/src/SDCCProject_DistributedMutualExclusion/src/coordinator/

RUN go mod tidy

COPY coordinator/*.go ./coordinator/
#COPY peer1/*.go ./peer1/

RUN go install ./coordinator
#RUN go install ./peer1

#CMD [ "sh", "-c", "${GOPATH}/bin/register" ]

