package cli

import (
	"os"
	"testing"
)

func TestCompletion_AllShells(t *testing.T) {
	shells := []string{"bash", "zsh", "fish", "powershell"}

	for _, shell := range shells {
		t.Run(shell, func(t *testing.T) {
			// Capture stdout
			oldStdout := os.Stdout
			r, w, err := os.Pipe()
			if err != nil {
				t.Fatalf("failed to create pipe: %v", err)
			}
			os.Stdout = w

			err = runCompletion(completionCmd, []string{shell})

			w.Close()
			os.Stdout = oldStdout

			if err != nil {
				t.Fatalf("runCompletion(%q) returned error: %v", shell, err)
			}

			buf := make([]byte, 1024)
			n, _ := r.Read(buf)
			r.Close()

			if n == 0 {
				t.Errorf("runCompletion(%q) produced no output", shell)
			}
		})
	}
}
