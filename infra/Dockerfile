FROM golang:latest

RUN curl -fsSL https://deb.nodesource.com/setup_20.x | bash - > /dev/null 2>&1 && \
    apt-get install -y -qq nodejs > /dev/null 2>&1 && \
    npm install -g aws-cdk aws-cdk-local > /dev/null 2>&1
