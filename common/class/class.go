package class

type Class interface {
	Inherit(Class)
}

type DefaultClass struct {
	derived Class
}

func (this *DefaultClass) Inherit(that Class) {
	this.derived = that
}
