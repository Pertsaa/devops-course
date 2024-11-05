FROM nginx:stable-alpine

COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY .htpasswd /etc/nginx/.htpasswd
COPY index.html /etc/nginx/html/index.html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]