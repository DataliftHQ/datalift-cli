package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRootCmd(t *testing.T) {
	mem := &exitMemento{}
	Execute("1.2.3", mem.Exit, []string{"-h"})
	require.Equal(t, 0, mem.code)
}

func TestRootCmdHelp(t *testing.T) {
	mem := &exitMemento{}
	cmd := newRootCmd("", mem.Exit).cmd
	cmd.SetArgs([]string{"-h"})
	require.NoError(t, cmd.Execute())
	require.Equal(t, 0, mem.code)
}

func TestRootCmdVersion(t *testing.T) {
	var b bytes.Buffer
	mem := &exitMemento{}
	cmd := newRootCmd("1.2.3", mem.Exit).cmd
	cmd.SetOut(&b)
	cmd.SetArgs([]string{"-v"})
	require.NoError(t, cmd.Execute())
	require.Equal(t, "updog version 1.2.3\n", b.String())
	require.Equal(t, 0, mem.code)
}

func TestRootCmdExitCode(t *testing.T) {
	mem := &exitMemento{}
	cmd := newRootCmd("", mem.Exit)
	args := []string{"check", "-f", "testdata/good.hcl"}
	cmd.Execute(args)
	require.Equal(t, 0, mem.code)
}

func TestRootPush(t *testing.T) {
	setup(t)
	mem := &exitMemento{}
	cmd := newRootCmd("", mem.Exit)
	cmd.Execute([]string{})
	require.Equal(t, 0, mem.code)
}

func TestRootPushDebug(t *testing.T) {
	setup(t)
	mem := &exitMemento{}
	cmd := newRootCmd("", mem.Exit)
	cmd.Execute([]string{"p", "--debug"})
	require.Equal(t, 0, mem.code)
}

func TestShouldPrependRelease(t *testing.T) {
	result := func(args []string) bool {
		return shouldPrependPush(newRootCmd("1", func(_ int) {}).cmd, args)
	}

	t.Run("no args", func(t *testing.T) {
		require.True(t, result([]string{}))
	})

	t.Run("release args", func(t *testing.T) {
		require.True(t, result([]string{"--skip-validate"}))
	})

	t.Run("several release args", func(t *testing.T) {
		require.True(t, result([]string{"--skip-validate", "--snapshot"}))
	})

	for _, s := range []string{"--help", "-h", "-v", "--version"} {
		t.Run(s, func(t *testing.T) {
			require.False(t, result([]string{s}))
		})
	}

	t.Run("check", func(t *testing.T) {
		require.False(t, result([]string{"check", "-f", "testdata/good.hcl"}))
	})

	t.Run("help", func(t *testing.T) {
		require.False(t, result([]string{"help"}))
	})

	t.Run("__complete", func(t *testing.T) {
		require.False(t, result([]string{"__complete"}))
	})

	t.Run("__completeNoDesc", func(t *testing.T) {
		require.False(t, result([]string{"__completeNoDesc"}))
	})
}
