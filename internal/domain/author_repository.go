package domain

type AuthorRepository interface {
	Create(author *Author) error
	GetByID(id int) (*Author, error)
	GetAll() ([]Author, error)
	GetSummary(authorID int) (*AuthorSummary, error)
}
