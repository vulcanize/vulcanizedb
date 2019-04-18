import { RequestHandler } from 'express';
import {PluginHookFn } from 'postgraphile/build/postgraphile/pluginHook';
import {Plugin} from 'postgraphile';

// NOTE: Shape of the middleware is not
// currently important to this application, but if a need arises,
// any needed shape can be assigned from a custom type here. For
// the time being, this is a named stub to provide clarity.
export interface PostgraphileMiddleware extends RequestHandler {}

export interface PostgraphileOptions {
  appendPlugins: Plugin[];
  disableDefaultMutations: boolean;
  enableCors: boolean;
  exportGqlSchemaPath: string;
  graphiql: boolean;
  ignoreRBAC: boolean;
  ownerConnectionString: string;
  pluginHook: PluginHookFn;
  watchPg: boolean;
  // NOTE Shape of the middlewares is not
  // currently important to this application, but if a need arises,
  // any needed shape can be assigned from a custom type here.
  webSocketMiddlewares: object[];
}

export type PostgraphileInitCallback = (
  databaseUrl: string,
  schemas: string[],
  options: PostgraphileOptions
) => PostgraphileMiddleware;

