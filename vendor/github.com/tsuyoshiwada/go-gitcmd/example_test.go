package gitcmd

import (
	"fmt"
	"log"
)

func Example() {
	git := New(nil)

	out, err := git.Exec("rev-parse", "--git-dir")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(out)
	// Output:
	// .git
}
