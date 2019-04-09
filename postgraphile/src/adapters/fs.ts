import { URL } from 'url';

export type ReadFileSyncCallback = (
  path: string | number | Buffer | URL,
  options?: { encoding?: null | undefined; flag?: string | undefined; }
    | null
    | undefined
) => string | Buffer;

export type TomlParseCallback
  = (fileContents: string) => { [key: string]: { [key: string]: string } };
