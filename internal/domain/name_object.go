package domain

type NameObject struct {
	value string
}

func NewNameObject(name string) NameObject {
	return NameObject{value: name}
}

func (o *NameObject) Value() string {
	return o.value
}

func (o *NameObject) IsGuest() bool {
	return o.value == ""
}
