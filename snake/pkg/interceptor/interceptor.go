package interceptor

import (
	"fmt"
	"os"

	"golang.org/x/sys/unix"
)

func InterceptKeystrokes(keys chan rune) error {
	oldState, err := enableRawMode()
	if err != nil {
		return fmt.Errorf("failed to enable raw mode: %w", err)
	}

	defer disableRawMode(oldState)

	for {
		b := make([]byte, 1)

		if _, err := os.Stdin.Read(b); err != nil {
			return fmt.Errorf("failed to read from stdin: %w", err)
		}

		keys <- rune(b[0])
	}
}

func enableRawMode() (*unix.Termios, error) {
	fd := int(os.Stdin.Fd())

	oldState, err := unix.IoctlGetTermios(fd, unix.TIOCGETA)
	if err != nil {
		return nil, fmt.Errorf("failed to get old terminal state: %w", err)
	}

	newState := *oldState
	newState.Lflag &^= unix.ECHO   // turn off echo
	newState.Lflag &^= unix.ICANON // disable canonical mode
	newState.Cc[unix.VMIN] = 1
	newState.Cc[unix.VTIME] = 0

	if err = unix.IoctlSetTermios(fd, unix.TIOCSETA, &newState); err != nil {
		return nil, fmt.Errorf("failed to set new terminal state: %w", err)
	}

	return oldState, nil
}

func disableRawMode(state *unix.Termios) error {
	fd := int(os.Stdin.Fd())

	return unix.IoctlSetTermios(fd, unix.TIOCSETA, state)
}
