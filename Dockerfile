# This container is designed for hosting the Orion
# code analysis application. It's not really suitable
# for much else.

FROM centos:7

ENV PATH=$PATH:/usr/local/go/bin:/root/go/bin
ENV GH_USERNAME=OrionGH
# allow token to be passed in at build time using `--build-arg TOKEN=...`
ARG TOKEN=""
ENV GH_AUTH_TOKEN=$TOKEN
ENV ORION=/root/go/src/github.com/ggrpmnn/orion

RUN mkdir -p /root/go/src/github.com/ggrpmnn && cd /root
RUN yum install gcc git make sqlite-devel which golang ruby ruby-devel -y -q -e 0

RUN cd /root/go/src/github.com/ggrpmnn && git clone https://github.com/ggrpmnn/orion && go install ./orion/src/orion
RUN go get -u github.com/ggrpmnn/orion/src github.com/GoASTScanner/gas
RUN gem install bundler 
#dawnscanner
