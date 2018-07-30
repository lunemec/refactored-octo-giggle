FROM alpine

LABEL version="0.1"
# Public API
EXPOSE 8888
ENV LOGXI="*"

RUN apk add --update ca-certificates

RUN mkdir -p /srv/refactored_octo_giggle

COPY refactored-octo-giggle /srv/refactored_octo_giggle
COPY app.toml /srv/refactored_octo_giggle/

CMD ["/srv/refactored_octo_giggle/refactored-octo-giggle", "--config", "/srv/refactored_octo_giggle/app.toml"]
