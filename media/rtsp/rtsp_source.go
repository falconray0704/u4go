package rtsp

import (
	"fmt"
	"net"
	"crypto/md5"
	"encoding/hex"
	"strings"
	"strconv"
	"time"
	"io"
	"net/url"
	b64 "encoding/base64"

	"github.com/pkg/errors"
	"context"
	"runtime/debug"

	"github.com/falconray0704/fsm"
	"github.com/falconray0704/u4go/media/rtsp/h264parser"
)


type RtspPacket struct {
	NalType byte
	Timestamp uint32
	Payload []byte
}

const RtspDropTriger = 30 * 1
const RtspBuf = RtspDropTriger * 3

const RtspPayloadMax = 1024 * 1024 * 2

var (
	// Rtsp client
	RtspCltOutChLen = 1024 * 4 * 100

	// fsm
	Status_Leave_Prefix string = "leave_"
	Status_Enter_Prefix string = "enter_"
	Event_Before_Prefix string = "before_"
	Event_After_Prefix  string = "after_"

	Status_Stop string = "sStop"
	Status_Play string = "sPlay"

	Event_Stop string = "eStop"
	Event_Play string = "ePlay"

	CBK_Leave_Status_Stop string = Status_Leave_Prefix + Status_Stop
	CBK_Leave_Status_Play string = Status_Leave_Prefix + Status_Play

)

type StartRtspClientResp struct {
	err error
	spsInfo h264parser.SPSInfo
}

type RtspSource struct {

	fsmRunner *fsm.FSM

	rCtxParent context.Context
	rCtxParentCclFunc context.CancelFunc

	rCtx context.Context
	rCtxCclFunc context.CancelFunc

	//rCtxClt context.Context
	//rCtxCltCclFunc context.CancelFunc

	Id uint32
	rtspUrl string
	outGoingStream chan []byte 	//out chanel
	OutputChan chan interface {}
	spsCh	chan h264parser.SPSInfo

	sps []byte
	pps []byte
	syncCount uint64
	naluCount uint64
	dfc			uint64

	RtspReader *RtspClient
	cltExit chan bool
	outPkgExit chan bool

}

func (src *RtspSource)handleNALU(nalType byte, payload []byte, ts uint32) {

	rtpPkg := &RtspPacket{
		NalType: nalType,
		Timestamp: ts,
		Payload: payload }

	if nalType == 7 {
		if len(src.sps) == 0 {
			src.sps = payload
		}
		rtpPkg.Payload = payload
		//fmt.Printf("+++ handleNALU() sps nalType:%x ++++++++++\n", nalType)
		if spsInfo, err := h264parser.ParseSPS(payload); err != nil {
			fmt.Printf("+++ id:%d handleNALU() ParseSPS err:%v +++\n", src.Id, err)
		} else {
			//src.SPSVideoW = uint32(spsInfo.Width)
			//src.SPSVideoH = uint32(spsInfo.Height)
			//fmt.Printf("+++ spsInfo:%+v +++\n",spsInfo)

			select {
			case <-src.spsCh:
				break
			default:
				break
			}

			src.spsCh <- spsInfo
		}
	} else if nalType == 8 {
		if len(src.pps) == 0 {
			src.pps = payload
		}
		rtpPkg.Payload = payload
		//fmt.Printf("+++ handleNALU() pps nalType:%x +++++\n", nalType)
	} else if nalType == 5 {
		// keyframe
		src.syncCount++
		src.naluCount++

		//writeNALU(true, int(ts), payload)
		rtpPkg.Payload = payload
		//log.Printf("=== naluCount:%d get I frame flvChk len:%d ===\n", rtspSource.naluCount - 1, len(rtspSource.OutputChan) )
	} else {
		if nalType != 1 && nalType != 6 {
			fmt.Printf("--- id:%d handleNALU() Got non-I frame, unknow nalType:%d +++\n", src.Id, nalType)
		}
		// non-keyframe
		if src.syncCount > 0 {
			src.naluCount++
			//writeNALU(false, int(ts), payload)
			rtpPkg.Payload = payload
			/*
			if rtspSource.syncCount%30 == 0 {
				log.Printf("=== naluCount:%d get non-I frame flvChk len:%d ===\n", rtspSource.naluCount - 1, len(rtspSource.OutputChan))
			}
			*/
		}
	}

	chLen := len(src.OutputChan)
	if chLen < RtspBuf - 1 {
		//fmt.Printf("+++ len(rtspSource.OutputChan): %d +++\n", chLen)
		if chLen >= RtspDropTriger && rtpPkg.NalType == 1 {
			if src.dfc % 2 == 0 {
				fmt.Printf("+++ Id:%d handleNALU() drop P frame, chlen:%d \n", src.Id, chLen)
			} else {
				select {
				case <-src.rCtx.Done():
					return
				case src.OutputChan <- rtpPkg:
					return
				}
			}
			src.dfc++
		} else {
			select {
			case <-src.rCtx.Done():
				return
			case src.OutputChan <- rtpPkg:
				return
			}
		}
	} else {
		fmt.Printf("+++ id:%d handleNALU() rtsp Output chan is blocked, len(rtspSource.OutputChan): %d +++\n", src.Id, chLen)
	}
}

