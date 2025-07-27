package coninput

import (
	"encoding/binary"
	"strings"
	"testing"
)

func TestInputRecordUnwrap(t *testing.T) {
	tests := []struct {
		name      string
		eventType EventType
		eventData EventUnion
		expected  EventRecord
	}{
		{
			name:      "focus event - set focus true",
			eventType: FocusEventType,
			eventData: func() EventUnion {
				var data EventUnion
				data[0] = 1
				return data
			}(),
			expected: FocusEventRecord{SetFocus: true},
		},
		{
			name:      "focus event - set focus false",
			eventType: FocusEventType,
			eventData: func() EventUnion {
				var data EventUnion
				data[0] = 0
				return data
			}(),
			expected: FocusEventRecord{SetFocus: false},
		},
		{
			name:      "key event - key down",
			eventType: KeyEventType,
			eventData: func() EventUnion {
				var data EventUnion
				binary.LittleEndian.PutUint32(data[0:4], 1)   // KeyDown = true
				binary.LittleEndian.PutUint16(data[4:6], 2)   // RepeatCount = 2
				binary.LittleEndian.PutUint16(data[6:8], 65)  // VirtualKeyCode = 'A'
				binary.LittleEndian.PutUint16(data[8:10], 30) // VirtualScanCode = 30
				binary.LittleEndian.PutUint16(data[10:12], 97) // Char = 'a'
				binary.LittleEndian.PutUint32(data[12:16], 0)  // ControlKeyState = 0
				return data
			}(),
			expected: KeyEventRecord{
				KeyDown:         true,
				RepeatCount:     2,
				VirtualKeyCode:  65,
				VirtualScanCode: 30,
				Char:            'a',
				ControlKeyState: 0,
			},
		},
		{
			name:      "mouse event",
			eventType: MouseEventType,
			eventData: func() EventUnion {
				var data EventUnion
				binary.LittleEndian.PutUint16(data[0:2], 10)  // X = 10
				binary.LittleEndian.PutUint16(data[2:4], 20)  // Y = 20
				binary.LittleEndian.PutUint32(data[4:8], 1)   // ButtonState = 1
				binary.LittleEndian.PutUint32(data[8:12], 0)  // ControlKeyState = 0
				binary.LittleEndian.PutUint32(data[12:16], 0) // EventFlags = 0
				return data
			}(),
			expected: MouseEventRecord{
				MousePositon:    Coord{X: 10, Y: 20},
				ButtonState:     1,
				ControlKeyState: 0,
				EventFlags:      0,
				WheelDirection:  0,
			},
		},
		{
			name:      "window buffer size event",
			eventType: WindowBufferSizeEventType,
			eventData: func() EventUnion {
				var data EventUnion
				binary.LittleEndian.PutUint16(data[0:2], 80) // X = 80
				binary.LittleEndian.PutUint16(data[2:4], 25) // Y = 25
				return data
			}(),
			expected: WindowBufferSizeEventRecord{
				Size: Coord{X: 80, Y: 25},
			},
		},
		{
			name:      "menu event",
			eventType: MenuEventType,
			eventData: func() EventUnion {
				var data EventUnion
				binary.LittleEndian.PutUint32(data[0:4], 100) // CommandID = 100
				return data
			}(),
			expected: MenuEventRecord{
				CommandID: 100,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ir := InputRecord{
				EventType: tt.eventType,
				Event:     tt.eventData,
			}
			result := ir.Unwrap()
			
			switch expected := tt.expected.(type) {
			case FocusEventRecord:
				if focus, ok := result.(FocusEventRecord); !ok || focus != expected {
					t.Errorf("Unwrap() = %v, want %v", result, expected)
				}
			case KeyEventRecord:
				if key, ok := result.(KeyEventRecord); !ok || key != expected {
					t.Errorf("Unwrap() = %v, want %v", result, expected)
				}
			case MouseEventRecord:
				if mouse, ok := result.(MouseEventRecord); !ok || mouse != expected {
					t.Errorf("Unwrap() = %v, want %v", result, expected)
				}
			case WindowBufferSizeEventRecord:
				if win, ok := result.(WindowBufferSizeEventRecord); !ok || win != expected {
					t.Errorf("Unwrap() = %v, want %v", result, expected)
				}
			case MenuEventRecord:
				if menu, ok := result.(MenuEventRecord); !ok || menu != expected {
					t.Errorf("Unwrap() = %v, want %v", result, expected)
				}
			}
		})
	}
}

