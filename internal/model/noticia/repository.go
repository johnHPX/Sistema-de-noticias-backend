package noticia

type Repository interface {
	Store(e *Entity) error
	List() (*[]Entity, error)
	// ListByTitOrCat(titCat string) ([]Entity, error)
	// Update(e *Entity) error
	// Remove(nid string) error
}
