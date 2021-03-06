package test161

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

// A random string input generator.  This generates between 2 and 10
// strings, each of length between 5 and 10 characters.
const input_template = "{{$x := randInt 2 10 | ranger}}{{range $index, $element := $x}}{{randString 5 10}}\n{{end}}"

func TestCommandArgTest(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	// Create a test
	args := []string{"arg1", "arg2", "arg3", "arg4"}
	cmdline := "$ /testbin/argtest"
	for _, a := range args {
		cmdline += " " + a
	}
	test, err := TestFromString(cmdline)
	assert.Nil(err)

	// Set the commands for argtest
	var argtest *Command
	for _, c := range test.Commands {
		if c.Id() == "/testbin/argtest" {
			argtest = c
			break
		}
	}

	assert.NotNil(argtest)
	if argtest == nil {
		t.FailNow()
	}

	// Create the command instance
	err = argtest.Instantiate(defaultEnv)

	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	for _, o := range argtest.ExpectedOutput {
		t.Log(o.Text)
	}

	// Assertions
	assert.Equal(3+len(args), len(argtest.ExpectedOutput))

	if len(argtest.ExpectedOutput) == 7 {
		assert.Equal("argc: 5", argtest.ExpectedOutput[0].Text)
		assert.Equal("argv[0]: /testbin/argtest", argtest.ExpectedOutput[1].Text)
		for i, arg := range args {
			assert.Equal(fmt.Sprintf("argv[%d]: %v", i+1, arg), argtest.ExpectedOutput[i+2].Text)
		}
	}

}

func TestCommandAdd(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	// Create a test
	test, err := TestFromString("$ /testbin/add 70 200")
	assert.Nil(err)

	// Set the commands for argtest
	var add *Command
	for _, c := range test.Commands {
		if c.Id() == "/testbin/add" {
			add = c
			break
		}
	}

	assert.NotNil(add)
	if add == nil {
		t.FailNow()
	}

	// Create the command instance
	err = add.Instantiate(defaultEnv)

	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	for _, o := range add.ExpectedOutput {
		t.Log(o.Text)
	}

	// Assertions
	assert.Equal(1, len(add.ExpectedOutput))
	if len(add.ExpectedOutput) == 1 {
		assert.Equal("270", add.ExpectedOutput[0].Text)
	}
}

func TestCommandFactorial(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	// Create a test
	test, err := TestFromString("$ /testbin/factorial 8")
	assert.Nil(err)

	// Set the commands for argtest
	var factorial *Command
	for _, c := range test.Commands {
		if c.Id() == "/testbin/factorial" {
			factorial = c
			break
		}
	}

	assert.NotNil(factorial)
	if factorial == nil {
		t.FailNow()
	}

	// Create the command instance
	err = factorial.Instantiate(defaultEnv)

	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	for _, o := range factorial.ExpectedOutput {
		t.Log(o.Text)
	}

	// Assertions
	assert.Equal(1, len(factorial.ExpectedOutput))
	if len(factorial.ExpectedOutput) == 1 {
		assert.Equal("40320", factorial.ExpectedOutput[0].Text)
	}
}

func addInputTest() (*TestEnvironment, error) {

	env, err := NewEnvironment("./fixtures", &DoNothingPersistence{})
	if err != nil {
		return nil, err
	}
	env.TestDir = "./fixtures/tests/nocycle/"

	// Create the Command Template for (fake) randinput.
	c := &CommandTemplate{
		Name:  "randinput",
		Input: []string{input_template},
	}

	env.Commands["randinput"] = c
	return env, nil
}

