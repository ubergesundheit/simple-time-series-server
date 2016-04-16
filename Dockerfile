FROM scratch

COPY bin/simple-time-series-server /
COPY simple-time-series-secrets /
COPY simple-time-series-db.db /

EXPOSE 8080

ENTRYPOINT ["/simple-time-series-server"]
