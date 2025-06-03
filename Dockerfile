FROM ubuntu:latest
LABEL authors="wanjalize"

ENTRYPOINT ["top", "-b"]