package network

import (
	"encoding/binary"
	"fmt"
)

// type AcceptPacket Packet
type AcceptPacket struct {
	Packet     Packet
	SessionCtx SessionContext
	Buffer     []byte
}

func (pck *AcceptPacket) Bytes() []byte {
	// ptkSize := 41
	// if pck.sessionCtx.Version < 315 {
	// 	ptkSize = 32
	// }
	output := pck.Packet.Bytes()
	//output := make([]byte, pck.dataOffset)
	//binary.BigEndian.PutUint16(output[0:], pck.packet.length)
	//output[4] = uint8(pck.packet.packetType)
	//output[5] = pck.packet.flag
	binary.BigEndian.PutUint16(output[8:], pck.SessionCtx.Version)
	binary.BigEndian.PutUint16(output[10:], pck.SessionCtx.Options)
	if pck.SessionCtx.Version < 315 {
		binary.BigEndian.PutUint16(output[12:], uint16(pck.SessionCtx.SessionDataUnit))
		binary.BigEndian.PutUint16(output[14:], uint16(pck.SessionCtx.TransportDataUnit))
	} else {
		binary.BigEndian.PutUint32(output[32:], pck.SessionCtx.SessionDataUnit)
		binary.BigEndian.PutUint32(output[36:], pck.SessionCtx.TransportDataUnit)
	}

	binary.BigEndian.PutUint16(output[16:], pck.SessionCtx.Histone)
	binary.BigEndian.PutUint16(output[18:], uint16(len(pck.Buffer)))
	binary.BigEndian.PutUint16(output[20:], pck.Packet.DataOffset)
	output[22] = pck.SessionCtx.ACFL0
	output[23] = pck.SessionCtx.ACFL1
	// s
	output = append(output, pck.Buffer...)
	return output
}
func (pck *AcceptPacket) GetPacketType() PacketType {
	return pck.Packet.PacketType
}

//func NewAcceptPacket(sessionCtx SessionContext, acceptData []byte) *AcceptPacket {
//	sessionCtx.Histone = 1
//	sessionCtx.ACFL0 = 4
//	sessionCtx.ACFL1 = 4
//	pck := AcceptPacket{
//		sessionCtx: sessionCtx,
//		dataOffset: 32,
//		length:        0,
//		packetType:       2,
//		flag:       0,
//		NSPFSID:    0,
//		buffer:     acceptData,
//		SID:        nil,
//	}
//	if len(acceptData) > 230 {
//		pck.length = uint16(len(acceptData)) + pck.dataOffset
//	}
//	return &pck
//}

func NewAcceptPacketFromData(packetData []byte, connOption *ConnectionOption) *AcceptPacket {
	if len(packetData) < 32 {
		return nil
	}
	reconAddStart := binary.BigEndian.Uint16(packetData[28:])
	reconAddLen := binary.BigEndian.Uint16(packetData[30:])
	reconAdd := ""
	if reconAddStart != 0 && reconAddLen != 0 && uint16(len(packetData)) > (reconAddStart+reconAddLen) {
		reconAdd = string(packetData[reconAddStart:(reconAddStart + reconAddLen)])
	}
	pck := AcceptPacket{
		Packet: Packet{
			DataOffset: binary.BigEndian.Uint16(packetData[20:]),
			Length:     uint32(binary.BigEndian.Uint16(packetData)),
			PacketType: PacketType(packetData[4]),
			Flag:       packetData[5],
		},
		SessionCtx: SessionContext{
			ConnOption:          connOption,
			SID:                 nil,
			Version:             binary.BigEndian.Uint16(packetData[8:]),
			LoVersion:           0,
			Options:             0,
			NegotiatedOptions:   binary.BigEndian.Uint16(packetData[10:]),
			OurOne:              0,
			Histone:             binary.BigEndian.Uint16(packetData[16:]),
			ReconAddr:           reconAdd,
			ACFL0:               packetData[22],
			ACFL1:               packetData[23],
			SessionDataUnit:     uint32(binary.BigEndian.Uint16(packetData[12:])),
			TransportDataUnit:   uint32(binary.BigEndian.Uint16(packetData[14:])),
			UsingAsyncReceivers: false,
			IsNTConnected:       false,
			OnBreakReset:        false,
			GotReset:            false,
		},
	}
	pck.Buffer = packetData[int(pck.Packet.DataOffset):]
	if pck.SessionCtx.Version >= 315 {
		pck.SessionCtx.SessionDataUnit = binary.BigEndian.Uint32(packetData[32:])
		pck.SessionCtx.TransportDataUnit = binary.BigEndian.Uint32(packetData[36:])
	}
	if (pck.Packet.Flag & 1) > 0 {
		fmt.Println("contain SID data")
		pck.Packet.Length -= 16
		pck.SessionCtx.SID = packetData[int(pck.Packet.Length):]
	}
	if pck.SessionCtx.TransportDataUnit < pck.SessionCtx.SessionDataUnit {
		pck.SessionCtx.SessionDataUnit = pck.SessionCtx.TransportDataUnit
	}
	if binary.BigEndian.Uint16(packetData[18:]) != uint16(len(pck.Buffer)) {
		return nil
	}
	return &pck
}

//func (pck *AcceptPacket) SessionCTX() SessionContext {
//	return pck.sessionCtx
//}
