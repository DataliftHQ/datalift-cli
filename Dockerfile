FROM golang:1.20-alpine

RUN apk add --no-cache bash \
	curl \
	docker-cli \
	docker-cli-buildx \
	git \
	mercurial \
	make \
	build-base \
	tini

ENTRYPOINT ["/sbin/tini", "--", "/entrypoint.sh"]
CMD [ "-h" ]

COPY scripts/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

COPY datalift_*.apk /tmp/
RUN apk add --allow-untrusted /tmp/datalift_*.apk