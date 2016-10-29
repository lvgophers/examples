package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
)

var overwrite = flag.Bool("w", false, "Overwrite existing file (default stdout) (not recommended)")
var oe = flag.Bool("oe", false, "Add omitempty")

func main() {
	flag.Parse()
	srcfile := flag.Arg(0)
	if srcfile == "" {
		flag.Usage()
		log.Fatal("Source file required")
	}
	src, err := ioutil.ReadFile(srcfile)
	if err != nil {
		log.Fatal(err)
	}
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		log.Fatal(err)
	}
	omitempty := ""
	if *oe {
		omitempty = ",omitempty"
	}
	for _, decl := range f.Decls {
		if td, ok := decl.(*ast.GenDecl); ok {
			for _, spec := range td.Specs {
				if ts, ok := spec.(*ast.TypeSpec); ok {
					if st, ok := ts.Type.(*ast.StructType); ok {
						r, _ := utf8.DecodeRuneInString(ts.Name.Name)
						if !unicode.IsUpper(r) {
							continue
						}
						// START OMIT
						for _, f := range st.Fields.List {
							if len(f.Names) == 1 {
								r, _ = utf8.DecodeRuneInString(f.Names[0].Name)
								if unicode.IsUpper(r) {
									f.Tag = &ast.BasicLit{Kind: token.STRING,
										Value: fmt.Sprintf("`json:\"%s%s\"`",
											strings.ToLower(f.Names[0].Name), omitempty)}
								}
							}
						}
						// END OMIT
					}
				}
			}
		}
	}
	out := os.Stdout
	if *overwrite {
		out, err = os.Create(srcfile)
		if err != nil {
			log.Fatal(err)
		}
	}
	format.Node(out, fset, f)
}
