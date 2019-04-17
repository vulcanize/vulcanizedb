import { PassportStatic } from 'passport';

import { PostgraphileMiddleware } from '../../src/adapters/postgraphile';
import { buildServerConfig } from '../../src/server/config';

import {
  DatabaseConfig,
  ServerConfig,
  ServerUtilities
} from '../../src/server/interface';

describe('buildServerConfig', () => {
  let configParser: jasmine.Spy;
  let postgraphileMiddleware: PostgraphileMiddleware;
  let expressSessionHandler: jasmine.Spy;
  let passportInitializer: jasmine.Spy;
  let passportSessionHandler: jasmine.Spy;
  let serverConfig: ServerConfig;

  let serverUtilities: ServerUtilities;
  let databaseConfig: DatabaseConfig;

  beforeEach(() => {
    databaseConfig = {
      host: 'example.com',
      database: 'example_database',
      schemas: ['public'],
      ownerConnectionString: 'postgres://admin:admin@host'
    };

    postgraphileMiddleware = jasmine
      .createSpyObj<PostgraphileMiddleware>(['call']);

    serverUtilities = {
      pluginHook: jasmine.createSpy('pluginHook'),
      enableSubscriptions: jasmine.createSpy('enableSubscriptions'),
      express: jasmine.createSpy('express'),
      expressSession: jasmine.createSpy('expressSession'),
      httpServerFactory: jasmine.createSpy('httpServerFactory'),
      passport: jasmine.createSpyObj<PassportStatic>(['initialize', 'session']),
      postgraphile: jasmine.createSpy('postgraphile')
    };

    const rawConfig: object = { exampleOption: 'example value' };

    configParser = jasmine.createSpy('configParser');
    configParser.and.returnValue(rawConfig);

    expressSessionHandler = jasmine.createSpy('expressSessionHandler');
    passportInitializer = jasmine.createSpy('passportInitializer');
    passportSessionHandler = jasmine.createSpy('passportSessionHandler');

    (serverUtilities.postgraphile as jasmine.Spy)
      .and.returnValue(postgraphileMiddleware);
    (serverUtilities.expressSession as jasmine.Spy)
      .and.returnValue(expressSessionHandler);
    (serverUtilities.passport.initialize as jasmine.Spy)
      .and.returnValue(passportInitializer);
    (serverUtilities.passport.session as jasmine.Spy)
      .and.returnValue(passportSessionHandler);

    serverConfig = buildServerConfig(
      serverUtilities, databaseConfig, undefined);
  });

  it('provides the Postgraphile options', () => {
    expect(serverConfig.options).not.toBeNull();
  });

  it('enables simple subscriptions', () => {
    expect(serverConfig.options.simpleSubscriptions).toBe(true);
  });

  it('it adds the express session handler as the first middleware', () => {
    expect(serverConfig.options.webSocketMiddlewares[0])
      .toBe(expressSessionHandler);
  });

  it('it adds the passport initializer as the second middleware', () => {
    expect(serverConfig.options.webSocketMiddlewares[1])
    .toBe(passportInitializer);
  });

  it('it adds the passport session handler as the third middleware', () => {
    expect(serverConfig.options.webSocketMiddlewares[2])
    .toBe(passportSessionHandler);
  });

  it('provides the database config to Postgraphile', () => {
    expect(serverUtilities.postgraphile).toHaveBeenCalledWith(
      `${databaseConfig.host}/${databaseConfig.database}`,
      databaseConfig.schemas,
      jasmine.any(Object));
  });

  it('provides the Postgraphile middleware', () => {
    expect(serverConfig.middleware).toBe(postgraphileMiddleware);
  });

  it('sets the default server port', () => {
    expect(serverConfig.port).toEqual(3000);
  });

  it('sets an explicit server port', () => {
    const serverConfigWithPort = buildServerConfig(
      serverUtilities, databaseConfig, '1234');

    expect(serverConfigWithPort.port).toEqual(1234);
  });
});