func (src *RtspSource)runProcessPkg(isLaunched chan error){

	defer func() {
		if p := recover(); p != nil {
			fmt.Printf("+++ rtsp source runProcessPkg() recover! p:%+v\n", p)
			debug.PrintStack()
		}

		close(src.outPkgExit)
		fmt.Printf("+++ Id:%d runProcessPkg() defer exited.\n", src.Id)
	}()

	close(isLaunched)

	fuBuffer := []byte{}
	for {
		select {
		case data := <-src.outGoingStream:
			//fmt.Printf("+++ Id:%d packet recive +++\n", src.Id)
			if data[0] == 36 && data[1] == 0 {
				cc := data[4] & 0xF
				//rtp header
				rtphdr := 12 + cc*4
				ts := (uint32(data[8]) << 24) + (uint32(data[9]) << 16) + (uint32(data[10]) << 8) + (uint32(data[11]))

				//packet number
				packno := (int64(data[6]) << 8) + int64(data[7])
				if false {
					fmt.Printf("+++ Id:%d packet num:%d +++\n", src.Id, packno)
				}

				nalType := data[4+rtphdr] & 0x1F

				if nalType >= 1 && nalType <= 23 {
					src.handleNALU(nalType, data[4+rtphdr:], ts)
					//fmt.Printf("+++ RtspReader.videow:%d , RtspReader.videoh:%d +++\n",src.RtspReader.videow, src.RtspReader.videoh)
				} else if nalType == 28 {
					isStart := data[4+rtphdr+1]&0x80 != 0
					isEnd := data[4+rtphdr+1]&0x40 != 0
					nalType := data[4+rtphdr+1] & 0x1F
					//nri := (data[4+rtphdr+1]&0x60)>>5
					nal := data[4+rtphdr]&0xE0 | data[4+rtphdr+1]&0x1F
					if isStart {
						fuBuffer = []byte{0}
					}
					fuBuffer = append(fuBuffer, data[4+rtphdr+2:]...)
					if isEnd {
						fuBuffer[0] = nal
						src.handleNALU(nalType, fuBuffer, ts)
					}
				}

			} else if data[0] == 36 && data[1] == 2 {
				// audio

				//cc := data[4] & 0xF
				//rtphdr := 12 + cc*4
				//or not payload := data[4+rtphdr:]
				//payload := data[4+rtphdr+4:]
				//outfileAAC.Write(payload)
				//log.Print("audio payload\n", hex.Dump(payload))
			}

			case <-src.rCtx.Done():
				return
		}
	}
}

