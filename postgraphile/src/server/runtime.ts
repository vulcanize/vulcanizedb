import { ServerUtilities, ServerConfig } from './interface';

export function bootServer(
  utilities: ServerUtilities,
  config: ServerConfig
): void {
  const expressApp = utilities.express();
  expressApp.use(config.middleware);

  const httpServer = utilities.httpServerFactory(expressApp);

  utilities.enableSubscriptions(
    httpServer,
    config.middleware,
    config.options);
  
  httpServer.listen(config.port);
}
