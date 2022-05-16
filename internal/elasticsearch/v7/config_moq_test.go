// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package v7

import (
	"sync"
)

// Ensure, that configMock does implement config.
// If this is not the case, regenerate this file with moq.
var _ config = &configMock{}

// configMock is a mock implementation of config.
//
// 	func TestSomethingThatUsesconfig(t *testing.T) {
//
// 		// make and configure a mocked config
// 		mockedconfig := &configMock{
// 			GetAPIKeyFunc: func() string {
// 				panic("mock out the GetAPIKey method")
// 			},
// 			GetCertificateFingerprintFunc: func() string {
// 				panic("mock out the GetCertificateFingerprint method")
// 			},
// 			GetCloudIDFunc: func() string {
// 				panic("mock out the GetCloudID method")
// 			},
// 			GetHostFunc: func() string {
// 				panic("mock out the GetHost method")
// 			},
// 			GetIndexFunc: func() string {
// 				panic("mock out the GetIndex method")
// 			},
// 			GetPasswordFunc: func() string {
// 				panic("mock out the GetPassword method")
// 			},
// 			GetServiceTokenFunc: func() string {
// 				panic("mock out the GetServiceToken method")
// 			},
// 			GetUsernameFunc: func() string {
// 				panic("mock out the GetUsername method")
// 			},
// 		}
//
// 		// use mockedconfig in code that requires config
// 		// and then make assertions.
//
// 	}
type configMock struct {
	// GetAPIKeyFunc mocks the GetAPIKey method.
	GetAPIKeyFunc func() string

	// GetCertificateFingerprintFunc mocks the GetCertificateFingerprint method.
	GetCertificateFingerprintFunc func() string

	// GetCloudIDFunc mocks the GetCloudID method.
	GetCloudIDFunc func() string

	// GetHostFunc mocks the GetHost method.
	GetHostFunc func() string

	// GetIndexFunc mocks the GetIndex method.
	GetIndexFunc func() string

	// GetPasswordFunc mocks the GetPassword method.
	GetPasswordFunc func() string

	// GetServiceTokenFunc mocks the GetServiceToken method.
	GetServiceTokenFunc func() string

	// GetUsernameFunc mocks the GetUsername method.
	GetUsernameFunc func() string

	// calls tracks calls to the methods.
	calls struct {
		// GetAPIKey holds details about calls to the GetAPIKey method.
		GetAPIKey []struct {
		}
		// GetCertificateFingerprint holds details about calls to the GetCertificateFingerprint method.
		GetCertificateFingerprint []struct {
		}
		// GetCloudID holds details about calls to the GetCloudID method.
		GetCloudID []struct {
		}
		// GetHost holds details about calls to the GetHost method.
		GetHost []struct {
		}
		// GetIndex holds details about calls to the GetIndex method.
		GetIndex []struct {
		}
		// GetPassword holds details about calls to the GetPassword method.
		GetPassword []struct {
		}
		// GetServiceToken holds details about calls to the GetServiceToken method.
		GetServiceToken []struct {
		}
		// GetUsername holds details about calls to the GetUsername method.
		GetUsername []struct {
		}
	}
	lockGetAPIKey                 sync.RWMutex
	lockGetCertificateFingerprint sync.RWMutex
	lockGetCloudID                sync.RWMutex
	lockGetHost                   sync.RWMutex
	lockGetIndex                  sync.RWMutex
	lockGetPassword               sync.RWMutex
	lockGetServiceToken           sync.RWMutex
	lockGetUsername               sync.RWMutex
}

// GetAPIKey calls GetAPIKeyFunc.
func (mock *configMock) GetAPIKey() string {
	if mock.GetAPIKeyFunc == nil {
		panic("configMock.GetAPIKeyFunc: method is nil but config.GetAPIKey was just called")
	}
	callInfo := struct {
	}{}
	mock.lockGetAPIKey.Lock()
	mock.calls.GetAPIKey = append(mock.calls.GetAPIKey, callInfo)
	mock.lockGetAPIKey.Unlock()
	return mock.GetAPIKeyFunc()
}

// GetAPIKeyCalls gets all the calls that were made to GetAPIKey.
// Check the length with:
//     len(mockedconfig.GetAPIKeyCalls())
func (mock *configMock) GetAPIKeyCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockGetAPIKey.RLock()
	calls = mock.calls.GetAPIKey
	mock.lockGetAPIKey.RUnlock()
	return calls
}

// GetCertificateFingerprint calls GetCertificateFingerprintFunc.
func (mock *configMock) GetCertificateFingerprint() string {
	if mock.GetCertificateFingerprintFunc == nil {
		panic("configMock.GetCertificateFingerprintFunc: method is nil but config.GetCertificateFingerprint was just called")
	}
	callInfo := struct {
	}{}
	mock.lockGetCertificateFingerprint.Lock()
	mock.calls.GetCertificateFingerprint = append(mock.calls.GetCertificateFingerprint, callInfo)
	mock.lockGetCertificateFingerprint.Unlock()
	return mock.GetCertificateFingerprintFunc()
}

