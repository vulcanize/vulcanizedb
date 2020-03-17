FROM node:alpine

RUN npm install -g postgraphile
RUN npm install -g postgraphile-plugin-connection-filter
RUN npm install -g @graphile/pg-pubsub

EXPOSE 5000
ENTRYPOINT ["postgraphile"]