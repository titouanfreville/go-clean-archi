# go-clean-archi

Go clean archi aims to propose a base skeletons for clean architecture projects with rest API or command line
interfaces.

I am not a specialist in the domain but this is the result of a year of experiencing clean architecture on a company
project, discussions with experienced people and idea from conferences on the subject.

## Organization

This proposition revolves on a main separation between domains and implementation regrouping all domains in a single
package and splitting implementation between:

- services for external resources
- core for global utilities
- infrastructure for direct feature implementation and transport split by transport (cmd, routes, grpc, ....)

### Domains

In domain will go everything that defines your application. Witch entities it provides and receive, the functional
errors it has to manage, usecases and interfaces used to describe usecases.

Each usecase function should provide a solution to a specific problem the application has to solve.

### Core

Core is not mandatory, but it is the best place to provide simple tool witch will be used all around your applications.

It is the place I would put a database connection management system for example as it will be same everywhere in the
application, or a JSON responder for HTTP requests.

As long as more than X non domains modules use the same code, it should go to core.

### Adapters

Adapters define the implementation required to resolve usecases. Here goes the SQL requests for example to persist data,
algorithms to compute delivery dates or anything your application has to do.

### Services

Service implements base structures for any external tools you are using. It would be the place to implement a manager
for google PubSub witch will resolve the desired instance and provide some global logic. Then you will just have to rely
on this bloc to defines your infrastructure code

### Transport: cmd/routes ....

Every action linking your users with the application should go in a transport package.

Transport could be group in a single package, but I decided to separate them here cause I did not encounter cases where
command line, rpc communication or rest objects shared definitions.

Fell free to rename and adapt this proposition to your need :)

## Tools around

### Dependency injection

To manage lifecycle and dependencies injection in this kind of project, I was
presented [uber fx](https://pkg.go.dev/go.uber.org/fx)
witch is not the simplest tool I found, but the most comfortable on witch keep go spirit.

### [Go Convey](https://github.com/smartystreets/goconvey)

While the releases are not legion on the project, I still love to use it to write clear tests. It provides a lot of
simple structure and assertions for testing though some are outdated. It is also easy to write your own asserters to
keep a clear test syntaxes.

### Godog + Kactus

[Godog](https://github.com/cucumber/godog) is the official library for writing Cucumber tests in go and I like gherkin
syntax to test globally a feature. As I do not really use it for BDD development as intended and have a lot of
repetitive steps, I add Kactus library to help.

[Kactus](https://github.com/elmagician/kactus) is a library I developed with the help
of [@bnadim](https://github.com/bnadim), [@ftoufet](),
[@jclebreton](https://github.com/jclebreton). It provides basic steps and tools to write cleaner test with gherkin. It
mainly provides a variable management, fixture management system, database/pubsub/rest helpers for testing.

### Wiremock

[Wiremock](http://wiremock.org/) is a great tool to mock external APIs so I like to use it for global and integration
tests. 
