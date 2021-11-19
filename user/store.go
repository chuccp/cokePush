package user

type Store struct {
	masterMap *indexMap
}

func NewStore() *Store {
	return &Store{masterMap:newIndexMap()}
}
func (store *Store)AddUser(user IUser)  {
	store.masterMap.add(user)
}
func (store *Store)GetUser(username string) IUser {
	return store.masterMap.get(username)
}