func (src *RtspSource)runPlayRtspStream(signalLaunched chan error)  {

	defer func() {
		close(src.cltExit)
		fmt.Printf("+++ Id:%d runPlayRtspStream() defer exited.\n", src.Id)
	}()

	var isFirstLaunch bool = true
	src.RtspReader = rtspClientNew(src.outGoingStream)

	for {
		err := src.RtspReader.buildUpStream(src.rtspUrl, time.Second * 3)
		if err != nil {
			if isFirstLaunch {
				signalLaunched <- err
				return
			}

			fmt.Printf("--- Id:%d buildUpStream() err:%s \n", src.Id, err.Error())
			select {
			case <-src.rCtx.Done():
				return
			case <-time.After(time.Second * 1):
				continue
			}
		}


		cltExit := make(chan error, 1)
		var sLaunched chan error
		if isFirstLaunch {
			isFirstLaunch = false
			sLaunched = signalLaunched
		} else {
			sLaunched = make(chan error, 1)
		}
		go src.RtspReader.pullStream(src.rCtx, cltExit, sLaunched)

		fmt.Printf("+++ Id:%d pullStream() running.................. \n", src.Id)

		var isClosed error
		select {
		case <-src.rCtx.Done():
			<-cltExit
			return
		case isClosed = <-cltExit:
			if _, ok := isClosed.(ContextClosed); ok {
				fmt.Printf("--- Id:%d pullStream() return err:%s \n", src.Id, isClosed.Error())
				return
			} else {
				fmt.Printf("--- Id:%d pullStream() err:%s \n", src.Id, isClosed.Error())
				continue
			}
		}
	}
}

func (src *RtspSource)Play(timeoutMs time.Duration) (h264parser.SPSInfo, error) {

	var spsInfo h264parser.SPSInfo
	err := src.fsmRunner.Event(Event_Play, &spsInfo)

	return spsInfo, err
}

func (src *RtspSource)Stop() (error){

	return src.fsmRunner.Event(Event_Stop)
}

func NewRtspSource(parent context.Context, Id uint32, rtspUrl string, outChLen uint32) (*RtspSource) {

	src := new(RtspSource)
	src.rCtxParent, src.rCtxParentCclFunc = context.WithCancel(parent)

	src.Id = Id
	src.rtspUrl = rtspUrl
	src.sps = []byte{}
	src.pps = []byte{}
	src.outGoingStream = make(chan []byte, RtspCltOutChLen)
	src.OutputChan = make(chan interface {}, outChLen)

	src.fsmRunner = fsm.NewFSM(
		Status_Stop,
		fsm.Events{
			{Name: Event_Play, Src: []string{Status_Stop}, Dst: Status_Play},
			{Name: Event_Stop, Src: []string{Status_Play}, Dst: Status_Stop},
		},
		fsm.Callbacks{
			CBK_Leave_Status_Stop: func(event *fsm.Event) {
				var err error
				var isLaunched chan error

				src.rCtx, src.rCtxCclFunc = context.WithCancel(src.rCtxParent)

				isLaunched = make(chan error)
				src.outPkgExit = make(chan bool)
				src.spsCh = make(chan h264parser.SPSInfo, 1)
				go src.runProcessPkg(isLaunched)
				err = <-isLaunched


				isLaunched = make(chan error)
				src.cltExit = make(chan bool, 1)
				go src.runPlayRtspStream(isLaunched)

				err = <-isLaunched
				if _, ok := err.(NonError); !ok {
					src.rCtxCclFunc()
					event.Cancel(err)
					return
				}

				var spsInfo	h264parser.SPSInfo
				select {
				case spsInfo = <-src.spsCh:
					spsArg, ok := event.Args[0].(*h264parser.SPSInfo)
					if ok {
						*spsArg = spsInfo
					}
					break
				case <-time.After(time.Second * 5):
					src.rCtxCclFunc()
					event.Cancel(NoVideoSizeError{})
					return
				}
			},
			CBK_Leave_Status_Play: func(event *fsm.Event) {
				src.rCtxCclFunc()
				select {
				case <-src.cltExit:
					break
				case <-time.After(time.Second * 2):
					event.Cancel(TimeoutError{Err:"rtsp clt exit timeout, try again later"})
				}

				select {
				case <-src.outPkgExit:
					break
				case <-time.After(time.Second * 2):
					event.Cancel(TimeoutError{Err:"rtsp outPkg exit timeout, try again later"})
				}
			},
		},
	)

	return src
}

type RtspClient struct {

	socket   net.Conn
	outGoingStream chan []byte 	//out chanel

	host     string      	//host
	port     string      	//port
	uri      string      	//url
	auth     bool        	//aut
	login    string
	password string   	//password
	session  string   	//rtsp session
	responce string   	//responce string
	bauth    string   	//string b auth
	track    []string 	//rtsp track
	cseq     int      	//qury number
	videow   int
	videoh   int

}

