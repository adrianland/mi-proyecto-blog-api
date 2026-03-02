package domain

type ArticleRepository interface {
	Create(article *Article) error
	GetByID(id int) (*Article, error)
	Update(article *Article) error
	ListPublished(page, pageSize int) ([]ArticleWithAuthor, int64, error)
	ListByAuthor(authorID int, status string, page, pageSize int) ([]Article, int64, error)
	CountPublishedByAuthor(authorID int) (int, error)
	GetAllPublished() ([]Article, error)
}
