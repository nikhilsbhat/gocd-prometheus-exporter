### Description: Dockerfile for gocd-prometheus-exporter
FROM alpine:3.16

COPY gocd-prometheus-exporter /

# Starting
ENTRYPOINT [ "/gocd-prometheus-exporter" ]