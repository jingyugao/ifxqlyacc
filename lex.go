package ifxqlyacc

import (
	"vitess.io/vitess/go/vt/log"
	"strconv"
	"time"
	"fmt"
	"github.com/pkg/errors"
)

type Tokenizer struct {
	query Query
	scanner *Scanner
}

func (tkn *Tokenizer) Lex(lval *yySymType) int{
	var typ int
	var val string

	for {
		typ, _, val  = tkn.scanner.Scan()
		if typ == EOF{
			return 0
		}
		if typ == MUL{
			lval.int = ILLEGAL
			//break
		}
		if typ >= EQ && typ <= GTE{
			//println("oo string is ",val)
			//lval.int,_ = strconv.Atoi(val)
			//println("oo is ",lval.int)
			lval.int = typ
			//break
		}
		if typ == NUMBER{
			lval.float64, _ = strconv.ParseFloat(val,64)
			//break
		}
		if typ == INTEGER{
			lval.int64, _ = strconv.ParseInt(val,10,64)
			//break
		}
		if typ == DURATIONVAL{
			time,_ := ParseDuration(val)
			lval.tdur = time
			//break
		}
		if typ == DESC{
			lval.bool = false
			//break
		}
		if typ == AND{
			lval.int = AND
			//break
		}
		if typ == OR{
			lval.int = OR
			//break
		}
		if typ == ASC{
			lval.bool = true
			//break
		}
		if typ !=WS{
			break
		}
	}
	lval.str = val
	return typ
}
func (tkn *Tokenizer) Error(err string){
	log.Fatal(err)
}

var ErrInvalidDuration = errors.New("invalid duration")
// ParseDuration parses a time duration from a string.
// This is needed instead of time.ParseDuration because this will support
// the full syntax that InfluxQL supports for specifying durations
// including weeks and days.
func ParseDuration(s string) (time.Duration, error) {
	// Return an error if the string is blank or one character
	if len(s) < 2 {
		return 0, ErrInvalidDuration
	}

	// Split string into individual runes.
	a := []rune(s)

	// Start with a zero duration.
	var d time.Duration
	i := 0

	// Check for a negative.
	isNegative := false
	if a[i] == '-' {
		isNegative = true
		i++
	}

	var measure int64
	var unit string

	// Parsing loop.
	for i < len(a) {
		// Find the number portion.
		start := i
		for ; i < len(a) && isDigit(a[i]); i++ {
			// Scan for the digits.
		}

		// Check if we reached the end of the string prematurely.
		if i >= len(a) || i == start {
			return 0, ErrInvalidDuration
		}

		// Parse the numeric part.
		n, err := strconv.ParseInt(string(a[start:i]), 10, 64)
		if err != nil {
			return 0, ErrInvalidDuration
		}
		measure = n

		// Extract the unit of measure.
		// If the last two characters are "ms" then parse as milliseconds.
		// Otherwise just use the last character as the unit of measure.
		unit = string(a[i])
		switch a[i] {
		case 'n':
			if i+1 < len(a) && a[i+1] == 's' {
				unit = string(a[i : i+2])
				d += time.Duration(n)
				i += 2
				continue
			}
			return 0, ErrInvalidDuration
		case 'u', 'µ':
			d += time.Duration(n) * time.Microsecond
		case 'm':
			if i+1 < len(a) && a[i+1] == 's' {
				unit = string(a[i : i+2])
				d += time.Duration(n) * time.Millisecond
				i += 2
				continue
			}
			d += time.Duration(n) * time.Minute
		case 's':
			d += time.Duration(n) * time.Second
		case 'h':
			d += time.Duration(n) * time.Hour
		case 'd':
			d += time.Duration(n) * 24 * time.Hour
		case 'w':
			d += time.Duration(n) * 7 * 24 * time.Hour
		default:
			return 0, ErrInvalidDuration
		}
		i++
	}

	// Check to see if we overflowed a duration
	if d < 0 && !isNegative {
		return 0, fmt.Errorf("overflowed duration %d%s: choose a smaller duration or INF", measure, unit)
	}

	if isNegative {
		d = -d
	}
	return d, nil
}