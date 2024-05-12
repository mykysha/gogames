package interceptor

import (
	"os"

	"golang.org/x/sys/unix"
)

func InterceptKeystrokes(keys chan rune) {
	oldState, err := enableRawMode()
	if err != nil {
		panic(err) // TODO
	}

	defer disableRawMode(oldState)

	for {
		b := make([]byte, 1)

		if _, err := os.Stdin.Read(b); err != nil {
			panic(err) // TODO
		}

		if b[0] == 'q' {
			panic("exit") // TODO
		}

		keys <- rune(b[0])
	}
}

func enableRawMode() (*unix.Termios, error) {
	fd := int(os.Stdin.Fd())
	oldState, err := unix.IoctlGetTermios(fd, unix.TIOCGETA)
	if err != nil {
		return nil, err
	}

	newState := *oldState
	newState.Lflag &^= unix.ECHO   // turn off echo
	newState.Lflag &^= unix.ICANON // disable canonical mode
	newState.Cc[unix.VMIN] = 1
	newState.Cc[unix.VTIME] = 0

	err = unix.IoctlSetTermios(fd, unix.TIOCSETA, &newState)
	if err != nil {
		return nil, err
	}

	return oldState, nil
}

func disableRawMode(state *unix.Termios) error {
	fd := int(os.Stdin.Fd())
	return unix.IoctlSetTermios(fd, unix.TIOCSETA, state)
}
