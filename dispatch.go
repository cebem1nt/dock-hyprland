package main

import (
	"fmt"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

func dispatchReplyFailed(reply []byte, err error) bool {
	if err != nil {
		return true
	}

	text := strings.TrimSpace(string(reply))
	return strings.HasPrefix(strings.ToLower(text), "error:")
}

func dispatchWithFallback(primary string, fallback string) ([]byte, error) {
	reply, err := hyprctl(primary)
	if !dispatchReplyFailed(reply, err) || fallback == "" {
		return reply, err
	}

	log.Debugf("%s -> %s", primary, reply)
	log.Debugf("primary dispatch failed, trying fallback: %s", fallback)
	return hyprctl(fallback)
}

func luaString(value string) string {
	return strconv.Quote(value)
}

func legacyFocusWindowCommand(address string) string {
	return fmt.Sprintf("dispatch focuswindow address:%s", address)
}

func luaFocusWindowCommand(address string) string {
	return fmt.Sprintf("dispatch hl.dsp.focus({ window = %s })", luaString("address:"+address))
}

func legacyBringActiveToTopCommand() string {
	return "dispatch bringactivetotop"
}

func luaBringActiveToTopCommand() string {
	return "dispatch hl.dsp.window.bring_to_top()"
}

func legacyToggleSpecialWorkspaceCommand(name string) string {
	return fmt.Sprintf("dispatch togglespecialworkspace %s", name)
}

func luaToggleSpecialWorkspaceCommand(name string) string {
	return fmt.Sprintf("dispatch hl.dsp.workspace.toggle_special(%s)", luaString(name))
}

func focusWindow(address string) {
	primary := legacyFocusWindowCommand(address)
	fallback := luaFocusWindowCommand(address)
	reply, err := dispatchWithFallback(primary, fallback)
	log.Debugf("focus window %s -> %s (%v)", address, reply, err)
}

func bringActiveToTop() {
	primary := legacyBringActiveToTopCommand()
	fallback := luaBringActiveToTopCommand()
	reply, err := dispatchWithFallback(primary, fallback)
	log.Debugf("bring active to top -> %s (%v)", reply, err)
}

func toggleSpecialWorkspace(name string) {
	primary := legacyToggleSpecialWorkspaceCommand(name)
	fallback := luaToggleSpecialWorkspaceCommand(name)
	reply, err := dispatchWithFallback(primary, fallback)
	log.Debugf("toggle special workspace %s -> %s (%v)", name, reply, err)
}

func focusClient(c client) {
	if strings.HasPrefix(c.Workspace.Name, "special") {
		_, specialName, found := strings.Cut(c.Workspace.Name, "special:")
		if found {
			toggleSpecialWorkspace(specialName)
		} else {
			toggleSpecialWorkspace("")
		}
	}

	focusWindow(c.Address)
	bringActiveToTop()
}
