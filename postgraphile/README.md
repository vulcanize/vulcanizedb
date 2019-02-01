# Vulcanize GraphQL API

This application utilizes Postgraphile to expose GraphQL endpoints for exposure of the varied data that VulcanizeDB tracks.

## Docker use
_Note: currently this image is ~500MB large (unpacked)_

Build the docker image in this directory. Start the `GraphiQL` frontend by:
* Setting the env variables for the database connection: `DATABASE_HOST`,
  `DATABASE_NAME`, `DATABASE_USER`, `DATABASE_PASSWORD` (and optionally
  `DATABASE_PORT` if running on non-standard port).
  * The specified user needs to be `superuser` on the vulcanizeDB database
* Run the container (ex. `docker run -e DATABASE_HOST=localhost -e DATABASE_NAME=vulcanize_public -e DATABASE_USER=vulcanize -e DATABASE_PASSWORD=vulcanize -d postgraphile:latest`)
* GraphiQL is available at `:3000/graphiql`


## Building

*This application assumes the use of the [Yarn package manager](https://yarnpkg.com/en/). The use of npm may produce unexpected results.*

Install dependencies with `yarn` and execute `yarn build`. The bundle produced by Webpack will be present in `build/dist/`.

This application currently uses the Postgraphile supporter plugin. This plugin is present in the `vendor/` directory and is copied to `node_modules/` after installation of packages. It is a fresh checkout of the plugin as of August 31st, 2018.

## Running

Provide the built bundle to node as a runnable script: `node ./build/dist/vulcanize-postgraphile-server.js`

## Testing

Tests are executed via Jasmine with a console reporter via the `yarn test` task.
