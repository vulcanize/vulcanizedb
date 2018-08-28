import { createServer } from 'http';
import { postgraphile } from 'postgraphile';
import { readFileSync } from 'fs';

import express = require('express');
import passport = require('passport');
import session = require('express-session');
import toml = require('toml');

const {
  default: PostGraphileSupporter,
  enhanceHttpServerWithSubscriptions,
} = require('@graphile/plugin-supporter');

import { ServerUtilities } from './server/interface';
import { bootServer } from './server/runtime';
import { parseConfig } from './config/parse';

import {
  buildServerConfig,
  CONFIG_PATH_KEY,
  SERVER_PORT_KEY
} from './server/config';

const configPath = process.env[CONFIG_PATH_KEY];
const serverPort = process.env[SERVER_PORT_KEY];

const serverUtilities: ServerUtilities = {
  enableSubscriptions: enhanceHttpServerWithSubscriptions,
  express,
  expressSession: session,
  httpServerFactory: createServer,
  passport,
  postgraphile
};

const databaseConfig = parseConfig(readFileSync, toml.parse, configPath);
const serverConfig = buildServerConfig(
  serverUtilities, databaseConfig, serverPort);

bootServer(serverUtilities, serverConfig);
