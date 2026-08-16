package careflow

func valid(d Dose) bool { return d.Units >= 0 && d.Medicine != "" }
