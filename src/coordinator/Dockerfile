FROM golang:1.16-alpine

WORKDIR /app

RUN apk add git
RUN apk add git nano vim


ENV GO111MODULE=on
RUN go mod init github.com/martina97/SDCCProject_DistributedMutualExclusion
COPY coordinator/*.go /usr/local/go/src/SDCCProject_DistributedMutualExclusion/src/coordinator/
COPY utilities/*.go /usr/local/go/src/SDCCProject_DistributedMutualExclusion/src/utilities/
COPY /peer/tokenAsking/*.go /usr/local/go/src/SDCCProject_DistributedMutualExclusion/src/peer/tokenAsking/

RUN go mod tidy

COPY coordinator/*.go ./coordinator/

RUN go install ./coordinator


#CMD [ "sh", "-c", "${GOPATH}/bin/register" ]

