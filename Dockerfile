ARG base=ubuntu:18.04
FROM ${base}

RUN apt-get update && apt-get install curl -y
# Install htmltest into $PATH
RUN curl https://htmltest.wjdp.uk | sudo bash -s -- -b /usr/local/bin
