package log

import (
	"os"
	"testing"
)

func TestSetLevel(t *testing.T) {

	setLevel(InfoLevel)
	if infoLog.Writer() != os.Stdout || errorLog.Writer() != os.Stdout {
		t.Fatal("failed to set log level")
	}

	setLevel(ErrorLevel)
	if infoLog.Writer() == os.Stdout || errorLog.Writer() != os.Stdout {
		t.Fatal("failed to set log level")
	}

	setLevel(Disabled)
	if infoLog.Writer() == os.Stdout || errorLog.Writer() == os.Stdout {
		t.Fatal("failed to set log level")
	}

}
