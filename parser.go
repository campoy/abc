package abc

import (
	"fmt"
	"strconv"
	"strings"
)

type Tune struct {
	Title           string
	key             string // TODO: parse and use key
	timesPerMeasure uint8
	lengthPerTime   uint8
	Notes           []Note
}

type Note struct {
	Oct    int8  // octave number
	Name   byte  // A to G
	Acc    int8  // -2 to +2
	Length uint8 // duration in 1/64 of quarter
}

var accNames = map[int8]string{
	-2: "ß",
	-1: "b",
	0:  "",
	1:  "#",
	2:  "x",
}

func (n Note) String() string {
	return fmt.Sprintf("%v%c%v[%v]",
		accNames[n.Acc],
		n.Name,
		n.Oct,
		n.Length,
	)
}

func Parse(input string) (*Tune, error) {
	t := &Tune{lengthPerTime: 8}

	l, items := lex(input)
	defer l.stop()

	for it := range items {
		switch it := it.(type) {
		case itemError:
			return nil, fmt.Errorf("at %s: %v", it.ctx, it.msg)
		case itemField:
			switch it.key {
			case 'T':
				t.Title = it.val
			case 'K':
				t.key = it.val
			case 'L':

			}
		case itemNote:
			n, err := parseNote(t, it)
			if err != nil {
				return nil, err
			}
			t.Notes = append(t.Notes, n)
		}
	}

	return t, nil
}

func parseNote(tune *Tune, in itemNote) (Note, error) {
	var note Note

	name, oct := in.note, 3
	if 'a' <= name && name <= 'g' {
		name = 'A' + name - 'a'
		oct++
	}
	note.Name = name

	for _, acc := range in.acc {
		switch acc {
		case '\'':
			oct++
		case ',':
			oct--
		}
	}
	note.Oct = int8(oct)

	switch in.acc {
	case "__":
		note.Acc = -2
	case "_":
		note.Acc = -1
	case "=":
		note.Acc = 0
	case "^":
		note.Acc = 1
	case "^^":
		note.Acc = 2
	case "":
		// TODO: depends on the key
	}

	parts := strings.Split(in.ryt, "/")
	if len(parts) > 2 {
		return note, fmt.Errorf("wrong duration %s", in.ryt)
	}

	num := 1
	if len(parts[0]) > 0 {
		n, err := strconv.Atoi(parts[0])
		if err != nil {
			return note, fmt.Errorf("wrong duration numerator %s", in.ryt)
		}
		num = n
	}

	den := int(tune.lengthPerTime)
	if len(parts) == 2 && len(parts[1]) > 0 {
		n, err := strconv.Atoi(parts[1])
		if err != nil {
			return note, fmt.Errorf("wrong duration denominator %s", in.ryt)
		}
		den = n
	}

	note.Length = uint8((64 * num) / den)
	return note, nil
}
