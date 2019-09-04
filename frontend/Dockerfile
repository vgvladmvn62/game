FROM node:alpine as builder

ARG PORT=80

WORKDIR /app

# install dependencies
COPY package.json /app/package.json
RUN npm install --no-optional

# copy project files
COPY /src /app/src
COPY /public /app/public

#build our app
RUN npm run build --production


FROM nginx:latest
COPY nginx.conf /etc/nginx/nginx.conf
COPY --from=builder /app/build var/public
EXPOSE ${PORT}
CMD ["nginx", "-g", "daemon off;"]
