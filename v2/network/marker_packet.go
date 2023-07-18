package network

import "encoding/binary"

type MarkerPacket struct {
	packet     Packet
	sessionCtx *SessionContext
	//length     uint16
	//packetType PacketType
	//flag       uint8
	markerData uint8
	markerType uint8
}

func (pck *MarkerPacket) Bytes() []byte {
	if pck.sessionCtx.handshakeComplete && pck.sessionCtx.Version >= 315 {
		return []byte{0, 0x0, 0, 0xB, 0xC, 0, 0, 0, pck.markerType, 0, pck.markerData}
	} else {
		return []byte{0, 0xB, 0, 0, 0xC, 0, 0, 0, pck.markerType, 0, pck.markerData}
	}
}

func (pck *MarkerPacket) GetPacketType() PacketType {
	return pck.packet.PacketType
}

func newMarkerPacket(markerData uint8, sessionCtx *SessionContext) *MarkerPacket {
	return &MarkerPacket{
		packet: Packet{
			DataOffset: 0,
			Length:     0xB,
			PacketType: MARKER,
			Flag:       0x20,
		},
		sessionCtx: sessionCtx,
		markerType: 1,
		markerData: markerData,
	}
}
func newMarkerPacketFromData(packetData []byte, sessionCtx *SessionContext) *MarkerPacket {
	if len(packetData) != 0xB {
		return nil
	}
	pck := MarkerPacket{
		packet: Packet{
			DataOffset: 0,
			PacketType: PacketType(packetData[4]),
			Flag:       packetData[5],
		},
		sessionCtx: sessionCtx,
		markerType: packetData[8],
		markerData: packetData[10],
	}
	if sessionCtx.handshakeComplete && sessionCtx.Version >= 315 {
		pck.packet.Length = binary.BigEndian.Uint32(packetData)
	} else {
		pck.packet.Length = uint32(binary.BigEndian.Uint16(packetData))
	}
	if pck.packet.PacketType != MARKER {
		return nil
	}
	return &pck
}