func TestCommandInput(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	env, err := addInputTest()

	assert.Nil(err)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	// Create a test
	test, err := TestFromString("randinput")
	assert.Nil(err)

	var randinput *Command
	for _, c := range test.Commands {
		if c.Id() == "randinput" {
			randinput = c
			break
		}
	}

	assert.NotNil(randinput)
	if randinput == nil {
		t.FailNow()
	}

	// Create the command instance
	err = randinput.Instantiate(env)

	t.Log(randinput.Input.Line)

	for _, o := range randinput.ExpectedOutput {
		t.Log(o.Text)
	}

	_, id, args := randinput.Input.splitCommand()

	t.Log(args)
	t.Log(id)

	// Assertions
	assert.True(len(args) >= 2)
	assert.True(len(args) <= 10)

	for _, o := range args {
		assert.True(len(o) >= 5)
		assert.True(len(o) <= 10)
	}

	// Now, check override
	randinput.Input.Line = "randinput 1"
	randinput.ExpectedOutput = nil

	randinput.Instantiate(defaultEnv)

	_, id, args = randinput.Input.splitCommand()

	assert.Equal(0, len(randinput.ExpectedOutput))
	assert.Equal(1, len(args))
	if len(args) == 1 {
		assert.Equal("1", args[0])
	}
	assert.Equal("randinput", id)
	assert.Equal("randinput 1", randinput.Input.Line)
}

func TestCommandTemplateLoad(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	text := `---
templates:
  - name: sy1
  - name: sy2
  - name: sy3
  - name: sy4
  - name: sy5
`
	cmds, err := CommandTemplatesFromString(text)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	assert.Equal(5, len(cmds.Templates))
	if len(cmds.Templates) == 5 {
		for i, tmpl := range cmds.Templates {
			assert.Equal(fmt.Sprintf("sy%v", i+1), tmpl.Name)
			assert.Equal(1, len(tmpl.Output))
			if len(tmpl.Output) == 1 {
				assert.Equal(tmpl.Name+": SUCCESS", tmpl.Output[0].Text)
				assert.Equal("true", tmpl.Output[0].Trusted)
				assert.Equal("false", tmpl.Output[0].External)
			}
		}
	}
}

func addExternalCmd() (*TestEnvironment, error) {
	env, err := NewEnvironment("./fixtures", &DoNothingPersistence{})
	if err != nil {
		return nil, err
	}
	env.TestDir = "./fixtures/tests/nocycle/"

	// Create the Command Template for (fake) randinput.
	c := &CommandTemplate{
		Name: "external",
		Output: []*TemplOutputLine{
			&TemplOutputLine{
				Text:     "sem1",
				Trusted:  "true",
				External: "true",
			},
			&TemplOutputLine{
				Text:     "lt1",
				Trusted:  "true",
				External: "true",
			},
		},
	}

	env.Commands["external"] = c
	return env, nil
}

func TestCommandExternal(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	env, err := addExternalCmd()
	assert.Nil(err)
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	test, err := TestFromString("external")
	assert.Nil(err)

	var cmd *Command
	for _, c := range test.Commands {
		if c.Input.Line == "external" {
			cmd = c
			break
		}
	}

	assert.NotNil(cmd)
	if cmd == nil {
		t.FailNow()
	}

	// Create the command instance
	err = cmd.Instantiate(env)

	for _, o := range cmd.ExpectedOutput {
		t.Log(o.Text)
	}

	t.Log(cmd.ExpectedOutput)

	// Assertions
	assert.Equal(2, len(cmd.ExpectedOutput))
	if len(cmd.ExpectedOutput) != 2 {
		t.FailNow()
	}

	assert.Equal("sem1: SUCCESS", cmd.ExpectedOutput[0].Text)
	assert.True(cmd.ExpectedOutput[0].Trusted)
	assert.Equal("sem1", cmd.ExpectedOutput[0].KeyName)

	assert.Equal("lt1: SUCCESS", cmd.ExpectedOutput[1].Text)
	assert.True(cmd.ExpectedOutput[1].Trusted)
	assert.Equal("lt1", cmd.ExpectedOutput[1].KeyName)
}

