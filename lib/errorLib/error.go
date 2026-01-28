package errorlib

type CanErrorFuncNoReturn func() error

func ExecMultipleCanError(funcs ...CanErrorFuncNoReturn) error {
	for _, f := range funcs {
		if err := f(); err != nil {
			return err
		}
	}
	return nil
}
