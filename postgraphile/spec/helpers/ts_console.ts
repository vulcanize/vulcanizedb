// NOTE (jchristie@8thlight.com) This file helps Jasmine
// comprehend TS sourcemaps by installing a reporter
// specifically for TypeScript
const TSConsoleReporter = require('jasmine-ts-console-reporter');

jasmine.getEnv().clearReporters();
jasmine.getEnv().addReporter(new TSConsoleReporter());
