gosanitize
==========

Sanitize and validate input using the JSON Schema v4 draft.

Provides a generic way of loading JSON schemas, converting posted data to pure go structs,
respecting datatypes, validating them against the given JSON schema (with a third party library),
and optionally applying secondary rules to either the whole dataset or to a single field.

todo
====
Expose JSON schema to clients so the same schema can be used for client and server side validation.
Helper function pass http request object.
Length of posted arrays must be consistent.
Detect excess form data.