func TestInputRecordUnwrapUnknownEvent(t *testing.T) {
	ir := InputRecord{
		EventType: EventType(999), // Unknown event type
		Event:     EventUnion{},
	}
	
	result := ir.Unwrap()
	if unknown, ok := result.(*UnknownEvent); !ok {
		t.Errorf("Expected UnknownEvent, got %T", result)
	} else if unknown.EventType != 999 {
		t.Errorf("Expected EventType 999, got %d", unknown.EventType)
	}
}

func TestInputRecordString(t *testing.T) {
	ir := InputRecord{
		EventType: FocusEventType,
		Event: func() EventUnion {
			var data EventUnion
			data[0] = 1
			return data
		}(),
	}
	
	result := ir.String()
	if !strings.Contains(result, "FocusEvent") {
		t.Errorf("Expected string to contain 'FocusEvent', got %q", result)
	}
}

func TestFocusEventRecord(t *testing.T) {
	event := FocusEventRecord{SetFocus: true}
	
	if event.Type() != "FocusEvent" {
		t.Errorf("Type() = %q, want %q", event.Type(), "FocusEvent")
	}
	
	str := event.String()
	if !strings.Contains(str, "FocusEvent") || !strings.Contains(str, "true") {
		t.Errorf("String() = %q, want to contain 'FocusEvent' and 'true'", str)
	}
}

func TestKeyEventRecord(t *testing.T) {
	event := KeyEventRecord{
		KeyDown:         true,
		RepeatCount:     1,
		VirtualKeyCode:  65,
		VirtualScanCode: 30,
		Char:            'A',
		ControlKeyState: SHIFT_PRESSED,
	}
	
	if event.Type() != "KeyEvent" {
		t.Errorf("Type() = %q, want %q", event.Type(), "KeyEvent")
	}
	
	str := event.String()
	expected := []string{"KeyEvent", "'A'", "down", "Shift", "KeyCode: 65", "ScanCode: 30"}
	for _, exp := range expected {
		if !strings.Contains(str, exp) {
			t.Errorf("String() = %q, want to contain %q", str, exp)
		}
	}
}

func TestMouseEventRecord(t *testing.T) {
	event := MouseEventRecord{
		MousePositon:    Coord{X: 10, Y: 20},
		ButtonState:     FROM_LEFT_1ST_BUTTON_PRESSED,
		ControlKeyState: NO_CONTROL_KEY,
		EventFlags:      CLICK,
		WheelDirection:  0,
	}
	
	if event.Type() != "MouseEvent" {
		t.Errorf("Type() = %q, want %q", event.Type(), "MouseEvent")
	}
	
	str := event.String()
	expected := []string{"MouseEvent", "(10, 20)", "Left", "Click"}
	for _, exp := range expected {
		if !strings.Contains(str, exp) {
			t.Errorf("String() = %q, want to contain %q", str, exp)
		}
	}
}

