package format_test

import (
	"context"
	"fmt"
	"strings"
	"time"

	. "github.com/onsi/ginkgo/v2"

	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/format"
	"github.com/onsi/gomega/types"
)

//recursive struct

const truncateHelpText = `
Gomega truncated this representation as it exceeds 'format.MaxLength'.
Consider having the object provide a custom 'GomegaStringer' representation
or adjust the parameters in Gomega's 'format' package.

Learn more here: https://onsi.github.io/gomega/#adjusting-output
`

type StringAlias string
type ByteAlias []byte
type IntAlias int

type AStruct struct {
	Exported string
}

type SimpleStruct struct {
	Name        string
	Enumeration int
	Veritas     bool
	Data        []byte
	secret      uint32
}

type ComplexStruct struct {
	Strings      []string
	SimpleThings []*SimpleStruct
	DataMaps     map[int]ByteAlias
}

type SecretiveStruct struct {
	boolValue      bool
	intValue       int
	uintValue      uint
	uintptrValue   uintptr
	floatValue     float32
	complexValue   complex64
	chanValue      chan bool
	funcValue      func()
	pointerValue   *int
	sliceValue     []string
	byteSliceValue []byte
	stringValue    string
	arrValue       [3]int
	byteArrValue   [3]byte
	mapValue       map[string]int
	structValue    AStruct
	interfaceValue interface{}
}

type CustomFormatted struct {
	Data  string
	Count int
}
type NotCustomFormatted struct {
	Data  string
	Count int
}

type CustomError struct {
	Details string
}

var _ error = &CustomError{}

func (c *CustomError) Error() string {
	return c.Details
}

func customFormatter(obj interface{}) (string, bool) {
	cf, ok := obj.(CustomFormatted)
	if !ok {
		return "", false
	}
	return fmt.Sprintf("%s (%d)", cf.Data, cf.Count), true
}

type GoStringer struct {
}

func (g GoStringer) GoString() string {
	return "go-string"
}

func (g GoStringer) String() string {
	return "string"
}

type Stringer struct {
}

func (g Stringer) String() string {
	return "string"
}

type gomegaStringer struct {
}

func (g gomegaStringer) GomegaString() string {
	return "gomegastring"
}

type gomegaStringerLong struct {
}

func (g gomegaStringerLong) GomegaString() string {
	return strings.Repeat("s", MaxLength*2)
}

type gomegaStringerMultiline struct {
}

func (g gomegaStringerMultiline) GomegaString() string {
	return "A\nB\nC"
}

