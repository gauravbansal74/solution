FROM scratch
MAINTAINER Gaurav Kumar <gauravbansal74@gmail.com>
ADD solution solution
ADD data data
ADD ca-certificates.crt /etc/ssl/certs/
ENV PORT 8000
EXPOSE 8000
ENTRYPOINT ["/solution"]
