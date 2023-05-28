package payload

func (rec *Record) Clone() *Record {
	clone := &Record{
		Id:       rec.Id,
		Request:  rec.Request.Clone(),
		Response: rec.Response.Clone(),
		Step:     rec.Step,
	}
	return clone
}