func TestCommandID(t *testing.T) {

	t.Parallel()
	assert := assert.New(t)

	tests := [][]string{
		[]string{
			"/hello/world", "/hello/world",
		},
		[]string{
			"/hello/world ", "/hello/world",
		},
		[]string{
			`/testbin/argtest 1 2 3`, "/testbin/argtest",
		},
		[]string{
			`/bin/space test`, "/bin/space",
		},
		[]string{
			`/bin/space\ test`, `/bin/space\ test`,
		},
		[]string{
			`"/bin/space test" 1 2 3`, `"/bin/space test"`,
		},
		[]string{
			`p /bin/something 1 2 3`, `/bin/something`,
		},
		[]string{
			`p s 1 2 3`, `s`,
		},
	}

	for _, test := range tests {
		line := &InputLine{Line: test[0]}
		pfx, base, args := line.splitCommand()
		assert.Equal(test[1], base)
		if strings.HasPrefix(test[0], "p ") {
			assert.Equal("p", pfx)
		}
		t.Log(base, test[1])
		t.Log(args)
	}

}

func TestCommandReplaceArgs(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	type replaceTest struct {
		original string
		newArgs  []string
		expected string
	}

	tests := []*replaceTest{
		&replaceTest{
			"/hello/world", []string{"arg1"}, "/hello/world arg1",
		},
		&replaceTest{
			"/foo/bar *", []string{"-r", "25"}, "/foo/bar -r 25",
		},
		&replaceTest{
			"/foo/bar -t 100", []string{}, "/foo/bar",
		},
		&replaceTest{
			"p /hello/world", []string{"arg1"}, "p /hello/world arg1",
		},
		&replaceTest{
			"p /foo/bar *", []string{"-r", "25"}, "p /foo/bar -r 25",
		},
		&replaceTest{
			"p /foo/bar -t 100", []string{}, "p /foo/bar",
		},
	}

	for _, test := range tests {
		line := &InputLine{Line: test.original}
		line.replaceArgs(test.newArgs)
		assert.Equal(line.Line, test.expected)
		t.Log(test.original, test.newArgs, test.expected)
	}
}

func TestCommandOverride(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	testString := `---
commandoverrides:
  - name: sem1
    output:
      - text: "Override SUCCESS"
---
sem1
`
	test, err := TestFromString(testString)
	assert.Nil(err)

	// find the sem1 command
	var cmd *Command
	for _, c := range test.Commands {
		if c.Id() == "sem1" {
			cmd = c
			break
		}
	}

	err = cmd.Instantiate(defaultEnv)
	assert.Nil(err)

	if cmd == nil {
		t.Fatalf("cmd == nil)")
	}

	if len(cmd.ExpectedOutput) != 1 {
		t.Fatalf("Command output != 1")
	}

	assert.Equal("Override SUCCESS", cmd.ExpectedOutput[0].Text)

	// These are overrides only, so if these aren't specified, they get the
	// default values.
	assert.Equal(false, cmd.ExpectedOutput[0].Trusted)
	assert.Equal("", cmd.ExpectedOutput[0].KeyName)

	testString = `---
commandoverrides:
  - name: sem1
    timeout: 1000.0
---
sem1
`

	test, err = TestFromString(testString)
	assert.Nil(err)

	// find the sem1 command
	cmd = nil
	for _, c := range test.Commands {
		if c.Id() == "sem1" {
			cmd = c
			break
		}
	}

	err = cmd.Instantiate(defaultEnv)
	assert.Nil(err)

	if cmd == nil {
		t.Fatalf("cmd == nil)")
	}

	if len(cmd.ExpectedOutput) != 1 {
		t.Fatalf("Command output != 1")
	}

	assert.Equal("sem1: SUCCESS", cmd.ExpectedOutput[0].Text)
	assert.Equal(true, cmd.ExpectedOutput[0].Trusted)
	assert.Equal("sem1", cmd.ExpectedOutput[0].KeyName)

	assert.Equal(float32(1000.0), cmd.Timeout)

}
