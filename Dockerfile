FROM alpine:latest

ADD . /opt/webapp/
WORKDIR /opt/webapp

# Run the image as a non-root user
RUN adduser -D myuser
USER myuser

# Run the app.  CMD is required to run on Heroku
# $PORT is set by Heroku			
CMD /opt/webapp/main 