func rtspClientNew(outGoingStream chan []byte) *RtspClient {

	if outGoingStream == nil {
		return nil
	}

	clt := new(RtspClient)
	clt.outGoingStream = outGoingStream
	clt.cseq = 1

	return clt
}

func (clt *RtspClient) reset4Reconnect()  {
	clt.cseq = 1

	clt.host = ""
	clt.port = ""
	clt.uri = ""
	clt.auth = false
	clt.login = ""
	clt.password = ""
	clt.session = ""
	clt.responce = ""
	clt.bauth = ""
	clt.track = []string{}
	clt.videow = 0
	clt.videoh = 0

}

func (clt *RtspClient) phaseOPTION() {

	wMsg := fmt.Sprintf("OPTIONS " + clt.uri + " RTSP/1.0\r\nCSeq: " + strconv.Itoa(clt.cseq) + "\r\n\r\n")
	clt.write(wMsg)
	rMsg := clt.read()

	if strings.Contains(rMsg, "Digest") {
		clt.authDigest("OPTIONS", rMsg)
	} else if strings.Contains(rMsg, "Basic") {
		clt.authBasic("OPTIONS", rMsg)
	} else if !strings.Contains(rMsg, "200") {
		panic(StreamInfoError{Info:fmt.Sprintf("Read OPTIONS not status code 200 OK, wMsg:%s, rMsg:%s", wMsg, rMsg)})
	}
}

func (clt *RtspClient) phaseDESCRIBE () {

	//fmt.Printf("+++ Id:%d DESCRIBE :%s  RTSP/1.0    CSeq: %s %s \r\n\r\n", clt.Id, clt.uri, strconv.Itoa(clt.cseq), clt.bauth)
	wMsg := fmt.Sprintf("DESCRIBE " + clt.uri + " RTSP/1.0\r\nCSeq: " + strconv.Itoa(clt.cseq) + clt.bauth + "\r\n\r\n")
	clt.write(wMsg)
	rMsg := clt.read()

	if strings.Contains(rMsg, "Digest") {
		clt.authDigest("DESCRIBE", rMsg)
	} else if strings.Contains(rMsg, "Basic") {
		clt.authBasic("DESCRIBE", rMsg)
	} else if !strings.Contains(rMsg, "200") {
		panic(StreamInfoError{Info:fmt.Sprintf("Read DESCRIBE not status code 200 OK, wMsg:%s, rMsg:%s", wMsg, rMsg )})
	} else {
		clt.track = clt.parseMedia(rMsg)
		if len(clt.track) == 0 {
			panic(StreamInfoError{Info:fmt.Sprintf("Track not found. wMsg:%s, rMsg:%s ", wMsg, rMsg)})
		}
	}
}

