# Vulcanize GraphQL API

This application utilizes Postgraphile to expose GraphQL endpoints for exposure of the varied data that VulcanizeDB tracks.

## Docker use
_Note: currently this image is ~500MB large (unpacked)_

Build the docker image in this directory. Start the `GraphiQL` frontend by:
* Setting the env variables for the database connection: `DATABASE_HOST`,
  `DATABASE_NAME`, `DATABASE_USER`, `DATABASE_PASSWORD` (and optionally
  `DATABASE_PORT` if running on non-standard port).
  * The specified user needs to be `superuser` on the vulcanizeDB database,
    so postgraphile can setup watch fixtures keeping track of live schema
    changes.
* To limit the amount of available queries in GraphQL, a restricted user can be used
  for postgraphile introspection by adding env variables `GQ_USER` and `GQ_PASSWORD`.
  * By doing `GRANT [SELECT | EXECUTE]` on tables/functions for this user,
    you can selectively assign things you want available in GraphQL.
  * You still need to pass in a superuser with `DATABASE_USER` & `DATABASE_PASSWORD` for
    for the postgraphile watch fixtures to work.
* By default, postgraphile publishes the `public` schema. This can be expanded with for example `GQ_SCHEMAS=public,maker`
* Run the container (ex. `docker run -e DATABASE_HOST=localhost -e DATABASE_NAME=my_database -e DATABASE_USER=superuser -e DATABASE_PASSWORD=superuser -e GQ_USER=graphql -e GQ_PASSWORD=graphql -e GQ_SCHEMAS=public,anotherSchema -d my-postgraphile-image`)
* GraphiQL frontend is available at `:3000/graphiql`
  GraphQL endpoint is available at `:3000/graphql`

By default, this build will expose only the "public" schema - to add other schemas, use either the env variables,
or a config file `config.toml` and set the env var `CONFIG_PATH` to point to its location. Example `toml`:

```
[database]
    name     = "vulcanize_public"
    hostname = "localhost"
    port = 5432
    gq_schemas = ["public", "yourschema"]
    gq_user = "graphql"
    gq_password = "graphql"
```

## Building

*This application assumes the use of the [Yarn package manager](https://yarnpkg.com/en/). The use of npm may produce unexpected results.*

Install dependencies with `yarn` and execute `yarn build`. The bundle produced by Webpack will be present in `build/dist/`.

## Running

Provide the built bundle to node as a runnable script: `node ./build/dist/vulcanize-postgraphile-server.js`

## Testing

Tests are executed via Jasmine with a console reporter via the `yarn test` task.
