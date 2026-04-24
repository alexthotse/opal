package domain

// State represents the pure data model of our application.
type State struct {
	Messages  []string
	RequestID int
	ThemeType string
}

// InitialState returns a pure initial state.
func InitialState() State {
	return State{
		Messages:  []string{"Starting Falcon Backend (Gleam)..."},
		RequestID: 0,
		ThemeType: "pi.dev",
	}
}

// AddMessage is a pure function that returns a new state with an appended message.
func AddMessage(s State, msg string) State {
	newMessages := make([]string, len(s.Messages), len(s.Messages)+1)
	copy(newMessages, s.Messages)
	newMessages = append(newMessages, msg)
	return State{
		Messages:  newMessages,
		RequestID: s.RequestID,
		ThemeType: s.ThemeType,
	}
}

// IncrementRequestID returns a new state with an incremented RequestID.
func IncrementRequestID(s State) State {
	return State{
		Messages:  s.Messages,
		RequestID: s.RequestID + 1,
		ThemeType: s.ThemeType,
	}
}

// SetTheme returns a new state with the given theme type.
func SetTheme(s State, themeType string) State {
	return State{
		Messages:  s.Messages,
		RequestID: s.RequestID,
		ThemeType: themeType,
	}
}
