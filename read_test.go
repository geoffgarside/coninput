//go:build windows
// +build windows

package coninput

import (
	"testing"

	"golang.org/x/sys/windows"
)

func TestReadNConsoleInputs(t *testing.T) {
	tests := []struct {
		name      string
		maxEvents uint32
		wantError bool
		errorMsg  string
	}{
		{
			name:      "zero max events",
			maxEvents: 0,
			wantError: true,
			errorMsg:  "maxEvents cannot be zero",
		},
		{
			name:      "valid max events",
			maxEvents: 1,
			wantError: false,
		},
		{
			name:      "multiple max events",
			maxEvents: 10,
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock handle - this will fail on actual Windows API call
			// but we're testing the validation logic
			handle := windows.Handle(0)
			
			_, err := ReadNConsoleInputs(handle, tt.maxEvents)
			
			if tt.wantError {
				if err == nil {
					t.Error("Expected error, got nil")
				} else if err.Error() != tt.errorMsg {
					t.Errorf("Expected error %q, got %q", tt.errorMsg, err.Error())
				}
			} else {
				// For non-zero maxEvents, we expect a Windows API error since we're not on Windows
				// or using a valid handle, but not our validation error
				if err != nil && err.Error() == tt.errorMsg {
					t.Errorf("Got validation error when expecting Windows API error: %v", err)
				}
			}
		})
	}
}

func TestReadConsoleInput(t *testing.T) {
	tests := []struct {
		name         string
		inputRecords []InputRecord
		wantError    bool
		errorMsg     string
	}{
		{
			name:         "empty input records",
			inputRecords: []InputRecord{},
			wantError:    true,
			errorMsg:     "size of input record buffer cannot be zero",
		},
		{
			name:         "single input record",
			inputRecords: make([]InputRecord, 1),
			wantError:    false,
		},
		{
			name:         "multiple input records",
			inputRecords: make([]InputRecord, 5),
			wantError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock handle
			handle := windows.Handle(0)
			
			_, err := ReadConsoleInput(handle, tt.inputRecords)
			
			if tt.wantError {
				if err == nil {
					t.Error("Expected error, got nil")
				} else if err.Error() != tt.errorMsg {
					t.Errorf("Expected error %q, got %q", tt.errorMsg, err.Error())
				}
			} else {
				// For non-empty slices, we expect a Windows API error since we're not on Windows
				// or using a valid handle, but not our validation error
				if err != nil && err.Error() == tt.errorMsg {
					t.Errorf("Got validation error when expecting Windows API error: %v", err)
				}
			}
		})
	}
}

func TestPeekNConsoleInputs(t *testing.T) {
	tests := []struct {
		name      string
		maxEvents uint32
		wantError bool
		errorMsg  string
	}{
		{
			name:      "zero max events",
			maxEvents: 0,
			wantError: true,
			errorMsg:  "maxEvents cannot be zero",
		},
		{
			name:      "valid max events",
			maxEvents: 1,
			wantError: false,
		},
		{
			name:      "multiple max events",
			maxEvents: 10,
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock handle
			handle := windows.Handle(0)
			
			_, err := PeekNConsoleInputs(handle, tt.maxEvents)
			
			if tt.wantError {
				if err == nil {
					t.Error("Expected error, got nil")
				} else if err.Error() != tt.errorMsg {
					t.Errorf("Expected error %q, got %q", tt.errorMsg, err.Error())
				}
			} else {
				// For non-zero maxEvents, we expect a Windows API error
				if err != nil && err.Error() == tt.errorMsg {
					t.Errorf("Got validation error when expecting Windows API error: %v", err)
				}
			}
		})
	}
}

