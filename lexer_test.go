package abc

import "testing"

func TestLexer(t *testing.T) {
	cases := []struct {
		in  string
		out []item
	}{
		{"", []item{}},
		{"T:title", []item{
			itemField{'T', "title"},
		}},
		{"X:ref\nT:title", []item{
			itemField{'X', "ref"},
			itemField{'T', "title"},
		}},
		{"A", []item{
			itemNote{note: 'A'},
		}},
		{" A", []item{
			itemNote{note: 'A'},
		}},
		{"AaCcdE", []item{
			itemNote{note: 'A'},
			itemNote{note: 'a'},
			itemNote{note: 'C'},
			itemNote{note: 'c'},
			itemNote{note: 'd'},
			itemNote{note: 'E'},
		}},
		{"_A", []item{
			itemNote{note: 'A', acc: "_"},
		}},
		{"=A", []item{
			itemNote{note: 'A', acc: "="},
		}},
		{"__A", []item{
			itemNote{note: 'A', acc: "__"},
		}},
		{"^A", []item{
			itemNote{note: 'A', acc: "^"},
		}},
		{"^^A", []item{
			itemNote{note: 'A', acc: "^^"},
		}},
		{"=A", []item{
			itemNote{note: 'A', acc: "="},
		}},
		{"^_A", []item{
			itemError{msg: "wrong accidental", ctx: "^_"},
		}},
		{"A'", []item{
			itemNote{note: 'A', oct: "'"},
		}},
		{"A,", []item{
			itemNote{note: 'A', oct: ","},
		}},
		{"A,,", []item{
			itemNote{note: 'A', oct: ",,"},
		}},
		{"A,,A''", []item{
			itemNote{note: 'A', oct: ",,"},
			itemNote{note: 'A', oct: "''"},
		}},
		{"a'_b','^C,,", []item{
			itemNote{note: 'a', oct: "'"},
			itemNote{note: 'b', acc: "_", oct: "','"},
			itemNote{note: 'C', acc: "^", oct: ",,"},
		}},
		{"G4G4G4_E3_B1G4_E3_B1G4", []item{
			itemNote{note: 'G', ryt: "4"},
			itemNote{note: 'G', ryt: "4"},
			itemNote{note: 'G', ryt: "4"},
			itemNote{note: 'E', acc: "_", ryt: "3"},
			itemNote{note: 'B', acc: "_", ryt: "1"},
			itemNote{note: 'G', ryt: "4"},
			itemNote{note: 'E', acc: "_", ryt: "3"},
			itemNote{note: 'B', acc: "_", ryt: "1"},
			itemNote{note: 'G', ryt: "4"},
		}},
		{"z/8 [c'/4a/4f/4e/4] z/8 [f/8D,/2] z/8", []item{
			itemNote{note: 'z', ryt: "/8"},
			itemNote{note: 'c', oct: "'", ryt: "/4"},
			itemNote{note: 'a', ryt: "/4"},
			itemNote{note: 'f', ryt: "/4"},
			itemNote{note: 'e', ryt: "/4"},
			itemNote{note: 'z', ryt: "/8"},
			itemNote{note: 'f', ryt: "/8"},
			itemNote{note: 'D', oct: ",", ryt: "/2"},
			itemNote{note: 'z', ryt: "/8"},
		}},
	}

	for _, c := range cases {
		_, items := lex(c.in)
		for _, exp := range c.out {
			got, ok := <-items
			if !ok {
				t.Errorf("expected %v, channel was close\n", exp)
				break
			}
			if got != exp {
				t.Errorf("expected %v, got %v\n", exp, got)
			}
		}
		for {
			got, ok := <-items
			if !ok {
				break
			}
			t.Errorf("unexpected %v", got)
		}
	}
}
