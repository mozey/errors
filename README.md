# mozey/errors

Define custom errors with an optional common cause.

This package is based on the archived package `github.com/pkg/errors`, [see roadmap](https://github.com/pkg/errors?tab=readme-ov-file#roadmap). 

Also consider 
- [Error handling in Go](https://github.com/mozey/solid?tab=readme-ov-file#error-handling-in-go)
- [Go2 error proposals](https://go.googlesource.com/proposal/+/master/design/go2draft.md), and try to keep code in here future compatible


## Usage

See tests for usage

```bash
git clone https://github.com/mozey/errors.git
cd errors
go test ./...
```


## Error handling in Go

**Summary**: 
- For maximum flexibility **treat all errors as opaque**
- In situations where you cannot do that, **assert errors for behaviour**, instead of type or value
- When comparing type or value, **minimise the number of Sentinel Errors**

Consider the discussion in [Go Time #16](https://changelog.com/gotime/16):

All the developers on the team are encouraged to think about errors while writing the code. **Exceptions** allow you to delay thinking about how the error will be handled.

Error values may be compared to predefined **Sentinel errors**, e.g. `io.EOF` or `sql.ErrNoRows`. Do this sparingly, avoid comparing errors everywhere in the code. Usually it's enough to return the error all the way up the caller stack. Wrap and unwrap error value as required.

**Who** is the error message (value) for, the *end-user or developers*? Programs shouldn't say something unless there is an error. Stack traces are useful to developers, and may be included with the error.

**Structured logs** are not useful to the end-user, but are useful in development. RequestID (*tracing*) and UserID (*audit logs*) are examples where structured logging is especially useful.

Common Go Mistakes: [Handling an error twice](https://100go.co/#handling-an-error-twice-52). In most situations, an error should be handled only once. Logging an error is handling an error. Therefore, you have to choose between logging or returning an error. In many cases, error wrapping is the solution as it allows you to provide additional context to an error and return the source error.

Don't put metrics in the logs, keep these in a separate system.

Common Go Mistakes: [Unnecessary nested code](https://100go.co/#unnecessary-nested-code-2). Avoid many nested levels, it causes [bad line of sight](https://www.youtube.com/watch?v=yeetIgNeIkc&t=330s), try to keep the happy path aligned on the left. If possible, make the happy return the last statement. This makes it easier to build a mental mode of the code

[Don't just check errors, handle them gracefully](https://dave.cheney.net/2016/04/27/dont-just-check-errors-handle-them-gracefully)
- Never inspect the output of `error.Error`
- **Sentinel errors** become part of your public API. They create a dependency between two packages
- An **error type** is a type that you create that implements the error interface, i.e. custom errors. Error types must be made public
- **Opaque errors**: is when you know an error occurred, but you don't have the ability to see inside the error
- Interactions with the world outside your process, like network activity, require that the caller investigate the nature of the error to decide if it is reasonable to retry the operation. *Assert errors for behaviour, not type*. Errors may implement the Temporary interface, then call `IsTemporary` to determine if the error could be retried
- Add context to errors, use [pkg/errors](github.com/pkg/errors). It has methods like `Wrap`, `Cause`, and `Is`. Can be used to add a stack trace to errors
