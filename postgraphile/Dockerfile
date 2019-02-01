FROM mhart/alpine-node:10
RUN apk --update --no-cache add make g++ python findutils postgresql-dev

WORKDIR /app
COPY . /app
run yarn install
RUN ["./node_modules/typescript/bin/tsc"]
EXPOSE 3000
CMD ["node", "/app/build/dist/src/index.js"]
