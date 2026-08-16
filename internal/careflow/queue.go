package careflow

// valid keeps only real doses: a positive unit count and a named medicine.
// Zero (or negative) unit counts are sync noise and must be dropped rather
// than carried into the rebuilt timeline.
func valid(d Dose) bool { return d.Units > 0 && d.Medicine != "" }
