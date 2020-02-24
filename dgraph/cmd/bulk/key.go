package bulk

import (
	"fmt"
	"io"
)

func GetKeyForMapEntry(dAtA []byte) ([]byte, error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				panic("interger overflow")
			}
			if iNdEx >= l {
				return nil, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return nil, fmt.Errorf("proto: MapEntry: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return nil, fmt.Errorf("proto: MapEntry: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return nil, fmt.Errorf("proto: wrong wireType = %d for field Key", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					panic("interger overflow")
				}
				if iNdEx >= l {
					return nil, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				panic("interger overflow")
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				panic("interger overflow")
			}
			if postIndex > l {
				return nil, io.ErrUnexpectedEOF
			}
			return dAtA[iNdEx:postIndex], nil
		default:
			panic("unknown filed number")
		}
	}

	if iNdEx > l {
		return nil, io.ErrUnexpectedEOF
	}
	return nil, io.ErrUnexpectedEOF
}
