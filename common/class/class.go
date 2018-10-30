package class

type Class interface {
	Inherit()
	GetDerived() Class
	Duplicate() Class
}

type DefaultClass struct {
	derived Class
}

func (this *DefaultClass) Inherit(that Class) {
	this.derived = that
}

func (this *DefaultClass) GetDerived() Class {
	return this.derived
}
