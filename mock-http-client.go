package httphelpers

import (
	"fmt"
	"net/http"
	"sync"
	"testing"
)

type MockHttpClient struct {
	mutex     *sync.Mutex
	callIndex int
	t         *testing.T
	Expects   []MockHttpClientDoExpect
	Calls     []MockHttpClientDoCall
}

func NewMockHttpClient(t *testing.T) *MockHttpClient {
	return &MockHttpClient{
		mutex:   &sync.Mutex{},
		t:       t,
		Expects: []MockHttpClientDoExpect{},
		Calls:   []MockHttpClientDoCall{},
	}
}

type MockHttpClientDoCall struct {
	Request  *http.Request
	Response *MockHttpClientDoExpect
}

type MockHttpClientDoExpect struct {
	WantResponse *http.Response
	WantError    error
}

func (m *MockHttpClient) Do(req *http.Request) (*http.Response, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.callIndex >= len(m.Expects) {
		m.t.Errorf("'Do' called more times than expected")
		return nil, fmt.Errorf("expected %d calls, got %d", len(m.Expects), m.callIndex+1)
	}

	expected := m.Expects[m.callIndex]

	m.Calls = append(m.Calls, MockHttpClientDoCall{
		Request:  req,
		Response: &expected,
	})

	m.callIndex++

	return expected.WantResponse, expected.WantError
}

func (m *MockHttpClient) OnDo(wantResponse *http.Response, wantErr error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.Expects = append(m.Expects, MockHttpClientDoExpect{WantResponse: wantResponse, WantError: wantErr})
}

func (m *MockHttpClient) VerifyCallCount() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if len(m.Expects) != m.callIndex {
		m.t.Errorf("expected %d calls, but got %d", len(m.Expects), m.callIndex)
	}
}
