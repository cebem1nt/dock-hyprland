package main

import (
	"errors"
	"testing"
)

func TestDispatchReplyFailed(t *testing.T) {
	tests := []struct {
		name  string
		reply string
		err   error
		want  bool
	}{
		{name: "ok", reply: "ok", want: false},
		{name: "empty", reply: "", want: false},
		{name: "other success text", reply: "some compositor reply", want: false},
		{name: "error reply", reply: "error: parser failed", want: true},
		{name: "spaced uppercase error reply", reply: "  ERROR: parser failed", want: true},
		{name: "transport error", err: errors.New("socket failed"), want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := dispatchReplyFailed([]byte(tt.reply), tt.err)
			if got != tt.want {
				t.Fatalf("dispatchReplyFailed(%q, %v) = %v, want %v", tt.reply, tt.err, got, tt.want)
			}
		})
	}
}

func TestFocusWindowCommands(t *testing.T) {
	address := "0x5c45b8eef120"

	if got, want := legacyFocusWindowCommand(address), "dispatch focuswindow address:0x5c45b8eef120"; got != want {
		t.Fatalf("legacyFocusWindowCommand() = %q, want %q", got, want)
	}

	if got, want := luaFocusWindowCommand(address), `dispatch hl.dsp.focus({ window = "address:0x5c45b8eef120" })`; got != want {
		t.Fatalf("luaFocusWindowCommand() = %q, want %q", got, want)
	}
}

func TestBringActiveToTopCommands(t *testing.T) {
	if got, want := legacyBringActiveToTopCommand(), "dispatch bringactivetotop"; got != want {
		t.Fatalf("legacyBringActiveToTopCommand() = %q, want %q", got, want)
	}

	if got, want := luaBringActiveToTopCommand(), "dispatch hl.dsp.window.bring_to_top()"; got != want {
		t.Fatalf("luaBringActiveToTopCommand() = %q, want %q", got, want)
	}
}

func TestToggleSpecialWorkspaceCommands(t *testing.T) {
	name := "scratchpad"

	if got, want := legacyToggleSpecialWorkspaceCommand(name), "dispatch togglespecialworkspace scratchpad"; got != want {
		t.Fatalf("legacyToggleSpecialWorkspaceCommand() = %q, want %q", got, want)
	}

	if got, want := luaToggleSpecialWorkspaceCommand(name), `dispatch hl.dsp.workspace.toggle_special("scratchpad")`; got != want {
		t.Fatalf("luaToggleSpecialWorkspaceCommand() = %q, want %q", got, want)
	}
}

func TestLuaStringEscapesValues(t *testing.T) {
	if got, want := luaString(`address:0xabc"def`), `"address:0xabc\"def"`; got != want {
		t.Fatalf("luaString() = %q, want %q", got, want)
	}
}
