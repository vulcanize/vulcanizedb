import { parseConfig } from '../../src/config/parse';

describe('parseConfig', () => {
  let configPath: string;
  let tomlContents: string;
  let parsedToml: object;
  let readCallback: jasmine.Spy;
  let tomlParseCallback: jasmine.Spy;

  beforeEach(() => {
    configPath = '/example/config/path.toml';
    tomlContents = `[database]\nname = 'example_database'\n `
      + `hostname = 'example.com'\nport = 1234`;

    parsedToml = {
      database: {
        hostname: 'example.com',
        name: 'example_database',
        port: '1234'
      }
    };

    readCallback = jasmine.createSpy('readCallback');
    readCallback.and.returnValue(tomlContents);

    tomlParseCallback = jasmine.createSpy('tomlParseCallback');
    tomlParseCallback.and.returnValue(parsedToml);
  });

  it('provides the database host', () => {
    const databaseConfig = parseConfig(
      readCallback, tomlParseCallback, configPath);

    expect(databaseConfig.host).toEqual('postgres://example.com:1234');
  });

  it('provides the database name', () => {
    const databaseConfig = parseConfig(
      readCallback, tomlParseCallback, configPath);

    expect(databaseConfig.database).toEqual('example_database');
  });

  it('handles a missing config path', () => {
    const failingCall = () =>
      parseConfig(readCallback, tomlParseCallback, '');

    tomlParseCallback.and.returnValue({
      database: { hostname: 'example.com', name: 'example_database' }
    });

    expect(failingCall).toThrow();    
  });

  it('handles a missing database host', () => {
    const failingCall = () =>
      parseConfig(readCallback, tomlParseCallback, configPath);

    tomlParseCallback.and.returnValue({
      database: { hostname: '', name: 'example_database' }
    });

    expect(failingCall).toThrow();    
  });

  it('handles a missing database name', () => {
    const failingCall = () =>
      parseConfig(readCallback, tomlParseCallback, configPath);

    tomlParseCallback.and.returnValue({
      database: { hostname: 'example.com', name: '', port: '1234' }
    });

    expect(failingCall).toThrow();    
  });
});
