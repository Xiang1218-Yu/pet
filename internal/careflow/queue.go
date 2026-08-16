package careflow

func accepted(c Checkin) bool { return c.Pet != "" && !c.At.IsZero() }
