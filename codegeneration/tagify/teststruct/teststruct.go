// Package teststruct is just spectacular.
//
/*
          ,_---~~~~~----._
   _,,_,*^____      _____``*g*\"*,
  / __/ /'     ^.  /      \ ^@q   f
 |  @f | @))    |  | @))   l  0 _/
  \`/   \~____ / __ \_____/    \
   |           _l__l_           I
   }          [______]           I
   ]            | | |            |
   ]             ~ ~             |
   |                            |
    |                           |
*/
//
package teststruct

//go:generate godocdown -output README.md

type T struct {
	Name        string
	ID          int
	privatedata map[string]interface{}
}
