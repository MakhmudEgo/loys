package api

func checkErrAndAppend(errs []error, err error) []error {
	if err != nil {
		errs = append(errs, err)
	}

	return errs
}
