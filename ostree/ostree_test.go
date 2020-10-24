package ostree

import (
	"testing"
)

func TestTree(t *testing.T) {
	actual := Print("testdata")
	expect := `testdata
├── a
│   └── b
│       └── b.txt
└── c
    └── c.txt
`
	if actual != expect {
		t.Errorf("\nactual\n%s\n != expect\n%s\n", actual, expect)
	}
}
