package main

import "testing"

func TestBackslashOutsideQuotes(t *testing.T) {
	stdout, stderr := runShell(t, []string{
		`echo world\ \ \ \ \ \ script`,
		`echo before\nafter`,
		`echo \'\"script world\"\'`,
		"exit",
	})

	expectStdout(
		t,
		stdout,
		[]string{
			"$ world      script",
			"$ beforenafter",
			"$ '\"script world\"'",
			"$ ",
		},
	)

	expectStderr(t, stderr, []string{})
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

	expectStderr(t, stderr, []string{})
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

	expectStderr(t, stderr, []string{})
}