func TestPeekConsoleInput(t *testing.T) {
	tests := []struct {
		name         string
		inputRecords []InputRecord
		wantError    bool
		errorMsg     string
	}{
		{
			name:         "empty input records",
			inputRecords: []InputRecord{},
			wantError:    true,
			errorMsg:     "size of input record buffer cannot be zero",
		},
		{
			name:         "single input record",
			inputRecords: make([]InputRecord, 1),
			wantError:    false,
		},
		{
			name:         "multiple input records",
			inputRecords: make([]InputRecord, 5),
			wantError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock handle
			handle := windows.Handle(0)
			
			_, err := PeekConsoleInput(handle, tt.inputRecords)
			
			if tt.wantError {
				if err == nil {
					t.Error("Expected error, got nil")
				} else if err.Error() != tt.errorMsg {
					t.Errorf("Expected error %q, got %q", tt.errorMsg, err.Error())
				}
			} else {
				// For non-empty slices, we expect a Windows API error
				if err != nil && err.Error() == tt.errorMsg {
					t.Errorf("Got validation error when expecting Windows API error: %v", err)
				}
			}
		})
	}
}

// Test helper functions that don't require Windows API calls

func TestNewStdinHandle(t *testing.T) {
	// This will fail on non-Windows systems, but we can test that it returns something
	_, err := NewStdinHandle()
	// We expect this to fail on non-Windows systems
	if err == nil {
		t.Log("NewStdinHandle() succeeded (running on Windows)")
	} else {
		t.Logf("NewStdinHandle() failed as expected on non-Windows: %v", err)
	}
}

// Test validation logic without actual Windows API calls
func TestValidationLogic(t *testing.T) {
	t.Run("ReadNConsoleInputs validation", func(t *testing.T) {
		// Test the validation logic by checking if zero maxEvents produces the right error
		handle := windows.Handle(0)
		_, err := ReadNConsoleInputs(handle, 0)
		if err == nil || err.Error() != "maxEvents cannot be zero" {
			t.Errorf("Expected 'maxEvents cannot be zero' error, got: %v", err)
		}
	})

	t.Run("ReadConsoleInput validation", func(t *testing.T) {
		// Test the validation logic by checking if empty slice produces the right error
		handle := windows.Handle(0)
		emptySlice := []InputRecord{}
		_, err := ReadConsoleInput(handle, emptySlice)
		if err == nil || err.Error() != "size of input record buffer cannot be zero" {
			t.Errorf("Expected 'size of input record buffer cannot be zero' error, got: %v", err)
		}
	})

	t.Run("PeekNConsoleInputs validation", func(t *testing.T) {
		handle := windows.Handle(0)
		_, err := PeekNConsoleInputs(handle, 0)
		if err == nil || err.Error() != "maxEvents cannot be zero" {
			t.Errorf("Expected 'maxEvents cannot be zero' error, got: %v", err)
		}
	})

	t.Run("PeekConsoleInput validation", func(t *testing.T) {
		handle := windows.Handle(0)
		emptySlice := []InputRecord{}
		_, err := PeekConsoleInput(handle, emptySlice)
		if err == nil || err.Error() != "size of input record buffer cannot be zero" {
			t.Errorf("Expected 'size of input record buffer cannot be zero' error, got: %v", err)
		}
	})
}

// Test error handling for Windows API functions that we can't easily mock
func TestWindowsAPIErrorHandling(t *testing.T) {
	// These tests verify that the functions properly handle invalid handles
	// and return appropriate errors rather than panicking
	invalidHandle := windows.Handle(0xFFFFFFFF) // Invalid handle

	t.Run("GetNumberOfConsoleInputEvents with invalid handle", func(t *testing.T) {
		_, err := GetNumberOfConsoleInputEvents(invalidHandle)
		// We expect an error, not a panic
		if err == nil {
			t.Error("Expected error with invalid handle")
		}
	})

	t.Run("FlushConsoleInputBuffer with invalid handle", func(t *testing.T) {
		err := FlushConsoleInputBuffer(invalidHandle)
		// We expect an error, not a panic
		if err == nil {
			t.Error("Expected error with invalid handle")
		}
	})
}