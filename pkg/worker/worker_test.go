package worker

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
)

// Mocks ---

var defMethod = "exec"

type mockFn struct {
	mock.Mock
}

func (m *mockFn) exec() (id string, err error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

// Tests ---

func TestItShouldExecuteFunctionAndReturnTheIdAndErrorNilWhenSuccess(t *testing.T) {
	mk1 := new(mockFn)
	mk1.On(defMethod).Return("01", nil)
	mk2 := new(mockFn)
	mk2.On(defMethod).Return("02", errors.New("OH NO"))
	mk3 := new(mockFn)
	mk3.On(defMethod).Return("03", nil)

	dispatcher := NewDispatcher(10)
	dispatcher.Run()

	dispatcher.Dispatch(mk1.exec)
	dispatcher.Dispatch(mk2.exec)
	dispatcher.Dispatch(mk3.exec)

	dispatcher.Quit()

	mk1.AssertExpectations(t)
	mk1.AssertNumberOfCalls(t, defMethod, 1)
	mk2.AssertExpectations(t)
	mk2.AssertNumberOfCalls(t, defMethod, 1)
	mk2.AssertExpectations(t)
	mk3.AssertNumberOfCalls(t, defMethod, 1)
}

func TestItShouldQuitThePoolEvenWhenThereIsNoJobsRunning(t *testing.T) {
	dispatcher := NewDispatcher(1)
	dispatcher.Run()
	dispatcher.Quit()
}