func (clt *RtspClient) phaseSETUP() {

	wMsg := fmt.Sprintf("SETUP " + clt.uri + "/" + clt.track[0] + " RTSP/1.0\r\nCSeq: " + strconv.Itoa(clt.cseq) + "\r\nTransport: RTP/AVP/TCP;unicast;interleaved=0-1" + clt.bauth + "\r\n\r\n")
	clt.write(wMsg)
	rMsg := clt.read()

	if !strings.Contains(rMsg, "200") {
		if strings.Contains(rMsg, "401") {
			str := clt.authDigest_Only("SETUP", rMsg)
			ado_wMsg := fmt.Sprintf("SETUP " + clt.uri + "/" + clt.track[0] + " RTSP/1.0\r\nCSeq: " + strconv.Itoa(clt.cseq) + "\r\nTransport: RTP/AVP/TCP;unicast;interleaved=0-1" + clt.bauth + str + "\r\n\r\n")
			clt.write(ado_wMsg)
			ado_rMsg := clt.read()

			if !strings.Contains(ado_rMsg, "200") {
				panic(StreamInfoError{Info:fmt.Sprintf("Read SETUP authDigest_Only() not 200. wMsg:%s, rMsg:%s ado_wMsg:%s, ado_rMsg:%s", wMsg, rMsg, ado_wMsg, ado_rMsg)})
			} else {
				clt.session = clt.parseSession(ado_rMsg)
			}
		} else {
			panic(StreamInfoError{Info:fmt.Sprintf("Read SETUP not 200 and 401, wMsg:%s rMsg:%s", wMsg, rMsg)})
		}
	} else {
		clt.session = clt.parseSession(rMsg)
	}

	if len(clt.track) > 1 {

		wMsg = fmt.Sprintf("SETUP " + clt.uri + "/" + clt.track[1] + " RTSP/1.0\r\nCSeq: " + strconv.Itoa(clt.cseq) + "\r\nTransport: RTP/AVP/TCP;unicast;interleaved=2-3" + "\r\nSession: " + clt.session + clt.bauth + "\r\n\r\n")
		clt.write(wMsg)
		rMsg = clt.read()

		if !strings.Contains(rMsg, "200") {
			if strings.Contains(rMsg, "401") {
				str := clt.authDigest_Only("SETUP", rMsg)
				ado_wMsg := fmt.Sprintf("SETUP " + clt.uri + "/" + clt.track[1] + " RTSP/1.0\r\nCSeq: " + strconv.Itoa(clt.cseq) + "\r\nTransport: RTP/AVP/TCP;unicast;interleaved=2-3" + clt.bauth + str + "\r\n\r\n")
				clt.write(ado_wMsg)
				ado_rMsg := clt.read()

				if !strings.Contains(ado_rMsg, "200") {
					panic(StreamInfoError{Info:fmt.Sprintf("Read SETUP authDigest_Only() Audio not 200, Write SETUP, wMsg:%s, rMsg:%s, ado_wMsg:%s, ado_rMsg:%s", wMsg, rMsg, ado_wMsg, ado_rMsg)})
				} else {
					clt.session = clt.parseSession(ado_rMsg)
				}
			} else {
				panic(StreamInfoError{Info:fmt.Sprintf("Read SETUP Audio not 200, wMsg:%s, rMsg:%s", wMsg, rMsg)})
			}
		} else {
			clt.session = clt.parseSession(rMsg)
		}
	}
}

func (clt *RtspClient) phasePLAY() {

	wMsg := fmt.Sprintf("PLAY " + clt.uri + " RTSP/1.0\r\nCSeq: " + strconv.Itoa(clt.cseq) + "\r\nSession: " + clt.session + clt.bauth + "\r\n\r\n")
	clt.write(wMsg)
	rMsg := clt.read()

	if !strings.Contains(rMsg, "200") {
		if strings.Contains(rMsg, "401") {
			str := clt.authDigest_Only("PLAY", rMsg)

			ado_wMsg := fmt.Sprintf("PLAY " + clt.uri + " RTSP/1.0\r\nCSeq: " + strconv.Itoa(clt.cseq) + "\r\nSession: " + clt.session + clt.bauth + str + "\r\n\r\n")
			clt.write(ado_wMsg)
			ado_rMsg := clt.read()

			if !strings.Contains(ado_rMsg, "200") {
				panic(StreamInfoError{Info:fmt.Sprintf("Read PLAY not 200, wMsg:%s, rMsg:%s, ado_wMsg:%s, ado_rMsg:%s", wMsg, rMsg, ado_wMsg, ado_rMsg )})
			}
			/*
			else {
				//clt.session = clt.parseSession(message)
				//log.Print(message)
				//fmt.Printf("+++ Id:%d rtsp client Start() go rtspRtpLoop(), msg:%s\n", clt.Id, message)
				//fmt.Printf("+++ Id:%d rtsp client Start() go rtspRtpLoop()\n", clt.Id)
				isLaunched := make(chan bool)
				go clt.rtspRtpLoop(parent, isExit, isLaunched)
				<-isLaunched
				//fmt.Printf("+++ Id:%d rtsp client Start() Read PLAY not 200 but 401, Write PLAY, Read PLAY ok, go rtspRtpLoop()\n", clt.Id)
				err = nil
				return isExit, nil
			}
			*/
		} else {
			panic(StreamInfoError{Info:fmt.Sprintf("Read PLAY not 200 and 401, wMsg:%s, rMsg:%s", wMsg, rMsg)})
		}
	}
	/*else {
		//log.Print(message)
		//fmt.Printf("+++ Id:%d rtsp client Start() go rtspRtpLoop(), msg:%s\n", clt.Id, message)
		//fmt.Printf("+++ Id:%d rtsp client Start() go rtspRtpLoop()\n", clt.Id)
		isLaunched := make(chan bool)
		go clt.rtspRtpLoop(parent, isExit, isLaunched)
		<-isLaunched
		//fmt.Printf("+++ Id:%d Read PLAY 200 ok, go rtspRtpLoop()\n", clt.Id)
		err = nil
		return isExit, nil
	}
	*/

}

