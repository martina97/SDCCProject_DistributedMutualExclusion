FROM golang:1.16-alpine

WORKDIR /app

RUN apk add git
RUN apk add git nano vim


ENV GO111MODULE=on


RUN go mod init github.com/martina97/SDCCProject_DistributedMutualExclusion
RUN go mod tidy
RUN go get -u github.com/manifoldco/promptui

COPY node/*.go ./node/
COPY utilities/*.go /usr/local/go/src/SDCCProject_DistributedMutualExclusion/app/utilities/
COPY node/*.go /usr/local/go/src/SDCCProject_DistributedMutualExclusion/app/node/

#COPY menu/*.go ./menu/

RUN go install ./node
#RUN go install ./menu

#CMD [ "sh", "-c", "${GOPATH}/bin/peer1" ]