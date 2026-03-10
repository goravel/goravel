<h1 align="center">AtomicGo | keyboard</h1>

<p align="center">

<a href="https://github.com/atomicgo/keyboard/releases">
<img src="https://img.shields.io/github/v/release/atomicgo/keyboard?style=flat-square" alt="Latest Release">
</a>

<a href="https://codecov.io/gh/atomicgo/keyboard" target="_blank">
<img src="https://img.shields.io/github/actions/workflow/status/atomicgo/keyboard/go.yml?style=flat-square" alt="Tests">
</a>

<a href="https://codecov.io/gh/atomicgo/keyboard" target="_blank">
<img src="https://img.shields.io/codecov/c/gh/atomicgo/keyboard?color=magenta&logo=codecov&style=flat-square" alt="Coverage">
</a>

<a href="https://codecov.io/gh/atomicgo/keyboard">
<!-- unittestcount:start --><img src="https://img.shields.io/badge/Unit_Tests-1-magenta?style=flat-square" alt="Unit test count"><!-- unittestcount:end -->
</a>

<a href="https://github.com/atomicgo/keyboard/issues">
<img src="https://img.shields.io/github/issues/atomicgo/keyboard.svg?style=flat-square" alt="Issues">
</a>

<a href="https://opensource.org/licenses/MIT" target="_blank">
<img src="https://img.shields.io/badge/License-MIT-yellow.svg?style=flat-square" alt="License: MIT">
</a>

</p>

---

<p align="center">
<strong><a href="#install">Get The Module</a></strong>
|
<strong><a href="https://pkg.go.dev/atomicgo.dev/keyboard#section-documentation" target="_blank">Documentation</a></strong>
|
<strong><a href="https://github.com/atomicgo/atomicgo/blob/main/CONTRIBUTING.md" target="_blank">Contributing</a></strong>
|
<strong><a href="https://github.com/atomicgo/atomicgo/blob/main/CODE_OF_CONDUCT.md" target="_blank">Code of Conduct</a></strong>
</p>

---

<p align="center">
  <img src="https://raw.githubusercontent.com/atomicgo/atomicgo/main/assets/header.png" alt="AtomicGo">
</p>

<p align="center">
<table>
<tbody>
<td align="center">
<img width="2000" height="0"><br>
  -----------------------------------------------------------------------------------------------------
<img width="2000" height="0">
</td>
</tbody>
</table>
</p>
<h3  align="center"><pre>go get atomicgo.dev/keyboard</pre></h3>
<p align="center">
<table>
<tbody>
<td align="center">
<img width="2000" height="0"><br>
   -----------------------------------------------------------------------------------------------------
<img width="2000" height="0">
</td>
</tbody>
</table>
</p>

## Description

Package keyboard can be used to read key presses from the keyboard, while in a
terminal application. It's crossplatform and keypresses can be combined to check
for ctrl+c, alt+4, ctrl-shift, alt+ctrl+right, etc. It can also be used to
simulate (mock) keypresses for CI testing.

Works nicely with https://atomicgo.dev/cursor

## Simple Usage

```go
keyboard.Listen(func(key keys.Key) (stop bool, err error) {
  if key.Code == keys.CtrlC {
    return true, nil // Stop listener by returning true on Ctrl+C
  }

  fmt.Println("\r" + key.String()) // Print every key press
  return false, nil // Return false to continue listening
})
```

## Advanced Usage

```go
// Stop keyboard listener on Escape key press or CTRL+C.
// Exit application on "q" key press.
// Print every rune key press.
// Print every other key press.
keyboard.Listen(func(key keys.Key) (stop bool, err error) {
  switch key.Code {
  case keys.CtrlC, keys.Escape:
    return true, nil // Return true to stop listener
  case keys.RuneKey: // Check if key is a rune key (a, b, c, 1, 2, 3, ...)
    if key.String() == "q" { // Check if key is "q"
      fmt.Println("\rQuitting application")
      os.Exit(0) // Exit application
    }
    fmt.Printf("\rYou pressed the rune key: %s\n", key)
  default:
    fmt.Printf("\rYou pressed: %s\n", key)
  }

  return false, nil // Return false to continue listening
})
```

## Simulate Key Presses (for mocking in tests)

```go
go func() {
  keyboard.SimulateKeyPress("Hello")             // Simulate key press for every letter in string
  keyboard.SimulateKeyPress(keys.Enter)          // Simulate key press for Enter
  keyboard.SimulateKeyPress(keys.CtrlShiftRight) // Simulate key press for Ctrl+Shift+Right
  keyboard.SimulateKeyPress('x')                 // Simulate key press for a single rune
      keyboard.SimulateKeyPress('x', keys.Down, 'a') // Simulate key presses for multiple inputs

  keyboard.SimulateKeyPress(keys.Escape) // Simulate key press for Escape, which quits the program
}()

keyboard.Listen(func(key keys.Key) (stop bool, err error) {
  if key.Code == keys.Escape || key.Code == keys.CtrlC {
    os.Exit(0) // Exit program on Escape
  }

  fmt.Println("\r" + key.String()) // Print every key press
  return false, nil                // Return false to continue listening
})
```

## Usage

#### func  Listen

```go
func Listen(onKeyPress func(key keys.Key) (stop bool, err error)) error
```
Listen calls a callback function when a key is pressed.

Simple example:

    keyboard.Listen(func(key keys.Key) (stop bool, err error) {
    	if key.Code == keys.CtrlC {
    		return true, nil // Stop listener by returning true on Ctrl+C
    	}

    	fmt.Println("\r" + key.String()) // Print every key press
    	return false, nil // Return false to continue listening
    })

#### func  SimulateKeyPress

```go
func SimulateKeyPress(input ...interface{}) error
```
SimulateKeyPress simulate a key press. It can be used to mock user input and
test your application.

Example:

    go func() {
    	keyboard.SimulateKeyPress("Hello")             // Simulate key press for every letter in string
    	keyboard.SimulateKeyPress(keys.Enter)          // Simulate key press for Enter
    	keyboard.SimulateKeyPress(keys.CtrlShiftRight) // Simulate key press for Ctrl+Shift+Right
    	keyboard.SimulateKeyPress('x')                 // Simulate key press for a single rune
    	keyboard.SimulateKeyPress('x', keys.Down, 'a') // Simulate key presses for multiple inputs
    }()

---

> [AtomicGo.dev](https://atomicgo.dev) &nbsp;&middot;&nbsp;
> with ❤️ by [@MarvinJWendt](https://github.com/MarvinJWendt) |
> [MarvinJWendt.com](https://marvinjwendt.com)
