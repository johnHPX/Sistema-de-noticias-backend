package noticiacategoria

type Repository interface {
	Store(e *Entity) error
	List() (*[]Entity, error)
}
