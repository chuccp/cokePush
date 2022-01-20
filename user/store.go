package user



type Store struct {
	masterMap *indexMap
}

func NewStore() *Store {
	return &Store{masterMap: NewIndexMap()}
}
func (store *Store) AddUser(user IUser)bool {
	return store.masterMap.add(user)
}

func (store *Store) DeleteUser(user IUser)bool{
	return store.masterMap.delete(user)
}

func (store *Store) GetUser(username string,f func(IUser)bool) bool {
	return store.masterMap.each(username,f)
}
func (store *Store) Has(username string) bool {
	return store.masterMap.has(username)
}
func (store *Store) GetUserNum() int32 {
	return store.masterMap.num
}