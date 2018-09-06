import { Express } from 'express';
import { PassportStatic } from 'passport';
import { Server } from 'http';

import { ServerUtilities, ServerConfig } from '../../src/server/interface';
import { bootServer } from '../../src/server/runtime';
import { PostgraphileMiddleware } from '../../src/adapters/postgraphile';

describe('bootServer', () => {
  let serverUtilities: ServerUtilities;
  let serverConfig: ServerConfig;
  let mockExpressApp: Express;
  let mockHttpServer: Server;

  beforeEach(() => {
    serverUtilities = {
      pluginHook: jasmine.createSpy('pluginHook'),
      enableSubscriptions: jasmine.createSpy('enableSubscriptions'),
      express: jasmine.createSpy('express'),
      expressSession: jasmine.createSpy('expressSession'),
      httpServerFactory: jasmine.createSpy('httpServerFactory'),
      passport: jasmine.createSpyObj<PassportStatic>(['initialize', 'session']),
      postgraphile: jasmine.createSpy('postgraphile')
    };

    serverConfig = {
      middleware: jasmine.createSpyObj<PostgraphileMiddleware>(['call']),
      options: { 
        pluginHook: jasmine.createSpy('pluginHook'),
        simpleSubscriptions: true,
        graphiql: true,
        webSocketMiddlewares: [] },
      port: 5678
    };

    mockExpressApp = jasmine.createSpyObj<Express>(['use']);
    (serverUtilities.express as jasmine.Spy)
      .and.returnValue(mockExpressApp);

    mockHttpServer = jasmine.createSpyObj<Server>(['listen']);
    (serverUtilities.httpServerFactory as jasmine.Spy)
      .and.returnValue(mockHttpServer);

    bootServer(serverUtilities, serverConfig);
  });

  it('builds a new, Node HTTP server', () => {
    expect(serverUtilities.httpServerFactory)
      .toHaveBeenCalledWith(mockExpressApp);
  });

  it('provides Postgraphile middleware to the Express app', () => {
    const useSpy = mockExpressApp.use as jasmine.Spy;
    expect(useSpy).toHaveBeenCalledWith(serverConfig.middleware);
  });

  it('enahances the Node HTTP server with Postgraphile subscriptions', () => {
    expect(serverUtilities.enableSubscriptions)
      .toHaveBeenCalledWith(
        mockHttpServer,
        serverConfig.middleware,
        serverConfig.options);
  });

  it('instructs the server to listen on the given port', () => {
    const listenSpy = mockHttpServer.listen as jasmine.Spy;
    expect(listenSpy).toHaveBeenCalledWith(serverConfig.port);
  });
});
