package utils

import (
	"encoding/binary"
	"github.com/google/uuid"
)

func UUId() string {
	uuid, err := uuid.NewUUID()
	if err != nil {
		return ""
	}
	return uuid.String()
}

// 利用uuid生成int64的全局id，
func UUidToInt64ID() int64 {
	uuidObj, _ := uuid.NewUUID()
	bytes, _ := uuidObj.MarshalBinary()
	int64ID := int64(binary.BigEndian.Uint64(bytes[8:]))
	return int64ID
}
