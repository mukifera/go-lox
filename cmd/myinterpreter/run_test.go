package main

import (
	"bytes"
	"testing"
)

func TestRuntime(t *testing.T) {
	tests := []struct {
		name string
		fileContents string
		expected string
	}{
		{"Print #1", `print "Hello, World!";`, "Hello, World!\n"},
		{"Print #2", "print 42;", "42\n"},
		{"Print #3", "print true;", "true\n"},
		{"Print #4", "print 12 + 24;", "36\n"},

		{
			"Multiple Statements #1",

`print "world" + "baz" + "bar";
print 27 - 26;
print "bar" == "quz";`,

`worldbazbar
1
false
`,
		},

		{
			"Multiple Statements #2",

`print "hello"; print true;
print false;
print "bar"; print 43;`,

`hello
true
false
bar
43
`,
		},

		{
			"Multiple Statements #3",
		
`print 81;
    print 81 + 46;
        print 81 + 46 + 19;`,
		
`81
127
146
`,
		},
		{
			"Multiple Statements #4",
		
`print true != true;

print "36
10
78
";

print "There should be an empty line above this.";`,
		
`false
36
10
78

There should be an empty line above this.
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exprs := getExpressions(tt.fileContents, t)
			var buf bytes.Buffer
			FRunExpressions(&buf, exprs)
			str := buf.String()

			if tt.expected != str {
				t.Errorf("Execution result mismatch\nExpected:\n\n%s\nGot:\n\n%s", tt.expected, str)
			}
		})
	}
}