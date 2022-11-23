FROM golang:1.16-alpine

WORKDIR /app

RUN apk add git
RUN apk add git nano vim


ENV GO111MODULE=on


RUN go mod init github.com/martina97/SDCCProject_DistributedMutualExclusion
RUN go mod tidy
RUN go get -u github.com/manifoldco/promptui

COPY peer/*.go ./peer/
COPY utilities/*.go /usr/local/go/src/SDCCProject_DistributedMutualExclusion/src/utilities/
COPY /peer/ricartAgrawala/*.go /usr/local/go/src/SDCCProject_DistributedMutualExclusion/src/peer/ricartAgrawala/
COPY /peer/tokenAsking/*.go /usr/local/go/src/SDCCProject_DistributedMutualExclusion/src/peer/tokenAsking/

#COPY menu/*.go ./menu/

RUN go install ./peer
#RUN go install ./menu


CMD [ "sh", "-c", "${GOPATH}/bin/peer" ]