package errs

func (ep *ErrPlace) Reset() {
	if ep == nil {
		return
	}

	ep.Message = ""

	ep.File = ""

	ep.Line = 0

}
