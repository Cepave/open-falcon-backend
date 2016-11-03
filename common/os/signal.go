package os

import (
	"fmt"
	"os"
	"os/signal"
)

// Callback with signal callback
type ExitCallback func(signal os.Signal)

func HoldingAndWaitSignal(exitCallback ExitCallback, signals ...os.Signal) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, signals...)

	go func() {
		signal := <-sigs

		fmt.Println()

		exitCallback(signal)

		os.Exit(0)
	} ()

	select {}
}
