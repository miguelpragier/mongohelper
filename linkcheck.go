package mongohelper

import "fmt"

func (l Link) linkCheck(routine string) error {
	if l.client == nil {
		l.log(routine, "use of uninitialized connection")

		return fmt.Errorf("use of uninitialized connection")
	}

	return nil
}
