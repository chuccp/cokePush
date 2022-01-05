package user



type Store struct {
	masterMap *indexMap
}

func NewStore() *Store {
	return &Store{masterMap: newIndexMap()}
}
func (store *Store) AddUser(user IUser) {
	store.masterMap.add(user)
}
func (store *Store) GetUser(username string,f func(IUser)bool) bool {
	return store.masterMap.each(username,f)
}
func (store *Store) GetUserNum() uint {
	return store.masterMap.num
}