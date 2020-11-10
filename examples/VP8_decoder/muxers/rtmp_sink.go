package muxers

import (
	"github.com/zhangpeihao/gortmp"
	"github.com/zhangpeihao/log"
	"fmt"
	"time"
)

type RtmpSink struct {
	InputChan chan interface {}
}

func NewRtmpSink(url, name string) *RtmpSink {
	sink := &RtmpSink{
		InputChan: make(chan interface {}),
	}

	go func() {
		var err error

		l := log.NewLogger("logs", "publisher", nil, 60, 3600*24, true)
		gortmp.InitLogger(l)

		handler := &RtmpSinkHandler{}
		handler.createStreamChan = make(chan gortmp.OutboundStream)
		handler.flvChan = sink.InputChan

		handler.obConn, err = gortmp.Dial(url, handler, 100)
		if err != nil {
			fmt.Println("Rtmp dial error", err)
			return
		}
		err = handler.obConn.Connect()
		if err != nil {
			fmt.Printf("Connect error: %s", err.Error())
			return
		}

		for {
			select {
			case stream := <-handler.createStreamChan:
				// Publish
				stream.Attach(handler)
				err = stream.Publish(name, "live")
				if err != nil {
					fmt.Printf("Publish error: %s", err.Error())
					return
				}

			case <-time.After(1 * time.Second):
				fmt.Printf("Audio size: %d bytes; Video size: %d bytes\n", handler.audioDataSize, handler.videoDataSize)
			}
		}
	}()

	return sink
}

type RtmpSinkHandler struct {
	status uint
	obConn gortmp.OutboundConn
	createStreamChan chan gortmp.OutboundStream
	videoDataSize int64
	audioDataSize int64
	flvChan chan interface {}
}

func (handler *RtmpSinkHandler) OnStatus(conn gortmp.OutboundConn) {
	var err error
	handler.status, err = handler.obConn.Status()
	fmt.Printf("@@@@@@@@@@@@@status: %d, err: %v\n", handler.status, err)
}

func (handler *RtmpSinkHandler) OnClosed(conn gortmp.Conn) {
	fmt.Printf("@@@@@@@@@@@@@Closed\n")
}

func (handler *RtmpSinkHandler) OnReceived(conn gortmp.Conn, message *gortmp.Message) {
}

func (handler *RtmpSinkHandler) OnReceivedRtmpCommand(conn gortmp.Conn, command *gortmp.Command) {
	fmt.Printf("ReceviedRtmpCommand: %+v\n", command)
}

func (handler *RtmpSinkHandler) OnStreamCreated(conn gortmp.OutboundConn, stream gortmp.OutboundStream) {
	fmt.Printf("Stream created: %d\n", stream.ID())
	handler.createStreamChan <- stream
}
func (handler *RtmpSinkHandler) OnPlayStart(stream gortmp.OutboundStream) {

}
func (handler *RtmpSinkHandler) OnPublishStart(stream gortmp.OutboundStream) {
	// Set chunk buffer size
	go func() {
		fmt.Println("Publish")
		for {
			flvTag := (<-handler.flvChan).(*FlvTag)
			//fmt.Println("Push data")
			if flvTag.TagType == TAG_VIDEO {
				stream.PublishVideoData(flvTag.Data, flvTag.Timestamp)
			} else {
				stream.PublishAudioData(flvTag.Data, flvTag.Timestamp)
			}
		}
	}()
}
