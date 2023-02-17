FROM node:14-alpine

RUN npm install -g @mockoon/cli@2.3.0
COPY mockoon-novomock.json ./mockoon-novomock.json

# Do not run as root.
RUN adduser --shell /bin/sh --disabled-password --gecos "" mockoon
RUN chown -R mockoon ./mockoon-novomock.json
USER mockoon

EXPOSE 3001

ENTRYPOINT ["mockoon-cli", "start", "--hostname", "0.0.0.0", "--daemon-off", "--data", "mockoon-novomock.json", "--container"]

# Usage: docker run -p <host_port>:<container_port> mockoon-test