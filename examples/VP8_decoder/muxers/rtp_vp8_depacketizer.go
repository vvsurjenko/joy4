package muxers

import "fmt"

type RtpVP8Depacketizer struct {
	InputChan chan interface {}
	OutputChan chan interface {}
}

func NewRtpVP8Depacketizer() *RtpH264Depacketizer {
	demuxer := &RtpH264Depacketizer{
		InputChan: make(chan interface {}),
		OutputChan: make(chan interface {}),
	}

	go func() {
		var buffer[] byte
		for {
			packet := (<-demuxer.InputChan).(*RtpPacket)
			header := packet.Payload[0]
			//fmt.Println("HEADER:",packet.Payload[:3],"SOME NEXT BYTES",packet.Payload[3:6])
			if (header & 127)>>4==1{ //first packet
				fmt.Println("144")
				if len(buffer)!=0{
					Payload := make([]byte, 0)
					Payload = append(Payload, buffer...)
					buffer=make([]byte, 0)
					buffer = append(buffer, packet.Payload[3:]...)

					packet.Payload = Payload
					demuxer.OutputChan <-packet
				} else {
					buffer = append(buffer, packet.Payload[3:]...)
				}

			} else { //additional packets
				fmt.Println("128")
				buffer = append(buffer, packet.Payload[3:]...)
			}
		}
	}()

	return demuxer
}
