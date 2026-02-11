package domain

// UserService defines the business logic contract for user operations.
type UserService interface {
	Create(user *User) error
	Get(id uint) (*User, error)
	Update(user *User) error
	Delete(id uint) error
	List(page, pageSize int) ([]User, int64, error)
}
