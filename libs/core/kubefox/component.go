package kubefox

import (
	"fmt"
)

type ComponentReg struct {
	Routes []*Route `json:"routes"`
}

func (lhs *Component) Equal(rhs *Component) bool {
	if rhs == nil {
		return false
	}
	return lhs.Name == rhs.Name &&
		lhs.Commit == rhs.Commit &&
		lhs.Id == rhs.Id
}

func (c *Component) Key() string {
	return fmt.Sprintf("%s-%s-%s", c.Name, c.ShortCommit(), c.Id)
}

func (c *Component) GroupKey() string {
	return fmt.Sprintf("%s-%s", c.Name, c.ShortCommit())
}

func (c *Component) Subject() string {
	if c.Id == "" {
		return c.GroupSubject()
	}
	return fmt.Sprintf("evt.js.%s.%s.%s", c.Name, c.ShortCommit(), c.Id)
}

func (c *Component) GroupSubject() string {
	return fmt.Sprintf("evt.js.%s.%s", c.Name, c.ShortCommit())
}

// DirectSubject returns the name of the JetStream subject that Events sent
// directly from Broker to Component should be placed so they are accessible for
// replay and lookup. Use of this subject is not required if Events are sent
// using JetStream as they will be available on that subject.
func (c *Component) DirectSubject() string {
	if c.Id == "" {
		return fmt.Sprintf("evt.direct.%s.%s", c.Name, c.ShortCommit())
	}
	return fmt.Sprintf("evt.direct.%s.%s.%s", c.Name, c.ShortCommit(), c.Id)
}

func (c *Component) ShortCommit() string {
	if len(c.Commit) >= 7 {
		return c.Commit[0:7]
	}

	return ""
}