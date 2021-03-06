// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package mocks

import (
	"github.com/jmgilman/gcli/ui"
	"sync"
)

var (
	lockPrompterMockRun sync.RWMutex
)

// Ensure, that PrompterMock does implement ui.Prompter.
// If this is not the case, regenerate this file with moq.
var _ ui.Prompter = &PrompterMock{}

// PrompterMock is a mock implementation of ui.Prompter.
//
//     func TestSomethingThatUsesPrompter(t *testing.T) {
//
//         // make and configure a mocked ui.Prompter
//         mockedPrompter := &PrompterMock{
//             RunFunc: func() (string, error) {
// 	               panic("mock out the Run method")
//             },
//         }
//
//         // use mockedPrompter in code that requires ui.Prompter
//         // and then make assertions.
//
//     }
type PrompterMock struct {
	// RunFunc mocks the Run method.
	RunFunc func() (string, error)

	// calls tracks calls to the methods.
	calls struct {
		// Run holds details about calls to the Run method.
		Run []struct {
		}
	}
}

// Run calls RunFunc.
func (mock *PrompterMock) Run() (string, error) {
	if mock.RunFunc == nil {
		panic("PrompterMock.RunFunc: method is nil but Prompter.Run was just called")
	}
	callInfo := struct {
	}{}
	lockPrompterMockRun.Lock()
	mock.calls.Run = append(mock.calls.Run, callInfo)
	lockPrompterMockRun.Unlock()
	return mock.RunFunc()
}

// RunCalls gets all the calls that were made to Run.
// Check the length with:
//     len(mockedPrompter.RunCalls())
func (mock *PrompterMock) RunCalls() []struct {
} {
	var calls []struct {
	}
	lockPrompterMockRun.RLock()
	calls = mock.calls.Run
	lockPrompterMockRun.RUnlock()
	return calls
}