func (clt *RtspClient) buildUpStream(rtsp_url string, connectTimeOut time.Duration) (err error) {

	defer func() {
		var errMsg string
		if p := recover(); p != nil {
			errMsg = fmt.Sprintf("Rtsp clt buildUpStream() recover! p:%+v", p)
			err = errors.New(fmt.Sprintf("panic:%s", errMsg))

			debug.PrintStack()
		}

		if err != nil {
			if clt.socket != nil {
				clt.socket.Close()
			}
		}
	}()

	// parse rtsp url
	clt.parseUrl(rtsp_url)

	// connect
	if err = clt.connect(connectTimeOut); err != nil{
		return err
	}

	//PHASE 1 OPTION
	clt.phaseOPTION()

	//PHASE 2 DESCRIBE
	clt.phaseDESCRIBE()

	//PHASE 3 SETUP
	clt.phaseSETUP()

	//PHASE 4 PLAY
	clt.phasePLAY()

	return nil
}

func (clt *RtspClient) pullStream(ctx context.Context, signalExit chan<- error, signalLaunched chan<- error) {

	var err error
	defer func() {
		if p := recover(); p != nil {
			pInfo := fmt.Sprintf("pullStream() recover! p:%+v", p)

			v, isOk := p.(SocketReadError)
			if isOk {
				err = v
			} else {
				debug.PrintStack()
				err = errors.New(pInfo)
			}
		}

		clt.socket.Close()
		signalExit<-err
		fmt.Printf("+++ pullStream() defer exited.\n")
	}()

	header := make([]byte, 4)
	payload := make([]byte, RtspPayloadMax)
	sync_b := make([]byte, 1)
	timer := time.Now()

	signalLaunched <- NonError{}

	for {
		//var rtspPanicInfo RtspPanicInfo
		if int(time.Now().Sub(timer).Seconds()) > 50 {
			clt.write("OPTIONS " + clt.uri + " RTSP/1.0\r\nCSeq: " + strconv.Itoa(clt.cseq) + "\r\nSession: " + clt.session + clt.bauth + "\r\n\r\n")
			timer = time.Now()
		}

		clt.socket.SetDeadline(time.Now().Add(2 * time.Second))

		//read rtp hdr 4
		var rCnt int
		rCnt, err = io.ReadFull(clt.socket, header)
		if err != nil {
			//rtspPanicInfo.errType = Err_SocketRead
			//rtspPanicInfo.errRaw = errors.New(fmt.Sprintf("rtspRtpLoop() read rtp hdr err:%s", err.Error()))
			//panic(rtspPanicInfo)
			panic(SocketReadError{Info:"pullStream() read rtp hdr", Err:err})
		} else if rCnt != 4 {
			panic(SocketReadError{Info:fmt.Sprintf("pullStream() read rtp hdr len:%d != 4", rCnt)})
		}

		if header[0] != 36 {
			//log.Println("desync?", clt.host)
			for {
				///////////////////////////skeep/////////////////////////////////////
				if rCnt, err = io.ReadFull(clt.socket, sync_b); err != nil && rCnt != 1 {
					//rtspPanicInfo.errType = Err_SocketRead
					//rtspPanicInfo.errRaw = errors.New(fmt.Sprintf("rtspRtpLoop() read sync_b err:%s n:%d ", err.Error(), n))
					//panic(rtspPanicInfo)
					panic(SocketReadError{Info:fmt.Sprintf("pullStream() read sync_b n:%d ", rCnt), Err:err})
				} else if sync_b[0] == 36 {
					header[0] = 36
					if rCnt, err = io.ReadFull(clt.socket, header[1:]); err != nil && rCnt == 3 {
						//rtspPanicInfo.errType = Err_SocketRead
						//rtspPanicInfo.errRaw = errors.New(fmt.Sprintf("rtspRtpLoop() sync_b[0] == 36 read err:%s n:%d ", err.Error(), n))
						//panic(rtspPanicInfo)
						panic(SocketReadError{Info:fmt.Sprintf("pullStream() sync_b[0] == 36 read n:%d ", rCnt), Err:err})
					}
					break
				}
			}
		}

		payloadLen := (int)(header[2])<<8 + (int)(header[3])
		if payloadLen > RtspPayloadMax || payloadLen < 12 {
			//panic(fmt.Sprintf("rtspRtpLoop() desync, uri:%s payloadLen:%d ", clt.uri, payloadLen))
			panic(StreamInfoError{Info:fmt.Sprintf("pullStream() desync, uri:%s payloadLen:%d ", clt.uri, payloadLen)})
		}
		if rCnt, err = io.ReadFull(clt.socket, payload[:payloadLen]); err != nil || rCnt != payloadLen {
			if err != nil {
				//rtspPanicInfo.errType = Err_SocketRead
				//rtspPanicInfo.errRaw = errors.New(fmt.Sprintf("rtspRtpLoop() read payload err:%s, n:%d ", err.Error(), payloadLen))
				//panic(rtspPanicInfo)
				panic(SocketReadError{Info:fmt.Sprintf("pullStream() read payload len:%d", payloadLen), Err:err})
			} else {
				//panic(fmt.Sprintf("pullStream() read payload err:%s, n:%d ", err.Error(), payloadLen))
				panic(StreamInfoError{Info:fmt.Sprintf("pullStream() read payload len:%d ", payloadLen)})
			}
		} else {
			select {
			case <-ctx.Done():
				err = ContextClosed{}
				return
			case clt.outGoingStream <- append(header, payload[:rCnt]...):
				continue
			}
		}
	}
}

