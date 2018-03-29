package zcloud

type Bucket interface {
	Create () error
	Delete () error
	Name () string
	Object (key string) Object
	Objects () ([]Object, error)
	ObjectsQuery (query *ObjectsQueryParams) ([]Object, error)
}
