package hocon

import (
	"reflect"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewStack(t *testing.T) {
	tests := []struct {
		name string
		want *Stack
	}{
		{
			name: "creates empty stack correctly",
			want: &Stack{sync.Mutex{}, make([]int, 0)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewStack(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStack() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStack_PushParallelReallyPushes(t *testing.T) {
	loops := 300
	stack := NewStack()

	var wg sync.WaitGroup

	for i := 0; i < loops; i++ {
		wg.Add(1)
		go func(value int) {
			stack.Push(value)
			wg.Done()
		}(i)
	}

	wg.Wait()

	for i := 0; i < loops; i++ {
		_, err := stack.Pop()
		assert.Nil(t, err)
	}

	_, err := stack.Pop()
	assert.Error(t, err)
}
