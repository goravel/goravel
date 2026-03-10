package keyboard

import (
	"fmt"
	"os"

	"github.com/containerd/console"

	"atomicgo.dev/keyboard/keys"
)

var windowsStdin *os.File
var con console.Console
var stdin = os.Stdin
var inputTTY *os.File
var mockChannel = make(chan keys.Key)

var mocking = false

func startListener() error {
	err := initInput()
	if err != nil {
		return err
	}

	if mocking {
		return nil
	}

	if con != nil {
		err := con.SetRaw()
		if err != nil {
			return fmt.Errorf("failed to set raw mode: %w", err)
		}
	}

	inputTTY, err = openInputTTY()
	if err != nil {
		return err
	}

	return nil
}

func stopListener() error {
	if con != nil {
		err := con.Reset()
		if err != nil {

			return fmt.Errorf("failed to reset console: %w", err)
		}
	}

	return restoreInput()
}

// Listen calls a callback function when a key is pressed.
//
// Simple example:
//
//	keyboard.Listen(func(key keys.Key) (stop bool, err error) {
//		if key.Code == keys.CtrlC {
//			return true, nil // Stop listener by returning true on Ctrl+C
//		}
//
//		fmt.Println("\r" + key.String()) // Print every key press
//		return false, nil // Return false to continue listening
//	})
func Listen(onKeyPress func(key keys.Key) (stop bool, err error)) error {
	cancel := make(chan bool)
	stopRoutine := false

	go func() {
		for {
			select {
			case c := <-cancel:
				if c {
					return
				}
			case keyInfo := <-mockChannel:
				stopRoutine, _ = onKeyPress(keyInfo)
				if stopRoutine {
					closeInput()
					inputTTY.Close()
				}
			}
		}
	}()

	err := startListener()
	if err != nil {
		if err.Error() != "provided file is not a console" {
			return err
		}
	}

	for !stopRoutine {
		key, err := getKeyPress()
		if err != nil {
			return err
		}

		// check if returned key is empty
		// if reflect.DeepEqual(key, keys.Key{}) {
		// 	return nil
		// }

		stop, err := onKeyPress(key)
		if err != nil {
			return err
		}

		if stop {
			closeInput()
			inputTTY.Close()
			break
		}
	}

	err = stopListener()
	if err != nil {
		return err
	}

	cancel <- true

	return nil
}

// SimulateKeyPress simulate a key press. It can be used to mock user stdin and test your application.
//
// Example:
//
//	go func() {
//		keyboard.SimulateKeyPress("Hello")             // Simulate key press for every letter in string
//		keyboard.SimulateKeyPress(keys.Enter)          // Simulate key press for Enter
//		keyboard.SimulateKeyPress(keys.CtrlShiftRight) // Simulate key press for Ctrl+Shift+Right
//		keyboard.SimulateKeyPress('x')                 // Simulate key press for a single rune
//		keyboard.SimulateKeyPress('x', keys.Down, 'a') // Simulate key presses for multiple inputs
//	}()
func SimulateKeyPress(input ...interface{}) error {
	for _, key := range input {
		// Check if key is a keys.Key
		if key, ok := key.(keys.Key); ok {
			mockChannel <- key
			return nil
		}

		// Check if key is a rune
		if key, ok := key.(rune); ok {
			mockChannel <- keys.Key{
				Code:  keys.RuneKey,
				Runes: []rune{key},
			}
			return nil
		}

		// Check if key is a string
		if key, ok := key.(string); ok {
			for _, r := range key {
				mockChannel <- keys.Key{
					Code:  keys.RuneKey,
					Runes: []rune{r},
				}
			}
			return nil
		}

		// Check if key is a KeyCode
		if key, ok := key.(keys.KeyCode); ok {
			mockChannel <- keys.Key{
				Code: key,
			}
			return nil
		}
	}

	return nil
}
