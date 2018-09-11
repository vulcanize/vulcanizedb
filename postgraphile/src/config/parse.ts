import { CONFIG_PATH_KEY } from '../server/config';
import { DatabaseConfig } from '../server/interface';
import { ReadFileSyncCallback, TomlParseCallback } from '../adapters/fs';

export const MISSING_PATH_MESSAGE = `No path to config toml file provided, `
  + `please check the value of ${CONFIG_PATH_KEY} in your environment`;

export const MISSING_HOST_MESSAGE = 'No database host provided in config toml';
export const MISSING_DATABASE_MESSAGE = 'No database name provided in config '
  + 'toml';

export function parseConfig(
  readCallback: ReadFileSyncCallback,
  tomlParseCallback: TomlParseCallback,
  configPath?: string
): DatabaseConfig {
  if (!configPath || configPath.length < 1) {
    throw new Error(MISSING_PATH_MESSAGE);
  }

  const tomlContents = readCallback(`${configPath}`).toString();
  const parsedToml = tomlParseCallback(tomlContents);

  const host = parsedToml['database']['hostname'];
  const port = parsedToml['database']['port'];
  const database = parsedToml['database']['name'];
  const user = parsedToml['database']['user'] || '';
  const password = parsedToml['database']['password'] || '';

  if (!host || host.length < 1) {
    throw new Error(MISSING_HOST_MESSAGE);
  }

  if (!database || database.length < 1) {
    throw new Error(MISSING_DATABASE_MESSAGE);
  }

  return { host: `postgres://${user}:${password}@${host}:${port}`, database };
}
