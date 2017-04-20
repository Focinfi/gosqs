package external

// UserStoreMock mock a user store with a map
type UserStoreMock map[string]int64

// GetUserIDByUniqueID implements UserStore
func (u UserStoreMock) GetUserIDByUniqueID(uniqueID string) (int64, error) {
	return u[uniqueID], nil
}

// CreateUserByUniqueID implements UserStore
func (u UserStoreMock) CreateUserByUniqueID(uniqueID string) (int64, error) {
	id := int64(len(u) + 1)
	u[uniqueID] = id
	return id, nil
}

// DefaultUserStoreMock default mock UserStore
var DefaultUserStoreMock = UserStoreMock{}
