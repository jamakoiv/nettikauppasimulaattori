# FROM alpine:latest
FROM debian:bookworm-slim

# Copy local src-folder to /src in container.
COPY . /prog
WORKDIR /prog

CMD ["/prog/main", "-target=test"]
#CMD ["ls", "-la", "/prog"]
