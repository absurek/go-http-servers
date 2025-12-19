package request

import "github.com/google/uuid"

func ParseOptionalUUID(s string) uuid.NullUUID {
	id, err := uuid.Parse(s)
	return uuid.NullUUID{UUID: id, Valid: err == nil}
}
