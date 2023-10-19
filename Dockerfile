# Start from Python 3.x base image
FROM python:3.11-alpine

RUN apk update update \
    && apk add apache2-utils \
    && apk add bash \
    && mkdir /pypi-server

WORKDIR /pypi-server

# Creating packages directory and Authentication Directory
RUN mkdir packages \
    && mkdir auth

# Installing PyPI server and passlib package for authentication
RUN python3 -m pip install pypiserver passlib

COPY ./docker_entry.sh /pypi-server

# Making the entrypoint file executable
RUN chmod +rx ./docker_entry.sh

ENTRYPOINT ["/pypi-server/docker_entry.sh"]
EXPOSE 80