# FROM alpine:latest
FROM debian:bookworm-slim

RUN apt-get update
RUN apt-get install apt-transport-https ca-certificates gnupg curl -y 

RUN curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | gpg --dearmor -o /usr/share/keyrings/cloud.google.gpg
RUN echo "deb [signed-by=/usr/share/keyrings/cloud.google.gpg] https://packages.cloud.google.com/apt cloud-sdk main" | tee -a /etc/apt/sources.list.d/google-cloud-sdk.list
RUN apt-get update -y && apt-get install google-cloud-cli -y 

# Copy local src-folder to /src in container.
COPY . /prog
COPY ap.json /prog/ap.json
ENV GOOGLE_APPLICATION_CREDENTIALS=/prog/ap.json

WORKDIR /prog

CMD ["/prog/main", "-target=prod"]
#CMD ["ls", "-la", "/prog"]
