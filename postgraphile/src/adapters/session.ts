import { SessionOptions } from 'express-session';
import { Express, RequestHandler, Handler } from 'express';

export type ExpressInitCallback = () => Express;

export type ExpressSessionInitCallback
  = (options?: SessionOptions) => RequestHandler ;

export type PassportInitCallback = (
  options?: { userProperty: string; } | undefined
) => Handler;

export type PassportSessionCallback = (
  options?: { pauseStream: boolean; } | undefined
) => Handler;

export interface StaticPassportProvider {
  initialize: PassportInitCallback;
  session: PassportSessionCallback;
}