var _ = Describe("Format", func() {
	match := func(typeRepresentation string, valueRepresentation string, args ...interface{}) types.GomegaMatcher {
		if len(args) > 0 {
			valueRepresentation = fmt.Sprintf(valueRepresentation, args...)
		}
		return Equal(fmt.Sprintf("%s<%s>: %s", Indent, typeRepresentation, valueRepresentation))
	}

	matchRegexp := func(typeRepresentation string, valueRepresentation string, args ...interface{}) types.GomegaMatcher {
		if len(args) > 0 {
			valueRepresentation = fmt.Sprintf(valueRepresentation, args...)
		}
		return MatchRegexp(fmt.Sprintf("%s<%s>: %s", Indent, typeRepresentation, valueRepresentation))
	}

	hashMatchingRegexp := func(entries ...string) string {
		entriesSwitch := "(" + strings.Join(entries, "|") + ")"
		arr := make([]string, len(entries))
		for i := range arr {
			arr[i] = entriesSwitch
		}
		return "{\\s*" + strings.Join(arr, ",\\s* ") + ",?\\s*}"
	}

	Describe("Message", func() {
		Context("with only an actual value", func() {
			BeforeEach(func() {
				MaxLength = 4000
			})

			It("should print out an indented formatted representation of the value and the message", func() {
				Expect(Message(3, "to be three.")).Should(Equal("Expected\n    <int>: 3\nto be three."))
			})

			It("should print out an indented formatted representation of the value and the message, and trucate it when too long", func() {
				tooLong := strings.Repeat("s", MaxLength+1)
				tooLongResult := strings.Repeat("s", MaxLength) + "...\n" + truncateHelpText
				Expect(Message(tooLong, "to be truncated")).Should(Equal("Expected\n    <string>: " + tooLongResult + "\nto be truncated"))
			})

			It("should print out an indented formatted representation of the value and the message, and not trucate it when MaxLength = 0", func() {
				MaxLength = 0
				tooLong := strings.Repeat("s", MaxLength+1)
				Expect(Message(tooLong, "to be truncated")).Should(Equal("Expected\n    <string>: " + tooLong + "\nto be truncated"))
			})
		})

		Context("with an actual and an expected value", func() {
			It("should print out an indented formatted representatino of both values, and the message", func() {
				Expect(Message(3, "to equal", 4)).Should(Equal("Expected\n    <int>: 3\nto equal\n    <int>: 4"))
			})
		})
	})

	Describe("MessageWithDiff", func() {
		It("shows the exact point where two long strings differ", func() {
			stringWithB := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaabaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
			stringWithZ := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaazaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

			Expect(MessageWithDiff(stringWithB, "to equal", stringWithZ)).Should(Equal(expectedLongStringFailureMessage))
		})

		It("truncates the start of long strings that differ only at their end", func() {
			stringWithB := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaab"
			stringWithZ := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaz"

			Expect(MessageWithDiff(stringWithB, "to equal", stringWithZ)).Should(Equal(expectedTruncatedStartStringFailureMessage))
		})

		It("truncates the start of long strings that differ only in length", func() {
			smallString := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
			largeString := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

			Expect(MessageWithDiff(largeString, "to equal", smallString)).Should(Equal(expectedTruncatedStartSizeFailureMessage))
			Expect(MessageWithDiff(smallString, "to equal", largeString)).Should(Equal(expectedTruncatedStartSizeSwappedFailureMessage))
		})

		It("truncates the end of long strings that differ only at their start", func() {
			stringWithB := "baaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
			stringWithZ := "zaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

			Expect(MessageWithDiff(stringWithB, "to equal", stringWithZ)).Should(Equal(expectedTruncatedEndStringFailureMessage))
		})

		It("handles multi-byte sequences correctly", func() {
			stringA := "• abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz1"
			stringB := "• abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz"

			Expect(MessageWithDiff(stringA, "to equal", stringB)).Should(Equal(expectedTruncatedMultiByteFailureMessage))
		})

		It("prints special characters", func() {
			stringA := "\n"
			stringB := "something_else"

			Expect(MessageWithDiff(stringA, "to equal", stringB)).Should(Equal(expectedSpecialCharacterFailureMessage))
		})

		It("handles negative padding length", func() {
			stringWithB := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaabaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
			stringWithZ := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaazaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
			longMessage := "to equal very long message"

			Expect(MessageWithDiff(stringWithB, longMessage, stringWithZ)).Should(Equal(expectedDiffLongMessage))
		})

		Context("With truncated diff disabled", func() {
			BeforeEach(func() {
				TruncatedDiff = false
			})

			AfterEach(func() {
				TruncatedDiff = true
			})

			It("should show the full diff", func() {
				stringWithB := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaabaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
				stringWithZ := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaazaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

				Expect(MessageWithDiff(stringWithB, "to equal", stringWithZ)).Should(Equal(expectedFullFailureDiff))
			})
		})

		Context("With alternate diff lengths", func() {
			initialValue := TruncateThreshold // 50 by default
			BeforeEach(func() {
				TruncateThreshold = 10000
			})

			AfterEach(func() {
				TruncateThreshold = initialValue
			})

			It("should show the full diff when truncate threshold is increased beyond length of strings", func() {
				stringWithB := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaabaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
				stringWithZ := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaazaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

				Expect(MessageWithDiff(stringWithB, "to equal", stringWithZ)).Should(Equal(expectedFullFailureDiff))
			})
		})

		Context("with alternative number of characters to include around mismatch", func() {
			initialValue := CharactersAroundMismatchToInclude // 5 by default
			BeforeEach(func() {
				CharactersAroundMismatchToInclude = 10
			})

			AfterEach(func() {
				CharactersAroundMismatchToInclude = initialValue
			})

			It("it shows more characters around a line length mismatch", func() {
				smallString := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
				largeString := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"

				Expect(MessageWithDiff(largeString, "to equal", smallString)).Should(Equal(expectedTruncatedStartSizeFailureMessageExtraDiff))
				Expect(MessageWithDiff(smallString, "to equal", largeString)).Should(Equal(expectedTruncatedStartSizeSwappedFailureMessageExtraDiff))
			})
		})

		Describe("At extremes of configurable values", func() {
			Context("with zero-length threshold", func() {
				initialValue := TruncateThreshold // 50 by default
				BeforeEach(func() {
					TruncateThreshold = 0
				})

				AfterEach(func() {
					TruncateThreshold = initialValue
				})

				It("should show the full diff when truncate threshold is increased beyond length of strings", func() {
					stringWithB := "aba"
					stringWithZ := "aza"
					Expect(MessageWithDiff(stringWithB, "to equal", stringWithZ)).Should(Equal(expectedDiffSmallThreshold))
				})
			})

			Context("with zero characters around mismatch", func() {
				initialValue := CharactersAroundMismatchToInclude // 5 by default
				BeforeEach(func() {
					CharactersAroundMismatchToInclude = 0
				})

				AfterEach(func() {
					CharactersAroundMismatchToInclude = initialValue
				})

				It("", func() {
					stringWithB := "aba"
					stringWithZ := "aza"
					Expect(MessageWithDiff(stringWithB, "to equal", stringWithZ)).Should(Equal(expectedDiffZeroMismatch))
				})
			})

			Context("with zero-length threshold and zero characters around mismatch", func() {
				initialCharactersAroundMismatch := CharactersAroundMismatchToInclude
				initialTruncateThreshold := TruncateThreshold
				BeforeEach(func() {
					CharactersAroundMismatchToInclude = 0
					TruncateThreshold = 0
				})

				AfterEach(func() {
					CharactersAroundMismatchToInclude = initialCharactersAroundMismatch
					TruncateThreshold = initialTruncateThreshold
				})

				It("", func() {
					stringWithB := "aba"
					stringWithZ := "aza"
					Expect(MessageWithDiff(stringWithB, "to equal", stringWithZ)).Should(Equal(expectedDiffSmallThresholdZeroMismatch))
				})
			})
		})
	})

	Describe("IndentString", func() {
		It("should indent the string", func() {
			Expect(IndentString("foo\n  bar\nbaz", 2)).Should(Equal("        foo\n          bar\n        baz"))
		})
	})

	Describe("Object", func() {
		Describe("formatting boolean values", func() {
			It("should give the type and format values correctly", func() {
				Expect(Object(true, 1)).Should(match("bool", "true"))
				Expect(Object(false, 1)).Should(match("bool", "false"))
			})
		})

		Describe("formatting numbers", func() {
			It("should give the type and format values correctly", func() {
				Expect(Object(int(3), 1)).Should(match("int", "3"))
				Expect(Object(int8(3), 1)).Should(match("int8", "3"))
				Expect(Object(int16(3), 1)).Should(match("int16", "3"))
				Expect(Object(int32(3), 1)).Should(match("int32", "3"))
				Expect(Object(int64(3), 1)).Should(match("int64", "3"))

				Expect(Object(uint(3), 1)).Should(match("uint", "3"))
				Expect(Object(uint8(3), 1)).Should(match("uint8", "3"))
				Expect(Object(uint16(3), 1)).Should(match("uint16", "3"))
				Expect(Object(uint32(3), 1)).Should(match("uint32", "3"))
				Expect(Object(uint64(3), 1)).Should(match("uint64", "3"))
			})

			It("should handle uintptr differently", func() {
				Expect(Object(uintptr(3), 1)).Should(match("uintptr", "0x3"))
			})
		})

		Describe("formatting channels", func() {
			It("should give the type and format values correctly", func() {
				c := make(chan<- bool, 3)
				c <- true
				c <- false
				Expect(Object(c, 1)).Should(match("chan<- bool | len:2, cap:3", "%v", c))
			})
		})

		Describe("formatting strings", func() {
			It("should give the type and format values correctly", func() {
				s := "a\nb\nc"
				Expect(Object(s, 1)).Should(match("string", `a
    b
    c`))
			})
		})

		Describe("formatting []byte slices", func() {
			When("the slice is made of printable bytes", func() {
				It("should present it as string", func() {
					b := []byte("a b c")
					Expect(Object(b, 1)).Should(matchRegexp(`\[\]uint8 \| len:5, cap:\d+`, `a b c`))
				})
			})
			When("the slice contains non-printable bytes", func() {
				It("should present it as slice", func() {
					b := []byte("a b c\n\x01\x02\x03\xff\x1bH")
					Expect(Object(b, 1)).Should(matchRegexp(`\[\]uint8 \| len:12, cap:\d+`, `\[97, 32, 98, 32, 99, 10, 1, 2, 3, 255, 27, 72\]`))
				})
			})
		})

		Describe("formatting functions", func() {
			It("should give the type and format values correctly", func() {
				f := func(a string, b []int) ([]byte, error) {
					return []byte("abc"), nil
				}
				Expect(Object(f, 1)).Should(match("func(string, []int) ([]uint8, error)", "%v", f))
			})
		})

		Describe("formatting pointers", func() {
			It("should give the type and dereference the value to format it correctly", func() {
				a := 3
				Expect(Object(&a, 1)).Should(match(fmt.Sprintf("*int | %p", &a), "3"))
			})

			When("there are pointers to pointers...", func() {
				It("should recursively deference the pointer until it gets to a value", func() {
					a := 3
					var b *int
					var c **int
					var d ***int
					b = &a
					c = &b
					d = &c

					Expect(Object(d, 1)).Should(match(fmt.Sprintf("***int | %p", d), "3"))
				})
			})

			When("the pointer points to nil", func() {
				It("should say nil and not explode", func() {
					var a *AStruct
					Expect(Object(a, 1)).Should(match("*format_test.AStruct | 0x0", "nil"))
				})
			})
		})

		Describe("formatting arrays", func() {
			It("should give the type and format values correctly", func() {
				w := [3]string{"Jed Bartlet", "Toby Ziegler", "CJ Cregg"}
				Expect(Object(w, 1)).Should(match("[3]string", `["Jed Bartlet", "Toby Ziegler", "CJ Cregg"]`))
			})

			Context("with byte arrays", func() {
				It("should give the type and format values correctly", func() {
					w := [3]byte{17, 28, 19}
					Expect(Object(w, 1)).Should(match("[3]uint8", `[17, 28, 19]`))
				})
			})
		})

		Describe("formatting slices", func() {
			It("should include the length and capacity in the type information", func() {
				s := make([]bool, 3, 4)
				Expect(Object(s, 1)).Should(match("[]bool | len:3, cap:4", "[false, false, false]"))
			})

			When("the slice contains long entries", func() {
				It("should format the entries with newlines", func() {
					w := []string{"Josiah Edward Bartlet", "Toby Ziegler", "CJ Cregg"}
					expected := `[
        "Josiah Edward Bartlet",
        "Toby Ziegler",
        "CJ Cregg",
    ]`
					Expect(Object(w, 1)).Should(match("[]string | len:3, cap:3", expected))
				})
			})
		})

		Describe("formatting maps", func() {
			It("should include the length in the type information", func() {
				m := make(map[int]bool, 5)
				m[3] = true
				m[4] = false
				Expect(Object(m, 1)).Should(matchRegexp(`map\[int\]bool \| len:2`, hashMatchingRegexp("3: true", "4: false")))
			})

			When("the slice contains long entries", func() {
				It("should format the entries with newlines", func() {
					m := map[string][]byte{}
					m["Josiah Edward Bartlet"] = []byte("Martin Sheen")
					m["Toby Ziegler"] = []byte("Richard Schiff")
					m["CJ Cregg"] = []byte("Allison Janney")
					expected := `{
        ("Josiah Edward Bartlet": "Martin Sheen"|"Toby Ziegler": "Richard Schiff"|"CJ Cregg": "Allison Janney"),
        ("Josiah Edward Bartlet": "Martin Sheen"|"Toby Ziegler": "Richard Schiff"|"CJ Cregg": "Allison Janney"),
        ("Josiah Edward Bartlet": "Martin Sheen"|"Toby Ziegler": "Richard Schiff"|"CJ Cregg": "Allison Janney"),
    }`
					Expect(Object(m, 1)).Should(matchRegexp(`map\[string\]\[\]uint8 \| len:3`, expected))
				})
			})
		})

		Describe("formatting structs", func() {
			It("should include the struct name and the field names", func() {
				s := SimpleStruct{
					Name:        "Oswald",
					Enumeration: 17,
					Veritas:     true,
					Data:        []byte("datum"),
					secret:      1983,
				}

				Expect(Object(s, 1)).Should(match("format_test.SimpleStruct", `{Name: "Oswald", Enumeration: 17, Veritas: true, Data: "datum", secret: 1983}`))
			})

			When("the struct contains long entries", func() {
				It("should format the entries with new lines", func() {
					s := &SimpleStruct{
						Name:        "Mithrandir Gandalf Greyhame",
						Enumeration: 2021,
						Veritas:     true,
						Data:        []byte("wizard"),
						secret:      3,
					}

					Expect(Object(s, 1)).Should(match(fmt.Sprintf("*format_test.SimpleStruct | %p", s), `{
        Name: "Mithrandir Gandalf Greyhame",
        Enumeration: 2021,
        Veritas: true,
        Data: "wizard",
        secret: 3,
    }`))
				})
			})
		})

		Describe("formatting nil values", func() {
			It("should print out nil", func() {
				Expect(Object(nil, 1)).Should(match("nil", "nil"))
				var typedNil *AStruct
				Expect(Object(typedNil, 1)).Should(match("*format_test.AStruct | 0x0", "nil"))
				var c chan<- bool
				Expect(Object(c, 1)).Should(match("chan<- bool | len:0, cap:0", "nil"))
				var s []string
				Expect(Object(s, 1)).Should(match("[]string | len:0, cap:0", "nil"))
				var m map[string]bool
				Expect(Object(m, 1)).Should(match("map[string]bool | len:0", "nil"))
			})
		})

		Describe("formatting aliased types", func() {
			It("should print out the correct alias type", func() {
				Expect(Object(StringAlias("alias"), 1)).Should(match("format_test.StringAlias", `alias`))
				Expect(Object(ByteAlias("alias"), 1)).Should(matchRegexp(`format_test\.ByteAlias \| len:5, cap:\d+`, `alias`))
				Expect(Object(IntAlias(3), 1)).Should(match("format_test.IntAlias", "3"))
			})
		})

		Describe("handling nested things", func() {
			It("should produce a correctly nested representation", func() {
				s := ComplexStruct{
					Strings: []string{"lots", "of", "short", "strings"},
					SimpleThings: []*SimpleStruct{
						{"short", 7, true, []byte("succinct"), 17},
						{"something longer", 427, true, []byte("designed to wrap around nicely"), 30},
					},
					DataMaps: map[int]ByteAlias{
						17:   ByteAlias("some substantially longer chunks of data"),
						1138: ByteAlias("that should make things wrap"),
					},
				}
				expected := `{
        Strings: \["lots", "of", "short", "strings"\],
        SimpleThings: \[
            {Name: "short", Enumeration: 7, Veritas: true, Data: "succinct", secret: 17},
            {
                Name: "something longer",
                Enumeration: 427,
                Veritas: true,
                Data: "designed to wrap around nicely",
                secret: 30,
            },
        \],
        DataMaps: {
            (17: "some substantially longer chunks of data"|1138: "that should make things wrap"),
            (17: "some substantially longer chunks of data"|1138: "that should make things wrap"),
        },
    }`
				Expect(Object(s, 1)).Should(matchRegexp(`format_test\.ComplexStruct`, expected))
			})
		})

		Describe("formatting nested interface{} types", func() {
			It("should print out the types of the container and value", func() {
				Expect(Object([]interface{}{"foo"}, 1)).
					To(match("[]interface {} | len:1, cap:1", `[<string>"foo"]`))

				Expect(Object(map[string]interface{}{"foo": true}, 1)).
					To(match("map[string]interface {} | len:1", `{"foo": <bool>true}`))

				Expect(Object(struct{ A interface{} }{A: 1}, 1)).
					To(match("struct { A interface {} }", "{A: <int>1}"))

				v := struct{ A interface{} }{A: struct{ B string }{B: "foo"}}
				Expect(Object(v, 1)).To(match(`struct { A interface {} }`, `{
        A: <struct { B string }>{B: "foo"},
    }`))
			})
		})

		Describe("formatting times", func() {
			It("should format time as RFC3339", func() {
				t := time.Date(2016, 10, 31, 9, 57, 23, 12345, time.UTC)
				Expect(Object(t, 1)).Should(match("time.Time", `2016-10-31T09:57:23.000012345Z`))
			})
		})

		Describe("formatting errors", func() {
			It("should include the error() representation", func() {
				err := fmt.Errorf("whoops: %w", fmt.Errorf("welp: %w", fmt.Errorf("ruh roh")))
				Expect(Object(err, 1)).Should(MatchRegexp(`    \<\*fmt\.wrapError \| 0x[0-9a-f]*\>\: 
    whoops\: welp\: ruh roh
    \{
        msg\: "whoops\: welp\: ruh roh",
        err\: \<\*fmt.wrapError \| 0x[0-9a-f]*\>\{
            msg\: "welp\: ruh roh",
            err\: \<\*errors.errorString \| 0x[0-9a-f]*\>\{s\: "ruh roh"\},
        \},
    \}`))
			})

			It("should not panic if the error is a boxed nil", func() {
				var err *CustomError
				Expect(Object(err, 1)).Should(Equal("    <*format_test.CustomError | 0x0>: nil"))
			})
		})
	})

	Describe("Handling unexported fields in structs", func() {
		It("should handle all the various types correctly", func() {
			a := int(5)
			s := SecretiveStruct{
				boolValue:      true,
				intValue:       3,
				uintValue:      4,
				uintptrValue:   5,
				floatValue:     6.0,
				complexValue:   complex(5.0, 3.0),
				chanValue:      make(chan bool, 2),
				funcValue:      func() {},
				pointerValue:   &a,
				sliceValue:     []string{"string", "slice"},
				byteSliceValue: []byte("bytes"),
				stringValue:    "a string",
				arrValue:       [3]int{11, 12, 13},
				byteArrValue:   [3]byte{17, 20, 32},
				mapValue:       map[string]int{"a key": 20, "b key": 30},
				structValue:    AStruct{"exported"},
				interfaceValue: map[string]int{"a key": 17},
			}

			expected := fmt.Sprintf(`{
        boolValue: true,
        intValue: 3,
        uintValue: 4,
        uintptrValue: 0x5,
        floatValue: 6,
        complexValue: \(5\+3i\),
        chanValue: %p,
        funcValue: %p,
        pointerValue: 5,
        sliceValue: \["string", "slice"\],
        byteSliceValue: "bytes",
        stringValue: "a string",
        arrValue: \[11, 12, 13\],
        byteArrValue: \[17, 20, 32\],
        mapValue: %s,
        structValue: {Exported: "exported"},
        interfaceValue: <map\[string\]int \| len:1>{"a key": 17},
    }`, s.chanValue, s.funcValue, hashMatchingRegexp(`"a key": 20`, `"b key": 30`))

			Expect(Object(s, 1)).Should(matchRegexp(`format_test\.SecretiveStruct`, expected))
		})
	})

	Describe("Handling interfaces", func() {
		It("should unpack the interface", func() {
			outerHash := map[string]interface{}{}
			innerHash := map[string]int{}

			innerHash["inner"] = 3
			outerHash["integer"] = 2
			outerHash["map"] = innerHash

			expected := hashMatchingRegexp(`"integer": <int>2`, `"map": <map\[string\]int \| len:1>{"inner": 3}`)
			Expect(Object(outerHash, 1)).Should(matchRegexp(`map\[string\]interface {} \| len:2`, expected))
		})
	})

	Describe("Handling recursive things", func() {
		It("should not go crazy...", func() {
			m := map[string]interface{}{}
			m["integer"] = 2
			m["map"] = m
			Expect(Object(m, 1)).Should(ContainSubstring("..."))
		})

		It("really should not go crazy...", func() {
			type complexKey struct {
				Value map[interface{}]int
			}

			complexObject := complexKey{}
			complexObject.Value = make(map[interface{}]int)

			complexObject.Value[&complexObject] = 2
			Expect(Object(complexObject, 1)).Should(ContainSubstring("..."))
		})
	})

	Describe("When instructed to use the Stringer representation", func() {
		BeforeEach(func() {
			UseStringerRepresentation = true
		})

		AfterEach(func() {
			UseStringerRepresentation = false
		})

		When("passed a GoStringer", func() {
			It("should use what GoString() returns", func() {
				Expect(Object(GoStringer{}, 1)).Should(ContainSubstring("<format_test.GoStringer>: go-string"))
			})
		})

		When("passed a stringer", func() {
			It("should use what String() returns", func() {
				Expect(Object(Stringer{}, 1)).Should(ContainSubstring("<format_test.Stringer>: string"))
			})
		})

		When("passed a GomegaStringer", func() {
			It("should use what GomegaString() returns", func() {
				Expect(Object(gomegaStringer{}, 1)).Should(ContainSubstring("<format_test.gomegaStringer>: gomegastring"))
				UseStringerRepresentation = false
				Expect(Object(gomegaStringer{}, 1)).Should(ContainSubstring("<format_test.gomegaStringer>: gomegastring"))
			})

			It("should use what GomegaString() returns, disregarding MaxLength", func() {
				Expect(Object(gomegaStringerLong{}, 1)).Should(Equal("    <format_test.gomegaStringerLong>: " + strings.Repeat("s", MaxLength*2)))
				UseStringerRepresentation = false
				Expect(Object(gomegaStringerLong{}, 1)).Should(Equal("    <format_test.gomegaStringerLong>: " + strings.Repeat("s", MaxLength*2)))
			})

			It("should indent what the GomegaString() returns", func() {
				Expect(Object(gomegaStringerMultiline{}, 1)).Should(Equal("    <format_test.gomegaStringerMultiline>: A\n        B\n        C"))
			})
		})

		Describe("when used with a registered CustomFormatter", func() {
			It("pases objects through the custom formatter and uses the returned format, if handled", func() {
				cf := CustomFormatted{"bob", 17}
				ncf := NotCustomFormatted{"bob", 17}
				Expect(Object(cf, 0)).To(Equal("<format_test.CustomFormatted>: {Data: bob, Count: 17}"))
				Expect(Object(ncf, 0)).To(Equal("<format_test.NotCustomFormatted>: {Data: bob, Count: 17}"))

				key := RegisterCustomFormatter(customFormatter)
				Expect(Object(cf, 0)).To(Equal("<format_test.CustomFormatted>: bob (17)"))
				Expect(Object(ncf, 0)).To(Equal("<format_test.NotCustomFormatted>: {Data: bob, Count: 17}"))

				UnregisterCustomFormatter(key)
				Expect(Object(cf, 0)).To(Equal("<format_test.CustomFormatted>: {Data: bob, Count: 17}"))
				Expect(Object(ncf, 0)).To(Equal("<format_test.NotCustomFormatted>: {Data: bob, Count: 17}"))
			})

			It("indents CustomFormatter output correctly", func() {
				cf := CustomFormatted{"hey\nbob", 17}
				DeferCleanup(UnregisterCustomFormatter, RegisterCustomFormatter(func(value interface{}) (string, bool) {
					cf, ok := value.(CustomFormatted)
					if !ok {
						return "", false
					}
					return fmt.Sprintf("The Data:\n%s\nThe Count:%d", cf.Data, cf.Count), true
				}))

				Ω(Object(cf, 1)).Should(Equal("    <format_test.CustomFormatted>: The Data:\n        hey\n        bob\n        The Count:17"))

				type Wrapped struct {
					MyObject   CustomFormatted
					OuterCount int
				}
				wrapped := Wrapped{
					MyObject:   cf,
					OuterCount: 10,
				}
				Ω(Object(wrapped, 1)).Should(Equal("    <format_test.Wrapped>: {\n        MyObject: The Data:\n            hey\n            bob\n            The Count:17,\n        OuterCount: 10,\n    }"))

			})
		})
	})

	Describe("Printing a context.Context field", func() {
		type structWithContext struct {
			Context context.Context
			Value   string
		}

		objWithContext := structWithContext{Value: "some-value", Context: context.TODO()}

		It("Suppresses the content by default", func() {
			Expect(Object(objWithContext, 1)).Should(ContainSubstring("<suppressed context>"))
		})

		It("Doesn't suppress the context if it's the object being printed", func() {
			Expect(Object(context.TODO(), 1)).ShouldNot(MatchRegexp("^.*<suppressed context>$"))
		})

		Context("PrintContextObjects is set", func() {
			BeforeEach(func() {
				PrintContextObjects = true
			})

			AfterEach(func() {
				PrintContextObjects = false
			})

			It("Prints the context", func() {
				Expect(Object(objWithContext, 1)).ShouldNot(ContainSubstring("<suppressed context>"))
			})
		})
	})
})

