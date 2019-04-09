import { RequestHandler } from 'express';
import { Server } from 'http';
import { PluginHookFn } from 'postgraphile/build/postgraphile/pluginHook';

// NOTE: Shape of the middleware is not
// currently important to this application, but if a need arises,
// any needed shape can be assigned from a custom type here. For
// the time being, this is a named stub to provide clarity.
export interface PostgraphileMiddleware extends RequestHandler {}

export interface PostgraphileOptions {
  pluginHook: PluginHookFn,
  simpleSubscriptions: boolean;
  watchPg: boolean;
  enableCors: boolean;
  graphiql: boolean;
  // NOTE: Shape of the middlewares is not
  // currently important to this application, but if a need arises,
  // any needed shape can be assigned from a custom type here.
  webSocketMiddlewares: object[];
}

export type PostgraphileInitCallback = (
  databaseUrl: string,
  schemas: string[],
  options: PostgraphileOptions
) => PostgraphileMiddleware;

export type AddSubscriptionsCallback = (
  httpServer: Server,
  middleware: PostgraphileMiddleware,
  options: PostgraphileOptions
) => void;
