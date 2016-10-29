package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

var output = flag.String("o", "", "Output file (default stdout)")
var floats = flag.Bool("floats", false, "Do not use json.Number for numbers")

var out = os.Stdout
var input *os.File

func fprint(i ...interface{}) {
	if _, err := fmt.Fprintln(out, i...); err != nil {
		panic(err)
	}
}

func printf(f string, i ...interface{}) {
	if _, err := fmt.Fprintf(out, f, i...); err != nil {
		panic(err)
	}
}

func init() {
	flag.Parse()
	infile := flag.Arg(0)
	if infile == "" {
		flag.Usage()
		log.Fatal("Input file required")
	}
	var err error
	input, err = os.Open(infile)
	if err != nil {
		log.Fatal(err)
	}
	if outfile := *output; outfile != "" {
		f, err := os.Create(outfile)
		if err != nil {
			log.Fatal(err)
		}
		out = f
	}
}

type structify struct {
	*json.Decoder
}

func (s *structify) findend() {
	n := 1
	for {
		tok, err := s.Token()
		if err != nil {
			panic(err)
		}
		if d, ok := tok.(json.Delim); ok {
			switch d {
			case '[':
				n++
			case ']':
				n--
			}
		}
		if n == 0 {
			return
		}
	}
}

func (s *structify) st(tok json.Token) (err error) {
	if d, ok := tok.(json.Delim); ok {
		if d != '{' {
			return fmt.Errorf("struct must start with {")
		}
		fprint("struct {")
		for {
			// START OMIT
			tok, err = s.Token()
			if err != nil {
				return
			}
			switch t := tok.(type) {
			case string:
				printf("%s ", strings.Title(t))
				tag := fmt.Sprintf("`json:\"%s,omitempty\"` ", t)
				tok, err = s.Token()
				if err != nil {
					return
				}
				// END OMIT
				switch t := tok.(type) {
				case nil:
					printf("interface{} %s // TODO: json null\n", tag)
				case json.Delim:
					switch t {
					case '{':
						s.st(t)
						fprint(tag)
					case '[':
						printf("[]interface{} %s // TODO: json array \n", tag)
						s.findend()
					}
				default:
					printf("%T %s\n", t, tag)
				}
			case json.Delim:
				switch t {
				case '}':
					fmt.Print("} ")
					return
				default:
					return fmt.Errorf("unexpected delimiter: %s", t)
				}
			default:
				return fmt.Errorf("unexpected type: %T", t)
			}
		}
	}
	return
}

func main() {
	defer out.Close()
	defer input.Close()
	dec := json.NewDecoder(input)
	if !*floats {
		dec.UseNumber()
	}
	t, err := dec.Token()
	if err != nil {
		log.Fatal(err)
	}
	st := &structify{dec}
	err = st.st(t)
	if err != nil {
		log.Fatal(err)
	}
	fprint("")
}
