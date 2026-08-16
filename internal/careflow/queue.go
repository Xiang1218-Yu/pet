package careflow

func usable(e Entry) bool { return e.ID != "" && !e.At.IsZero() }
