import { ServerUtilities, DatabaseConfig, ServerConfig } from './interface';

import {
  PostgraphileMiddleware,
  PostgraphileOptions
} from '../adapters/postgraphile';

export const CONFIG_PATH_KEY = 'CONFIG_PATH';
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

  const options: PostgraphileOptions = {
    pluginHook: pluginHook,
    simpleSubscriptions: true,
    watchPg: true,
    enableCors: true,
    graphiql: true,
    webSocketMiddlewares: [
      expressSessionHandler,
      passportInitializer,
      passportSessionHandler
    ]
  };

  const middleware: PostgraphileMiddleware = utilities.postgraphile(
    `${databaseConfig.host}/${databaseConfig.database}`,
    ["public", "maker"],
    options
  );

  return { middleware, options, port: parseInt(port, 10) };
}
