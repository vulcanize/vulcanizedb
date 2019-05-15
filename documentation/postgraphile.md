# Postgraphile

You can expose VulcanizeDB data via [Postgraphile](https://github.com/graphile/postgraphile).
Check out [their documentation](https://www.graphile.org/postgraphile/) for the most up-to-date instructions on installing, running, and customizing Postgraphile.

## Simple Setup

As of April 30, 2019, you can run Postgraphile pointed at the default `vulcanize_public` database with the following commands:

```
npm install -g postgraphile
postgraphile --connection postgres://localhost/vulcanize_public --schema=public,custom --disable-default-mutations --no-ignore-rbac
```

Arguments:
- `--connection` specifies the database. The above command connects to the default `vulcanize_public` database 
defined in [the example config](../environments/public.toml.example).
- `--schema` defines what schema(s) to expose. The above exposes the `public` schema (for core VulcanizeDB data) as well as a `custom` schema (where `custom` is the name of a schema defined in executed transformers).
- `--disable-default-mutations` prevents Postgraphile from exposing create, update, and delete operations on your data, which are otherwise enabled by default.
- `--no-ignore-rbac` ensures that Postgraphile will only expose the tables, columns, fields, and query functions that the user has explicit access to.

## Customizing Postgraphile

By default, Postgraphile will expose queries for all tables defined in your chosen database/schema(s), including [filtering](https://www.graphile.org/postgraphile/filtering/) and [auto-discovered relations](https://www.graphile.org/postgraphile/relations/).

If you'd like to expose more customized windows into your data, there are some techniques you can apply when writing migrations:

- [Computed columns](https://www.graphile.org/postgraphile/computed-columns/) enable you to derive additional fields from types defined in your database.
For example, you could write a function to expose a block header's state root over Postgraphile with a computed column - without modifying the `public.headers` table.
- [Custom queries](https://www.graphile.org/postgraphile/custom-queries/) enable you to provide on-demand access to more complex data (e.g. the product of joining and filtering several tables' data based on a passed argument).
For example, you could write a custom query to get the block timestamp for every transaction originating from a given address.
- [Subscriptions](https://www.graphile.org/postgraphile/subscriptions/) enable you to publish data as it is coming into your database.

The above list is not exhaustive - please see the Postgraphile documentation for a more comprehensive and up-to-date description of available features.