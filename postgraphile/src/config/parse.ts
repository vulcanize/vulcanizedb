import { CONFIG_PATH_KEY } from '../server/config';
import { DatabaseConfig } from '../server/interface';
import { ReadFileSyncCallback, TomlParseCallback } from '../adapters/fs';

export const MISSING_PATH_MESSAGE = `No path to config toml file provided, `
  + `please check the value of ${CONFIG_PATH_KEY} in your environment`;

export const MISSING_HOST_MESSAGE = 'No database host provided in config toml';
export const MISSING_USER_MESSAGE = 'No database user & password '
  + 'provided in config toml';
export const MISSING_DATABASE_MESSAGE = 'No database name provided '
  + 'in config toml';

export function parseConfig(
  readCallback: ReadFileSyncCallback,
  tomlParseCallback: TomlParseCallback,
  configPath?: string
): DatabaseConfig {
  let host = '';
  let port = '';
  let database = '';
  let user = '';
  let password = '';
  let schemas = ['public'];

  if (configPath) {
    const tomlContents = readCallback(`${configPath}`).toString();
    const parsedToml = tomlParseCallback(tomlContents);

    host = parsedToml['database']['hostname'];
    port = parsedToml['database']['port'];
    database = parsedToml['database']['name'];
    user = parsedToml['database']['user'];
    password = parsedToml['database']['password'];
    schemas = parsedToml['database']['schemas'];
  }

  // Overwrite config values with env. vars if such are set
  host = process.env.DATABASE_HOST || host;
  port = process.env.DATABASE_PORT || port;
  database = process.env.DATABASE_NAME || database;
  user = process.env.DATABASE_USER || user;
  password = process.env.DATABASE_PASSWORD || password;

  if (!host) {
    throw new Error(MISSING_HOST_MESSAGE);
  }

  if (!database) {
    throw new Error(MISSING_DATABASE_MESSAGE);
  }

  if (!user || !password) {
    throw new Error(MISSING_USER_MESSAGE);
  }

  return {
    host: `postgres://${user}:${password}@${host}:${port}`,
    database,
    schemas
  };
}
