import { ServerUtilities, DatabaseConfig, ServerConfig } from './interface';

import {
  PostgraphileMiddleware,
  PostgraphileOptions
} from '../adapters/postgraphile';

export const CONFIG_PATH_KEY = 'POSTGRAPHILE_CONFIG_PATH';
export const SERVER_PORT_KEY = 'SERVER_PORT';

const DEFAULT_SERVER_PORT = '3000';

export function buildServerConfig(
  utilities: ServerUtilities,
  databaseConfig: DatabaseConfig,
  port?: string
): ServerConfig {
  if (!port || port.length < 1) {
    port = DEFAULT_SERVER_PORT;
  }

  const expressSessionHandler = utilities.expressSession();
  const passportInitializer = utilities.passport.initialize();
  const passportSessionHandler = utilities.passport.session();
  const pluginHook = utilities.pluginHook;
  const PgSimplifyInflectorPlugin = require('@graphile-contrib/pg-simplify-inflector');

  const options: PostgraphileOptions = {
    appendPlugins: [PgSimplifyInflectorPlugin],
    disableDefaultMutations: databaseConfig.disableDefaultMutations,
    enableCors: true,
    exportGqlSchemaPath: 'schema.graphql',
    graphiql: true,
    ignoreRBAC: false,
    ownerConnectionString: databaseConfig.ownerConnectionString,
    pluginHook: pluginHook,
    watchPg: true,
    webSocketMiddlewares: [
      expressSessionHandler,
      passportInitializer,
      passportSessionHandler
    ]
  };

  const middleware: PostgraphileMiddleware = utilities.postgraphile(
    `${databaseConfig.host}/${databaseConfig.database}`,
    databaseConfig.schemas,
    options
  );

  return { middleware, options, port: parseInt(port, 10) };
}