func (clt *RtspClient) connect(timeOut time.Duration) (err error) {

	d := &net.Dialer{Timeout: timeOut}
	conn, err := d.Dial("tcp", clt.host+":"+clt.port)

	if err != nil {
		return ConnectError{Err:err}
	}
	clt.socket = conn
	return nil
}

func (clt *RtspClient) write(message string) {
	clt.cseq += 1
	if _, e := clt.socket.Write([]byte(message)); e != nil {
		panic(SocketWriteError{Err:e})
	}
}

func (clt *RtspClient) read() string {
	buffer := make([]byte, 4096)
	if nb, err := clt.socket.Read(buffer); err != nil || nb <= 0 {
		if err == nil {
			panic(SocketReadError{Info:"len <= 0"})
		} else {
			panic(SocketReadError{Err:err})
		}
	} else {
		return string(buffer[:nb])
	}
}

func (clt *RtspClient) authBasic(phase string, message string) {

	clt.bauth = "\r\nAuthorization: Basic " + b64.StdEncoding.EncodeToString([]byte(clt.login+":"+clt.password))
	wMsg := phase + " " + clt.uri + " RTSP/1.0\r\nCSeq: " + strconv.Itoa(clt.cseq) + clt.bauth + "\r\n\r\n"
	clt.write(phase + " " + clt.uri + " RTSP/1.0\r\nCSeq: " + strconv.Itoa(clt.cseq) + clt.bauth + "\r\n\r\n")
	rMsg := clt.read()

	if strings.Contains(rMsg, "200") {
		clt.track = clt.parseMedia(rMsg)
		if len(clt.track) == 0 {
			panic(StreamInfoError{Info:fmt.Sprintf("authBasic() track not found. wMsg:%s, rMsg:%s ", wMsg, rMsg)})
		}
	} else {
		panic(StreamInfoError{Info:fmt.Sprintf("authBasic() response not contain 200. wMsg:%s rMsg:%s", wMsg, rMsg)})
	}
}

