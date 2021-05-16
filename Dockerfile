FROM golang:1.16.0

RUN apt-get update && apt-get install -y \
    bash-completion \
    --no-install-recommends && rm -rf /var/lib/apt/lists/*

RUN mkdir /app
WORKDIR /app
COPY makefile .
RUN make setup-tools

CMD ["bash"]