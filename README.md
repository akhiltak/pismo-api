# Transaction API

# Development Dependencies (binary files)
 - Docker **MUST** be installed on the system
 - Install following for development:
   - go install go.uber.org/mock/mockgen@latest
   - go install github.com/amacneil/dbmate@latest
   - go install github.com/swaggo/swag/cmd/swag@latest
   <!-- - go get github.com/golang/mock/mockgen/model -->


### How to run
 - `make run` to run the server inside docker with two containers (app and postgres)
 - Access swagger api docs at: http://localhost:2090/swagger/index.html
 - `make test` to run tests
 - `make stop` to stop the containers
 - There are other make cmds that I use for development (like to gen swagger docs and mocks) - `make swagger` and `make mocks`
 - Please also see screenshot of a test run I did

 ## Error Handling

All endpoints return appropriate HTTP status codes and error messages in case of failures. All handlers and services have debug logs.
Common error scenarios include:

- `400 Bad Request`: For invalid inputs or missing required parameters.
- `404 Bad Request`: For resource not found.
- `500 Internal Server Error`: For server-side errors or issues.


### Assumptions and Tradeoffs:
 - For `operation_types`, spent some time thinking about it. Ideally I would want to keep them as enums if they are predefined but went with having a table for it due to two reasons
   - Create txn POST API needs an operation type id
   - Having as a table, we can have CRUD for it and adding new values does not requires code/schema changes
   - For now, since it wasn't mentioned in the problem statement, just seeded the operation types in the DB via migrations
 - There are `unit tests` for handlers and service layer where the majority of validation and business logic will reside
 - Quickly added `integration tests` now that runs test docker containers for app and postgres (didn't really spent too much time into it though to ensure consistent runs)
 - To see logs for integration tests, just comment out appropriate lines in makefile to avoid cleaning up containers afterwords
 - All the request are logged in server logs
 - The sql queries are logged on each request in debug mode
 - Created handlers separately for account and transaction to keep routes and validation separate
 - But created a single service for now which can easily be split later
 - DB access is abstracted via repo layer
 - Used generics in repo layer so that it's easy to switch to any different ORM if needed, it also avoids a *lot* of duplication
 - `Please note`: Since this was a quick assignment, I didn't follow the git practices (using develop/features branches and creating commit for each task since it was a single task) but I generally create a ticket for each task/feature and tag commit as `JIRA-1234/strong with the force, you are`
