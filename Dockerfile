FROM ubuntu:18.04

VOLUME /root/input

COPY ./main /root/app/
RUN chmod +x /root/app/main

WORKDIR /root/app
ENTRYPOINT ["/home/linfan.wty/tianchi/pilotset/src/demo/main", "daemon"]
