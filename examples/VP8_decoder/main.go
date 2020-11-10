package main

import (
	"../../cgo/ffmpeg"
	"fmt"
	"strconv"
	"./muxers"
	"time"
	"image"
	"image/png"
	"os"
)
func main() {
 	decoder,err:=ffmpeg.NewVP8Decoder()
 	fmt.Println("err=",err)
	decoder.Setup()
 	go vreader("5001",decoder)
 	for ;;{
		time.Sleep(time.Second)
	}
}

func vreader(url string,dec *ffmpeg.VideoDecoder) {
	//var part int = 16
	//var buff []byte
	//var prevnumber uint16 = 0
	//audioUdpSource := muxers.NewUdpSource(5002)
	port, _ := strconv.Atoi(url)
	videoUdpSource := muxers.NewUdpSource(port)
	demuxer:=muxers.NewRtpDemuxer()
	vp8Depaketizer:=muxers.NewRtpVP8Depacketizer()
	//depacketizer:=muxers.NewRtpH264Depacketizer()
	muxers.Bridge(videoUdpSource.OutputChan,demuxer.InputChan)
	muxers.Bridge(demuxer.OutputChan,vp8Depaketizer.InputChan)
	//muxers.Bridge(demuxer.OutputChan,depacketizer.InputChan)
	//reader, writer := io.Pipe()


	//var frameProcessed bool=false

	//var buf bytes.Buffer
	var framenum int=0
	for ; ; {
		data := (<-vp8Depaketizer.OutputChan).(*muxers.RtpPacket)
		fmt.Println("VIDEO DATA~~~~~~~~~~~~~~~~~~~~~")
		fr,err:=dec.Decode(data.Payload)
		if err==nil{
			str:=fmt.Sprintf("FRAME%d",framenum)
			encodeImg(fr.Image,str)
			framenum++
		}

		fmt.Println("~~~~~~~~~~~~~~~~~~~~~")

	}
}
func encodeImg(img image.YCbCr, framenum string) {
	bounds := img.Bounds()
	bimj := image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy()))
	for a := 0; a < bounds.Max.X; a++ {
		for b := 0; b < bounds.Max.Y; b++ {
			bimj.Set(a, b, img.At(a, b))
		}
	}
	filepng, _ := os.Create(fmt.Sprintf("frame%s.png", framenum))
	png.Encode(filepng, bimj)
}
