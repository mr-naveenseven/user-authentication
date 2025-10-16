package user

type UserRepoPort interface {
	Create(user User) (User, error)
	GetByID(userID int) (User, error)
	Get() ([]User, error)
	GetByUsername(username string) (User, error)
	Update(user User) (User, error)
}