func (clt *RtspClient) authDigest(phase string, message string) {
	var rMsg string

	nonce := clt.parseDirective(message, "nonce")
	realm := clt.parseDirective(message, "realm")
	hs1 := clt.getMD5Hash(clt.login + ":" + realm + ":" + clt.password)
	hs2 := clt.getMD5Hash(phase + ":" + clt.uri)
	response := clt.getMD5Hash(hs1 + ":" + nonce + ":" + hs2)
	dauth := "\r\n" + `Authorization: Digest username="` + clt.login + `", realm="` + realm + `", nonce="` + nonce + `", uri="` + clt.uri + `", response="` + response + `"`

	wMsg := phase + " " + clt.uri + " RTSP/1.0\r\nCSeq: " + strconv.Itoa(clt.cseq) + dauth + "\r\n\r\n"
	clt.write(wMsg)
	rMsg = clt.read()

	if strings.Contains(rMsg, "200") {
		clt.track = clt.parseMedia(rMsg)
		if len(clt.track) == 0 {
			panic(StreamInfoError{fmt.Sprintf("authDigest() track not found. wMsg:%s, rMsg:%s ", wMsg, rMsg)})
		}
	} else {
		panic(StreamInfoError{fmt.Sprintf("authDigest() read response doesn't contain 200. wMsg:%s rMsg:%s", wMsg, rMsg)})
	}
}

func (clt *RtspClient) authDigest_Only(phase string, message string) string {
	nonce := clt.parseDirective(message, "nonce")
	realm := clt.parseDirective(message, "realm")
	hs1 := clt.getMD5Hash(clt.login + ":" + realm + ":" + clt.password)
	hs2 := clt.getMD5Hash(phase + ":" + clt.uri)
	responce := clt.getMD5Hash(hs1 + ":" + nonce + ":" + hs2)
	dauth := "\r\n" + `Authorization: Digest username="` + clt.login + `", realm="` + realm + `", nonce="` + nonce + `", uri="` + clt.uri + `", response="` + responce + `"`
	return dauth
}

func (clt *RtspClient) parseUrl(rtsp_url string) {

	u, err := url.Parse(rtsp_url)
	if err != nil {
		panic(fmt.Sprintf("rtsp url incorrect, url:%s", rtsp_url))
	}

	phost := strings.Split(u.Host, ":")
	clt.host = phost[0]
	if len(phost) == 2 {
		clt.port = phost[1]
	} else {
		clt.port = "554"
	}
	if u.User != nil {
		clt.login = u.User.Username()
		clt.password, clt.auth = u.User.Password()
	}
	if u.RawQuery != "" {
		clt.uri = "rtsp://" + clt.host + ":" + clt.port + u.Path + "?" + string(u.RawQuery)
	} else {
		clt.uri = "rtsp://" + clt.host + ":" + clt.port + u.Path
	}
}

func (clt *RtspClient) parseDirective(header, name string) string {
	index := strings.Index(header, name)
	if index == -1 {
		return ""
	}
	start := 1 + index + strings.Index(header[index:], `"`)
	end := start + strings.Index(header[start:], `"`)
	return strings.TrimSpace(header[start:end])
}

func (clt *RtspClient) parseSession(header string) string {
	mparsed := strings.Split(header, "\r\n")
	for _, element := range mparsed {
		if strings.Contains(element, "Session:") {
			if strings.Contains(element, ";") {
				fist := strings.Split(element, ";")[0]
				return fist[9:]
			} else {
				return element[9:]
			}
		}
	}
	return ""
}

func (clt *RtspClient) getMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func (clt *RtspClient) parseMedia(header string) []string {
	letters := []string{}
	mparsed := strings.Split(header, "\r\n")
	paste := ""
	for _, element := range mparsed {
		if strings.Contains(element, "a=control:") && !strings.Contains(element, "*") && strings.Contains(element, "tra") {
			paste = element[10:]
			if strings.Contains(element, "/") {
				striped := strings.Split(element, "/")
				paste = striped[len(striped)-1]
			}
			letters = append(letters, paste)
		}

		dimensionsPrefix := "a=x-dimensions:"
		if strings.HasPrefix(element, dimensionsPrefix) {
			dims := []int{}
			for _, s := range strings.Split(element[len(dimensionsPrefix):], ",") {
				v := 0
				fmt.Sscanf(s, "%d", &v)
				if v <= 0 {
					break
				}
				dims = append(dims, v)
			}
			if len(dims) == 2 {
				clt.videow = dims[0]
				clt.videoh = dims[1]
			}
		}
	}
	return letters
}


