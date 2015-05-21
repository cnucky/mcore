# VALID
Internally this library works by recursively parsing a struct, getting all the "validate" tags from it,
then creating tokens from the tags (see: valid/lexer.go), actual validation is then performed by valid/valdsl.go

* tests/valid_test.go - Test case.
* valid/context.go    - Context class for calling validation functions, has the actual struct we're validating.
* valid/helper.go     - Helper functions for validation functions defined in predefined.go
* valid/lexer.go      - Parses a struct for "validate" tags and tokenizes it.
* valid/predefined.go - Validation functions.
* valid/valdsl.go     - Parses tokens and performs validation.

# EXAMPLES
* See tests/valid_test.go

# TODO
* Refactor valdsl.go token parser.
* Add a syntax checker to valdsl.
* Finish predefined functions.
* Cleanup helper.go functions.
