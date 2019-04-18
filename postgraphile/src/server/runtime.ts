import { ServerUtilities, ServerConfig } from './interface';

export function bootServer(
  utilities: ServerUtilities,
  config: ServerConfig
): void {
  const expressApp = utilities.express();
  expressApp.use(config.middleware);

  const httpServer = utilities.httpServerFactory(expressApp);
  
  httpServer.listen(config.port);
}
