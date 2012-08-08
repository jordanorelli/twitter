this is experimental.  I don't recommend it.

In particular, the biggest issue with this library is that the Twitter API
sends malformed json from the perspective of Go; in some cases where a field is
expected to be a string, it is given as null, which causes an error because Go
does not see null and empty string as equivalent.  Go is technically correct in
this regard, but the lack of a `permitnull` option or something similar in a
json struct tag is regrettable.

This isn't really ready to be released; I'm only posting it in response to a
question on the golang-nuts google group (which you should join).  This API is
subject to change and I do not recommend it for anything other than pedagogical
purposes at this time.
