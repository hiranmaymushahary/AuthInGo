package repositories

type Storage struct { //FACILITATES DEPENDENCY INJECTION FOR REPOSITORY
	UserRepository UserRepository
}

func NewStorage() *Storage {
	return &Storage{
		UserRepository: &UserRepositoryImpl{},
	}
}