// GetCertificateFingerprintCalls gets all the calls that were made to GetCertificateFingerprint.
// Check the length with:
//     len(mockedconfig.GetCertificateFingerprintCalls())
func (mock *configMock) GetCertificateFingerprintCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockGetCertificateFingerprint.RLock()
	calls = mock.calls.GetCertificateFingerprint
	mock.lockGetCertificateFingerprint.RUnlock()
	return calls
}

// GetCloudID calls GetCloudIDFunc.
func (mock *configMock) GetCloudID() string {
	if mock.GetCloudIDFunc == nil {
		panic("configMock.GetCloudIDFunc: method is nil but config.GetCloudID was just called")
	}
	callInfo := struct {
	}{}
	mock.lockGetCloudID.Lock()
	mock.calls.GetCloudID = append(mock.calls.GetCloudID, callInfo)
	mock.lockGetCloudID.Unlock()
	return mock.GetCloudIDFunc()
}

// GetCloudIDCalls gets all the calls that were made to GetCloudID.
// Check the length with:
//     len(mockedconfig.GetCloudIDCalls())
func (mock *configMock) GetCloudIDCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockGetCloudID.RLock()
	calls = mock.calls.GetCloudID
	mock.lockGetCloudID.RUnlock()
	return calls
}

// GetHost calls GetHostFunc.
func (mock *configMock) GetHost() string {
	if mock.GetHostFunc == nil {
		panic("configMock.GetHostFunc: method is nil but config.GetHost was just called")
	}
	callInfo := struct {
	}{}
	mock.lockGetHost.Lock()
	mock.calls.GetHost = append(mock.calls.GetHost, callInfo)
	mock.lockGetHost.Unlock()
	return mock.GetHostFunc()
}

// GetHostCalls gets all the calls that were made to GetHost.
// Check the length with:
//     len(mockedconfig.GetHostCalls())
func (mock *configMock) GetHostCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockGetHost.RLock()
	calls = mock.calls.GetHost
	mock.lockGetHost.RUnlock()
	return calls
}

// GetIndex calls GetIndexFunc.
func (mock *configMock) GetIndex() string {
	if mock.GetIndexFunc == nil {
		panic("configMock.GetIndexFunc: method is nil but config.GetIndex was just called")
	}
	callInfo := struct {
	}{}
	mock.lockGetIndex.Lock()
	mock.calls.GetIndex = append(mock.calls.GetIndex, callInfo)
	mock.lockGetIndex.Unlock()
	return mock.GetIndexFunc()
}

// GetIndexCalls gets all the calls that were made to GetIndex.
// Check the length with:
//     len(mockedconfig.GetIndexCalls())
func (mock *configMock) GetIndexCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockGetIndex.RLock()
	calls = mock.calls.GetIndex
	mock.lockGetIndex.RUnlock()
	return calls
}

// GetPassword calls GetPasswordFunc.
func (mock *configMock) GetPassword() string {
	if mock.GetPasswordFunc == nil {
		panic("configMock.GetPasswordFunc: method is nil but config.GetPassword was just called")
	}
	callInfo := struct {
	}{}
	mock.lockGetPassword.Lock()
	mock.calls.GetPassword = append(mock.calls.GetPassword, callInfo)
	mock.lockGetPassword.Unlock()
	return mock.GetPasswordFunc()
}

// GetPasswordCalls gets all the calls that were made to GetPassword.
// Check the length with:
//     len(mockedconfig.GetPasswordCalls())
func (mock *configMock) GetPasswordCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockGetPassword.RLock()
	calls = mock.calls.GetPassword
	mock.lockGetPassword.RUnlock()
	return calls
}

// GetServiceToken calls GetServiceTokenFunc.
func (mock *configMock) GetServiceToken() string {
	if mock.GetServiceTokenFunc == nil {
		panic("configMock.GetServiceTokenFunc: method is nil but config.GetServiceToken was just called")
	}
	callInfo := struct {
	}{}
	mock.lockGetServiceToken.Lock()
	mock.calls.GetServiceToken = append(mock.calls.GetServiceToken, callInfo)
	mock.lockGetServiceToken.Unlock()
	return mock.GetServiceTokenFunc()
}

// GetServiceTokenCalls gets all the calls that were made to GetServiceToken.
// Check the length with:
//     len(mockedconfig.GetServiceTokenCalls())
func (mock *configMock) GetServiceTokenCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockGetServiceToken.RLock()
	calls = mock.calls.GetServiceToken
	mock.lockGetServiceToken.RUnlock()
	return calls
}

// GetUsername calls GetUsernameFunc.
func (mock *configMock) GetUsername() string {
	if mock.GetUsernameFunc == nil {
		panic("configMock.GetUsernameFunc: method is nil but config.GetUsername was just called")
	}
	callInfo := struct {
	}{}
	mock.lockGetUsername.Lock()
	mock.calls.GetUsername = append(mock.calls.GetUsername, callInfo)
	mock.lockGetUsername.Unlock()
	return mock.GetUsernameFunc()
}

// GetUsernameCalls gets all the calls that were made to GetUsername.
// Check the length with:
//     len(mockedconfig.GetUsernameCalls())
func (mock *configMock) GetUsernameCalls() []struct {
} {
	var calls []struct {
	}
	mock.lockGetUsername.RLock()
	calls = mock.calls.GetUsername
	mock.lockGetUsername.RUnlock()
	return calls
}
