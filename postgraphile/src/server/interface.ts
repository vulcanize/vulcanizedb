import { CreateHttpServerCallback } from '../adapters/http';

import {
  ExpressInitCallback,
  ExpressSessionInitCallback,
  StaticPassportProvider
} from '../adapters/session';

import {
  AddSubscriptionsCallback,
  PostgraphileInitCallback,
  PostgraphileMiddleware,
  PostgraphileOptions
} from '../adapters/postgraphile';
import { PluginHookFn } from 'postgraphile/build/postgraphile/pluginHook';

export interface DatabaseConfig {
  host: string;
  database: string;
  schemas: string[];
  ownerConnectionString: string;
}

export interface ServerConfig {
  middleware: PostgraphileMiddleware;
  options: PostgraphileOptions;
  port: number;
}

export interface ServerUtilities {
  enableSubscriptions: AddSubscriptionsCallback;
  express: ExpressInitCallback;
  expressSession: ExpressSessionInitCallback;
  httpServerFactory: CreateHttpServerCallback;
  passport: StaticPassportProvider;
  postgraphile: PostgraphileInitCallback;
  pluginHook: PluginHookFn;
}
