package domain

type ContentRepository interface {
	Save(content *Content) error
	FindAll() ([]*Content, error)
	FindByID(id string) (*Content, error)
	FindFree() ([]*Content, error)
	FindByPsychologist(psychID string) ([]*Content, error)
}
