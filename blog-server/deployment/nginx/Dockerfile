FROM nginx

RUN mkdir -p /usr/share/nginx/html/vue/blog
RUN mkdir -p /usr/share/nginx/html/vue/admin
RUN mkdir -p /usr/share/upload/static
RUN mkdir -p /etc/nginx/cert
COPY ./cert/ /etc/nginx/cert/
COPY ./frontend/admin/ /usr/share/nginx/html/vue/blog/
COPY ./frontend/blog/ /usr/share/nginx/html/vue/admin/