import { RequestHandler } from 'express';
import { Server } from 'http';

// NOTE (jchristie@8thlight.com) Shape of the middleware is not
// currently important to this application, but if a need arises,
// any needed shape can be assigned from a custom type here. For
// the time being, this is a named stub to provide clarity.
export interface PostgraphileMiddleware extends RequestHandler {}

export interface PostgraphileOptions {
  simpleSubscriptions: boolean;
  // NOTE (jchristie@8thlight.com) Shape of the middlewares is not
  // currently important to this application, but if a need arises,
  // any needed shape can be assigned from a custom type here.
  webSocketMiddlewares: object[];
}

export type PostgraphileInitCallback = (
  databaseUrl: string,
  databaseName: string,
  options: PostgraphileOptions
) => PostgraphileMiddleware;

export type AddSubscriptionsCallback = (
  httpServer: Server,
  middleware: PostgraphileMiddleware,
  options: PostgraphileOptions
) => void;
