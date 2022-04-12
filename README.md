# Practice Route handler

[![Go](https://github.com/practice-golang/router-practice/actions/workflows/go.yml/badge.svg)](https://github.com/practice-golang/router-practice/actions/workflows/go.yml)
[![codecov](https://codecov.io/gh/practice-golang/router-practice/branch/main/graph/badge.svg?token=MQKFGED93S)](https://codecov.io/gh/practice-golang/router-practice)

## Practice what
* own router, group-router
* regexp, path/filepath, testing
* tooling usage: rs/cors, alexedwards/scs, rs/zerolog, lestrrat-go/jwx, mitchellh/gox, gobwas/ws, github action, gocov, codecov

## Build
* see [Makefile](/Makefile)

## Route
* see [requests.http](/requests.http), [setup.go](/setup.go)

## Files
|        | embed | Find fs first |
|--------|-------|---------------|
| embed  | yes   | no            |
| html   | yes   | yes           |
| static | no    | yes           |

## Coverage checking
* see [Makefile](/Makefile), [cover.cmd](/cover.cmd)

## Source
* https://gist.github.com/reagent/043da4661d2984e9ecb1ccb5343bf438
