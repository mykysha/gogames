package interceptor

import (
	"fmt"
	"os"

	"golang.org/x/sys/unix"
)

func InterceptKeystrokes(keys chan byte) error {
	oldState, err := enableRawMode()
	if err != nil {
		return fmt.Errorf("failed to enable raw mode: %w", err)
	}

	defer func() {
		if err := disableRawMode(oldState); err != nil {
			panic(err)
		}
	}()

	for {
		key := make([]byte, 1)

		if _, err := os.Stdin.Read(key); err != nil {
			return fmt.Errorf("failed to read from stdin: %w", err)
		}

		keys <- key[0]
	}
}

func enableRawMode() (*unix.Termios, error) {
	fileDescriptor := int(os.Stdin.Fd())

	oldState, err := unix.IoctlGetTermios(fileDescriptor, unix.TIOCGETA)
	if err != nil {
		return nil, fmt.Errorf("failed to get old terminal state: %w", err)
	}

	newState := *oldState
	newState.Lflag &^= unix.ECHO   // turn off echo
	newState.Lflag &^= unix.ICANON // disable canonical mode
	newState.Cc[unix.VMIN] = 1
	newState.Cc[unix.VTIME] = 0

	if err = unix.IoctlSetTermios(fileDescriptor, unix.TIOCSETA, &newState); err != nil {
		return nil, fmt.Errorf("failed to set new terminal state: %w", err)
	}

	return oldState, nil
}

func disableRawMode(state *unix.Termios) error {
	fd := int(os.Stdin.Fd())

	if err := unix.IoctlSetTermios(fd, unix.TIOCSETA, state); err != nil {
		return fmt.Errorf("failed to restore old terminal state: %w", err)
	}

	return nil
}
