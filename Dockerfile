FROM golang:1.14-alpine

# do update and add certs
RUN apk update && apk upgrade -U -a && apk --no-cache add ca-certificates

# create appuser manager user
RUN addgroup -S appuser && adduser -S appuser -G appuser

# server port
EXPOSE 8080

# copy maxmind db to proper location
RUN mkdir -p /usr/local/geocheck/maxmind
COPY data/GeoLite2-Country.mmdb /usr/local/geocheck/maxmind/

# copy static files
RUN mkdir /home/appuser/public
COPY public/ /home/appuser/public/

# change perms on maxmind file and static web files
RUN chown -R appuser:appuser /usr/local/geocheck/maxmind
RUN chown -R appuser:appuser /home/appuser/public

# copy server executable
ADD bin/server /usr/local/bin

# set location of maxmind file
ENV DB_DIR /usr/local/geocheck/maxmind/
ENV DB_NAME GeoLite2-Country.mmdb
ENV HUMANREADABLE True

# set user to local before running app
USER appuser

# set home dir as working dir
WORKDIR /home/appuser

ENTRYPOINT ["/usr/local/bin/server"]
