FROM ubuntu:18.04
RUN apt-get update && apt-get upgrade -y
RUN apt-get install ca-certificates -y
COPY server/nebular /app/nebular
COPY server/public /app/public
WORKDIR /app
RUN mkdir -p config
ENTRYPOINT ["./nebular"]