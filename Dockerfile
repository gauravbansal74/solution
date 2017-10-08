FROM scratch
MAINTAINER Gaurav Kumar <gauravbansal74@gmail.com>
ADD solution solution
ADD data data
ADD ca-certificates.crt /etc/ssl/certs/
ENV PORT 9000
EXPOSE 9000
ENTRYPOINT ["/solution"]
