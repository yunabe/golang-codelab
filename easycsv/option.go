package easycsv

import "errors"

// Option specifies the spec of Reader.
type Option struct {
	Comma   rune
	Comment rune
	// TODO: Use AutoIndex
	AutoIndex bool
	// TODO: Use AutoName
	AutoName bool
}

func (a *Option) mergeOption(b Option) {
	if b.Comma != 0 {
		a.Comma = b.Comma
	}
	if b.Comment != 0 {
		a.Comment = b.Comment
	}
	if b.AutoIndex {
		a.AutoIndex = true
	}
	if b.AutoName {
		a.AutoName = true
	}
}

func (a *Option) validate() error {
	if a.AutoIndex && a.AutoName {
		return errors.New("You can not set both AutoIndex and AutoName to easycsv.Reader.")
	}
	return nil
}

func mergeOptions(opts []Option) (Option, error) {
	var opt Option
	for _, o := range opts {
		opt.mergeOption(o)
	}
	return opt, opt.validate()
}