var expectedLongStringFailureMessage = strings.TrimSpace(`
Expected
    <string>: "...aaaaabaaaaa..."
to equal               |
    <string>: "...aaaaazaaaaa..."
`)
var expectedTruncatedEndStringFailureMessage = strings.TrimSpace(`
Expected
    <string>: "baaaaa..."
to equal       |
    <string>: "zaaaaa..."
`)
var expectedTruncatedStartStringFailureMessage = strings.TrimSpace(`
Expected
    <string>: "...aaaaab"
to equal               |
    <string>: "...aaaaaz"
`)
var expectedTruncatedStartSizeFailureMessage = strings.TrimSpace(`
Expected
    <string>: "...aaaaaa"
to equal               |
    <string>: "...aaaaa"
`)
var expectedTruncatedStartSizeFailureMessageExtraDiff = strings.TrimSpace(`
Expected
    <string>: "...aaaaaaaaaaa"
to equal                    |
    <string>: "...aaaaaaaaaa"
`)
var expectedTruncatedStartSizeSwappedFailureMessage = strings.TrimSpace(`
Expected
    <string>: "...aaaa"
to equal              |
    <string>: "...aaaaa"
`)
var expectedTruncatedStartSizeSwappedFailureMessageExtraDiff = strings.TrimSpace(`
Expected
    <string>: "...aaaaaaaaa"
to equal                   |
    <string>: "...aaaaaaaaaa"
`)
var expectedTruncatedMultiByteFailureMessage = strings.TrimSpace(`
Expected
    <string>: "...tuvwxyz1"
to equal                 |
    <string>: "...tuvwxyz"
`)
var expectedFullFailureDiff = strings.TrimSpace(`
Expected
    <string>: aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaabaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
to equal
    <string>: aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaazaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
`)
var expectedSpecialCharacterFailureMessage = strings.TrimSpace(`
Expected
    <string>: \n
to equal
    <string>: something_else

`)
var expectedDiffSmallThreshold = strings.TrimSpace(`
Expected
    <string>: "aba"
to equal        |
    <string>: "aza"
`)
var expectedDiffZeroMismatch = strings.TrimSpace(`
Expected
    <string>: aba
to equal
    <string>: aza
`)
var expectedDiffSmallThresholdZeroMismatch = strings.TrimSpace(`
Expected
    <string>: "...b..."
to equal          |
    <string>: "...z..."
`)

var expectedDiffLongMessage = strings.TrimSpace(`
Expected
    <string>: "...aaaaabaaaaa..."
to equal very long message
    <string>: "...aaaaazaaaaa..."
`)
