package infrastructure

type LinkRepository struct{}

func NewLinkRepository() *LinkRepository {
	return &LinkRepository{}
}

func (r *LinkRepository) GetShortLink() {
	r.storage.
}

func (r *LinkRepository) CreateShortLink() {}

func (r *LinkRepository) BatchCreateLinks() {}

func (r *LinkRepository) GetUserLinks() {}
