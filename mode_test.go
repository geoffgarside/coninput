//go:build windows
// +build windows

package coninput

import (
	"reflect"
	"testing"

	"golang.org/x/sys/windows"
)

func TestAddInputModes(t *testing.T) {
	tests := []struct {
		name        string
		mode        uint32
		enableModes []uint32
		expected    uint32
	}{
		{
			name:        "add single mode",
			mode:        0,
			enableModes: []uint32{windows.ENABLE_ECHO_INPUT},
			expected:    windows.ENABLE_ECHO_INPUT,
		},
		{
			name:        "add multiple modes",
			mode:        0,
			enableModes: []uint32{windows.ENABLE_ECHO_INPUT, windows.ENABLE_LINE_INPUT},
			expected:    windows.ENABLE_ECHO_INPUT | windows.ENABLE_LINE_INPUT,
		},
		{
			name:        "add to existing mode",
			mode:        windows.ENABLE_ECHO_INPUT,
			enableModes: []uint32{windows.ENABLE_LINE_INPUT},
			expected:    windows.ENABLE_ECHO_INPUT | windows.ENABLE_LINE_INPUT,
		},
		{
			name:        "add duplicate mode",
			mode:        windows.ENABLE_ECHO_INPUT,
			enableModes: []uint32{windows.ENABLE_ECHO_INPUT},
			expected:    windows.ENABLE_ECHO_INPUT,
		},
		{
			name:        "add no modes",
			mode:        windows.ENABLE_ECHO_INPUT,
			enableModes: []uint32{},
			expected:    windows.ENABLE_ECHO_INPUT,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := AddInputModes(tt.mode, tt.enableModes...)
			if result != tt.expected {
				t.Errorf("AddInputModes() = %d, want %d", result, tt.expected)
			}
		})
	}
}

func TestRemoveInputModes(t *testing.T) {
	tests := []struct {
		name         string
		mode         uint32
		disableModes []uint32
		expected     uint32
	}{
		{
			name:         "remove single mode",
			mode:         windows.ENABLE_ECHO_INPUT,
			disableModes: []uint32{windows.ENABLE_ECHO_INPUT},
			expected:     0,
		},
		{
			name:         "remove multiple modes",
			mode:         windows.ENABLE_ECHO_INPUT | windows.ENABLE_LINE_INPUT,
			disableModes: []uint32{windows.ENABLE_ECHO_INPUT, windows.ENABLE_LINE_INPUT},
			expected:     0,
		},
		{
			name:         "remove from multiple modes",
			mode:         windows.ENABLE_ECHO_INPUT | windows.ENABLE_LINE_INPUT | windows.ENABLE_MOUSE_INPUT,
			disableModes: []uint32{windows.ENABLE_LINE_INPUT},
			expected:     windows.ENABLE_ECHO_INPUT | windows.ENABLE_MOUSE_INPUT,
		},
		{
			name:         "remove non-existent mode",
			mode:         windows.ENABLE_ECHO_INPUT,
			disableModes: []uint32{windows.ENABLE_LINE_INPUT},
			expected:     windows.ENABLE_ECHO_INPUT,
		},
		{
			name:         "remove no modes",
			mode:         windows.ENABLE_ECHO_INPUT,
			disableModes: []uint32{},
			expected:     windows.ENABLE_ECHO_INPUT,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := RemoveInputModes(tt.mode, tt.disableModes...)
			if result != tt.expected {
				t.Errorf("RemoveInputModes() = %d, want %d", result, tt.expected)
			}
		})
	}
}

func TestToggleInputModes(t *testing.T) {
	tests := []struct {
		name        string
		mode        uint32
		toggleModes []uint32
		expected    uint32
	}{
		{
			name:        "toggle single mode on",
			mode:        0,
			toggleModes: []uint32{windows.ENABLE_ECHO_INPUT},
			expected:    windows.ENABLE_ECHO_INPUT,
		},
		{
			name:        "toggle single mode off",
			mode:        windows.ENABLE_ECHO_INPUT,
			toggleModes: []uint32{windows.ENABLE_ECHO_INPUT},
			expected:    0,
		},
		{
			name:        "toggle multiple modes",
			mode:        windows.ENABLE_ECHO_INPUT,
			toggleModes: []uint32{windows.ENABLE_ECHO_INPUT, windows.ENABLE_LINE_INPUT},
			expected:    windows.ENABLE_LINE_INPUT,
		},
		{
			name:        "toggle no modes",
			mode:        windows.ENABLE_ECHO_INPUT,
			toggleModes: []uint32{},
			expected:    windows.ENABLE_ECHO_INPUT,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToggleInputModes(tt.mode, tt.toggleModes...)
			if result != tt.expected {
				t.Errorf("ToggleInputModes() = %d, want %d", result, tt.expected)
			}
		})
	}
}

func TestListInputModes(t *testing.T) {
	tests := []struct {
		name     string
		mode     uint32
		expected []uint32
	}{
		{
			name:     "no modes",
			mode:     0,
			expected: []uint32{},
		},
		{
			name:     "single mode",
			mode:     windows.ENABLE_ECHO_INPUT,
			expected: []uint32{windows.ENABLE_ECHO_INPUT},
		},
		{
			name:     "multiple modes",
			mode:     windows.ENABLE_ECHO_INPUT | windows.ENABLE_LINE_INPUT,
			expected: []uint32{windows.ENABLE_ECHO_INPUT, windows.ENABLE_LINE_INPUT},
		},
		{
			name:     "all common modes",
			mode:     windows.ENABLE_ECHO_INPUT | windows.ENABLE_INSERT_MODE | windows.ENABLE_LINE_INPUT,
			expected: []uint32{windows.ENABLE_ECHO_INPUT, windows.ENABLE_INSERT_MODE, windows.ENABLE_LINE_INPUT},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ListInputModes(tt.mode)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ListInputModes() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestListInputModeNames(t *testing.T) {
	tests := []struct {
		name     string
		mode     uint32
		expected []string
	}{
		{
			name:     "no modes",
			mode:     0,
			expected: []string{},
		},
		{
			name:     "single mode",
			mode:     windows.ENABLE_ECHO_INPUT,
			expected: []string{"ENABLE_ECHO_INPUT"},
		},
		{
			name:     "multiple modes",
			mode:     windows.ENABLE_ECHO_INPUT | windows.ENABLE_LINE_INPUT,
			expected: []string{"ENABLE_ECHO_INPUT", "ENABLE_LINE_INPUT"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ListInputModeNames(tt.mode)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("ListInputModeNames() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestDescribeInputMode(t *testing.T) {
	tests := []struct {
		name     string
		mode     uint32
		expected string
	}{
		{
			name:     "no modes",
			mode:     0,
			expected: "",
		},
		{
			name:     "single mode",
			mode:     windows.ENABLE_ECHO_INPUT,
			expected: "ENABLE_ECHO_INPUT",
		},
		{
			name:     "multiple modes",
			mode:     windows.ENABLE_ECHO_INPUT | windows.ENABLE_LINE_INPUT,
			expected: "ENABLE_ECHO_INPUT|ENABLE_LINE_INPUT",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DescribeInputMode(tt.mode)
			if result != tt.expected {
				t.Errorf("DescribeInputMode() = %q, want %q", result, tt.expected)
			}
		})
	}
}