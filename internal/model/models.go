package model

type Models struct {
	Lookup ILookup
}

func NewModels(lu ILookup) *Models {
	return &Models{
		Lookup: lu,
	}
}
