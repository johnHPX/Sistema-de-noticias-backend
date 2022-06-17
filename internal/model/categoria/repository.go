package categoria

type Repository interface {
	Store(e *Entity) error
	List() (*[]Entity, error)
	Find(cid string) (*Entity, error)
}
