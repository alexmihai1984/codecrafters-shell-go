package main

import "testing"

func TestBackslashOutsideQuotes(t *testing.T) {
	stdout, stderr := runShell(t, []string{
		`echo world\ \ \ \ \ \ script`,
		`echo before\nafter`,
		`echo \'\"script world\"\'`,
		"exit",
	})

	expectedStdout := "$ world      script\n" +
		"$ beforenafter\n" +
		"$ '\"script world\"'\n" +
		"$ "

	if stdout != expectedStdout {
		t.Errorf("unexpected stdout:\nexpected: %q\ngot:      %q", expectedStdout, stdout)
	}
	if stderr != "" {
		t.Errorf("unexpected stderr: %q", stderr)
	}
}

func TestBackslashWithinSingleQuotes(t *testing.T) {
	stdout, stderr := runShell(t, []string{
		`echo 'hello\"world'`,
		`echo 'hello\\world'`,
		`echo 'hello\world'`,
		`exit`,
	})

	expectStdout(
		t,
		stdout,
		[]string{
			`$ hello\"world`,
			`$ hello\\world`,
			`$ hello\world`,
			`$ `,
		},
	)

	expectStderr(
		t,
		stderr,
		[]string{},
	)
}

func TestBackslashWithinDoubleQuotes(t *testing.T) {
	stdout, stderr := runShell(t, []string{
		`echo "hello\"world"`,
		`echo "hello\'world"`,
		`echo "hello\\world"`,
		`echo "hello\world"`,
		"exit",
	})

	expectStdout(
		t,
		stdout,
		[]string{
			"$ hello\"world",
			"$ hello\\'world",
			"$ hello\\world",
			"$ hello\\world",
			"$ ",
		},
	)

	expectStderr(
		t,
		stderr,
		[]string{},
	)
}
