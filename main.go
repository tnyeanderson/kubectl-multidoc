package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

var (
	multidocSep = []byte("---\n")
	lineSep     = []byte("\n")[0]
	arrSep      = "- "
	indent      = "  "
)

// splitToMultidoc reads a Kubernetes API response in YAML format (usually from
// "kubectl get <resource> -oyaml"), and outputs it as a YAML multidoc with
// each member of the "items" array as its own document.
//
// Importantly, this function does not understand YAML. It does not attempt to
// load the YAML into memory, or to ensure it is valid. It checks for the
// beginning of the "items" array definition, then goes line by line and
// changes any array start token ("- ") to a multidoc separator, and unindents
// the lines by one level (two spaces).
//
// Therefore, this function is highly dependent on the formatting of its input.
// It expects:
//
//   - Two space indentation
//   - Non-indented array start tokens (e.g. the hyphen starting an array item
//     should be at the same level as its parent--in the case of the "items"
//     array, not indented at all)
//
// The output of this function is not guaranteed to be valid YAML.
func splitToMultidoc(r io.Reader, w io.Writer) error {
	// If this was any more complicated, I'd build a lexer
	br := bufio.NewReader(r)
	inItems := false
	for {
		line, err := br.ReadBytes(lineSep)
		if !inItems {
			if string(line) == "items:\n" {
				inItems = true
			}
		} else {
			if len(line) > 2 {
				switch string(line[:2]) {
				case arrSep:
					w.Write(multidocSep)
					w.Write(line[2:])
				case indent:
					w.Write(line[2:])
				default:
					// We have left the items array
					break
				}
			}
		}
		if err == io.EOF {
			if !inItems {
				return fmt.Errorf("kubernetes list response does not appear valid")
			}
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	if err := splitToMultidoc(os.Stdin, os.Stdout); err != nil {
		log.Fatal(err)
	}
}
