
<a name="v0.5.3"></a>
## [v0.5.3](https://github.com/rotisserie/eris/compare/v0.5.2...v0.5.3) (2022-04-17)


<a name="v0.5.2"></a>
## [v0.5.2](https://github.com/rotisserie/eris/compare/v0.5.1...v0.5.2) (2022-03-09)


<a name="v0.5.1"></a>
## [v0.5.1](https://github.com/rotisserie/eris/compare/v0.5.0...v0.5.1) (2021-06-22)

### Bug Fixes

* move type reflection so As will work with external errors ([#100](https://github.com/rotisserie/eris/issues/100))


<a name="v0.5.0"></a>
## [v0.5.0](https://github.com/rotisserie/eris/compare/v0.4.1...v0.5.0) (2020-12-25)

### Features

* provide error assertion with As function ([#90](https://github.com/rotisserie/eris/issues/90))


<a name="v0.4.1"></a>
## [v0.4.1](https://github.com/rotisserie/eris/compare/v0.4.0...v0.4.1) (2020-07-08)

### Bug Fixes

* check for zero length wrapPCs when inserting into stacks ([#92](https://github.com/rotisserie/eris/issues/92))


<a name="v0.4.0"></a>
## [v0.4.0](https://github.com/rotisserie/eris/compare/v0.3.0...v0.4.0) (2020-05-20)

### Bug Fixes

* add helper methods to show correct trace in sentry ([#85](https://github.com/rotisserie/eris/issues/85))

### Features

* wrap external errors instead of changing them to root errors ([#84](https://github.com/rotisserie/eris/issues/84))
* try to unwrap external errors during error wrapping ([#80](https://github.com/rotisserie/eris/issues/80))


<a name="v0.3.0"></a>
## [v0.3.0](https://github.com/rotisserie/eris/compare/v0.2.1...v0.3.0) (2020-02-13)

### Bug Fixes

* return correct stack for local/global vars and add stack tests ([#74](https://github.com/rotisserie/eris/issues/74))

### Code Refactoring

* insert frames during error wrapping instead of unpacking ([#70](https://github.com/rotisserie/eris/issues/70))

### Features

* allow error output and stack trace inversion ([#73](https://github.com/rotisserie/eris/issues/73))


<a name="v0.2.1"></a>
## [v0.2.1](https://github.com/rotisserie/eris/compare/v0.2.0...v0.2.1) (2020-01-28)

### Code Refactoring

* check for global stack traces instead of forcing NewGlobal ([#71](https://github.com/rotisserie/eris/issues/71))


<a name="v0.2.0"></a>
## [v0.2.0](https://github.com/rotisserie/eris/compare/v0.1.1...v0.2.0) (2020-01-17)

### Bug Fixes

* add discord invite link ([#65](https://github.com/rotisserie/eris/issues/65))
* copy global root errors to ensure stack traces are isolated ([#58](https://github.com/rotisserie/eris/issues/58), [#59](https://github.com/rotisserie/eris/issues/59))

### Features

* eris errors are now compatible with sentry error tracing ([#60](https://github.com/rotisserie/eris/issues/60))
* improve default formatters and add custom format support ([#57](https://github.com/rotisserie/eris/issues/57))
* improve error wrapping, stack trace management, and formatting ([#46](https://github.com/rotisserie/eris/issues/46))


<a name="v0.1.1"></a>
## [v0.1.1](https://github.com/rotisserie/eris/compare/v0.1.0...v0.1.1) (2019-12-26)

### Bug Fixes

* update mod file


<a name="v0.1.0"></a>
## v0.1.0 (2019-12-25)

### Features

* json and str formats support with custom error type
* finish implementing and testing Is and Cause
* improve default printer and add json printer
* add error type checking
* improve error wrapping and integrate basic default printer
* add stack trace implementation and root error types

