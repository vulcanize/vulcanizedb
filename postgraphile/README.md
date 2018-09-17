# Vulcanize GraphQL API

This application utilizes Postgraphile to expose GraphQL endpoints for exposure of the varied data that VulcanizeDB tracks.

## Building

*This application assumes the use of the [Yarn package manager](https://yarnpkg.com/en/). The use of npm may produce unexpected results.*

Install dependencies with `yarn` and execute `yarn build`. The bundle produced by Webpack will be present in `build/dist/`.

This application currently uses the Postgraphile supporter plugin. This plugin is present in the `vendor/` directory and is copied to `node_modules/` after installation of packages. It is a fresh checkout of the plugin as of August 31st, 2018.

## Running

Provide the built bundle to node as a runnable script: `node ./build/dist/vulcanize-postgraphile-server.js`

## Testing

Tests are executed via Jasmine with a console reporter via the `yarn test` task.