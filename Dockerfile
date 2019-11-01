FROM gcr.io/distroless/base
COPY pinkie /
ENTRYPOINT ["/pinkie"]