func TestMouseEventRecordWheelDirection(t *testing.T) {
	tests := []struct {
		name           string
		eventFlags     EventFlags
		wheelDirection int
		expected       string
	}{
		{
			name:           "wheel forward",
			eventFlags:     MOUSE_WHEELED,
			wheelDirection: 1,
			expected:       "Forward",
		},
		{
			name:           "wheel backward",
			eventFlags:     MOUSE_WHEELED,
			wheelDirection: -1,
			expected:       "Backward",
		},
		{
			name:           "horizontal wheel right",
			eventFlags:     MOUSE_HWHEELED,
			wheelDirection: 1,
			expected:       "Right",
		},
		{
			name:           "horizontal wheel left",
			eventFlags:     MOUSE_HWHEELED,
			wheelDirection: -1,
			expected:       "Left",
		},
		{
			name:           "no wheel event",
			eventFlags:     CLICK,
			wheelDirection: 0,
			expected:       "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := MouseEventRecord{
				EventFlags:     tt.eventFlags,
				WheelDirection: tt.wheelDirection,
			}
			result := event.WheelDirectionName()
			if result != tt.expected {
				t.Errorf("WheelDirectionName() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestCoordString(t *testing.T) {
	coord := Coord{X: 42, Y: 24}
	expected := "(42, 24)"
	result := coord.String()
	
	if result != expected {
		t.Errorf("String() = %q, want %q", result, expected)
	}
}

func TestButtonState(t *testing.T) {
	tests := []struct {
		name     string
		state    ButtonState
		expected string
	}{
		{"left button", FROM_LEFT_1ST_BUTTON_PRESSED, "Left"},
		{"right button", RIGHTMOST_BUTTON_PRESSED, "Right"},
		{"second button", FROM_LEFT_2ND_BUTTON_PRESSED, "2"},
		{"third button", FROM_LEFT_3RD_BUTTON_PRESSED, "3"},
		{"fourth button", FROM_LEFT_4TH_BUTTON_PRESSED, "4"},
		{"no button", ButtonState(0), "No Button"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.state.String()
			if result != tt.expected {
				t.Errorf("String() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestButtonStateContains(t *testing.T) {
	state := FROM_LEFT_1ST_BUTTON_PRESSED | FROM_LEFT_2ND_BUTTON_PRESSED
	
	if !state.Contains(FROM_LEFT_1ST_BUTTON_PRESSED) {
		t.Error("Expected state to contain FROM_LEFT_1ST_BUTTON_PRESSED")
	}
	
	if !state.Contains(FROM_LEFT_2ND_BUTTON_PRESSED) {
		t.Error("Expected state to contain FROM_LEFT_2ND_BUTTON_PRESSED")
	}
	
	if state.Contains(RIGHTMOST_BUTTON_PRESSED) {
		t.Error("Expected state to not contain RIGHTMOST_BUTTON_PRESSED")
	}
}

func TestControlKeyStateString(t *testing.T) {
	tests := []struct {
		name     string
		state    ControlKeyState
		expected string
	}{
		{"no control keys", NO_CONTROL_KEY, ""},
		{"caps lock", CAPSLOCK_ON, "CapsLock"},
		{"shift", SHIFT_PRESSED, "Shift"},
		{"ctrl", LEFT_CTRL_PRESSED, "CTRL"},
		{"alt", LEFT_ALT_PRESSED, "Alt"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.state.String()
			if result != tt.expected {
				t.Errorf("String() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestControlKeyStateContains(t *testing.T) {
	state := SHIFT_PRESSED | LEFT_CTRL_PRESSED
	
	if !state.Contains(SHIFT_PRESSED) {
		t.Error("Expected state to contain SHIFT_PRESSED")
	}
	
	if !state.Contains(LEFT_CTRL_PRESSED) {
		t.Error("Expected state to contain LEFT_CTRL_PRESSED")
	}
	
	if state.Contains(LEFT_ALT_PRESSED) {
		t.Error("Expected state to not contain LEFT_ALT_PRESSED")
	}
}

func TestEventFlagsString(t *testing.T) {
	tests := []struct {
		name     string
		flags    EventFlags
		expected string
	}{
		{"click", CLICK, "Click"},
		{"double click", DOUBLE_CLICK, "DoubleClick"},
		{"mouse moved", MOUSE_MOVED, "Moved"},
		{"mouse wheeled", MOUSE_WHEELED, "Wheeled"},
		{"mouse hwheeled", MOUSE_HWHEELED, "HWheeld"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.flags.String()
			if result != tt.expected {
				t.Errorf("String() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestEventFlagsContains(t *testing.T) {
	flags := MOUSE_WHEELED | DOUBLE_CLICK
	
	if !flags.Contains(MOUSE_WHEELED) {
		t.Error("Expected flags to contain MOUSE_WHEELED")
	}
	
	if !flags.Contains(DOUBLE_CLICK) {
		t.Error("Expected flags to contain DOUBLE_CLICK")
	}
	
	if flags.Contains(MOUSE_MOVED) {
		t.Error("Expected flags to not contain MOUSE_MOVED")
	}
}