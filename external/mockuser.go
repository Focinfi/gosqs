package external

import "sync"

var mux sync.RWMutex

// UserStoreMock mock a user store with a map
type UserStoreMock map[string]int64

// GetUserIDByUniqueID implements UserStore
func (u UserStoreMock) GetUserIDByUniqueID(uniqueID string) (int64, error) {
	mux.RLock()
	id, ok := u[uniqueID]
	if !ok {
		mux.RUnlock()
		return u.createUserByUniqueID(uniqueID)
	}

	mux.RUnlock()
	return id, nil
}

// CreateUserByUniqueID implements UserStore
func (u UserStoreMock) CreateUserByUniqueID(uniqueID string) (int64, error) {
	if id, err := u.GetUserIDByUniqueID(uniqueID); err == nil {
		return id, nil
	}

	return u.createUserByUniqueID(uniqueID)
}

func (u UserStoreMock) createUserByUniqueID(uniqueID string) (int64, error) {
	mux.Lock()
	mux.Unlock()

	id := int64(len(u) + 1)
	u[uniqueID] = id
	return id, nil
}

// DefaultUserStoreMock default mock UserStore
var DefaultUserStoreMock = UserStoreMock{"test": 1}
