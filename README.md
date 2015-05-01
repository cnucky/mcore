# mcore - Microservice Core
Re-usable codebase for Go microservices.
This is the codebase used at XSNews for quickly creating
microservices without re-inventing the wheel everytime with
regards to the outside/infrastructure.

* config
Abstract the repetetive JSON config parsing.

* dates
Simple string interval parser. Converting '1w', time.Now() to next week's time.Time.

* log
Abstract logging policy (stdout/stderr/prefix/debugmsg)

* inspect
Validate user input by validating the in-memory struct's annotations.
