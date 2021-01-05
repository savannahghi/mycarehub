FROM nginx:latest
LABEL Author="DavidDexter dmwangi@savannahinformatics.com"
ENV NGINX_PORT=8080
EXPOSE 8080
COPY site /opt/
RUN rm -v /etc/nginx/conf.d/default.conf 
COPY site.conf /etc/nginx/conf.d
