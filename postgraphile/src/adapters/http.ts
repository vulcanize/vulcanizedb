import { IncomingMessage, ServerResponse, Server } from 'http';

export type RequestListenerCallback
  = (request: IncomingMessage, response: ServerResponse) => void;

export type CreateHttpServerCallback
  = (requestListener?: RequestListenerCallback) => Server;
