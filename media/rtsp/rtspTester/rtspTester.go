package main

import (
	"flag"
	"fmt"

	"context"
	"time"
	"runtime/debug"

	"github.com/falconray0704/u4go/media/rtsp"
)

func runRtspSourceTestOnce(url string, isExitOutCh chan <- bool)  {

	var tc, fc int
	defer func() {
		if p := recover(); p != nil {

			fmt.Printf("+++ runRtspSourceTestOnce() recover! p:%+v", p)

			debug.PrintStack()

		}
		close(isExitOutCh)
		fmt.Printf("+++ runRtspSourceTestOnce defer exited , tc:%d fc:%d ++++++++++++++++++\n", tc, fc)
	}()

	src := rtsp.NewRtspSource(context.Background(), "0", url, 30)
	defer func() {
		if err := src.Stop(); err != nil {
			fmt.Printf("+++ runRtspSourceTestOnce defer src.Stop() err:%+v ++++++++++++++++++\n", err)
		} else {
		}
	}()

	okExit := make(chan bool, 10)
	go func() {
		defer func() {
			//fmt.Printf("+++ Get frame defer exited-------------------------------\n")
		}()
		var ifc = 0
		for {
			select {
			case val := <- src.OutputChan:
				if val != nil {
					pkg := val.(*rtsp.RtspPacket)
					if pkg.NalType != 1 {
						//fmt.Printf("+++ Get one frame, type:%d \n", pkg.NalType)
					}
					if pkg.NalType == 5 {
						ifc++
					}
					if ifc == 2 {
						okExit <- true
					}
				} else {
					fmt.Printf("--- Get frame returned\n")
					return
				}
			}
		}
	}()

	if spsinfo, err := src.Play(time.Second * 10); err != nil {
		fmt.Printf("--- rtsp source fail, err:%s \n", err.Error())
		fc++
	} else {
		if spsinfo.Width == 0 || spsinfo.Height == 0 {
			fmt.Printf("--- rtsp sps info incorrect, spsInfo:%+v \n", spsinfo)
		}

		<-okExit
	}

	tc++
	fmt.Printf("========================== TC: %d fc:%d \n", tc, fc)
}

func runRtspSourceTestNetworkReconnect(url string, isExitOutCh chan <- bool)  {

	var tc, fc int
	defer func() {
		if p := recover(); p != nil {
			fmt.Printf("+++ runRtspSourceTestNetworkReconnect() recover! p:%+v", p)
			debug.PrintStack()
		}
		close(isExitOutCh)
		fmt.Printf("+++ runRtspSourceTestNetworkReconnect() defer exited , tc:%d fc:%d ++++++++++++++++++\n", tc, fc)
	}()

	src := rtsp.NewRtspSource(context.Background(), "0", url, 30)

	okExit := make(chan bool, 10)
	isFinished := make(chan bool)
	go func() {
		defer func() {
			fmt.Printf("+++ runRtspSourceTestNetworkReconnect() Get frame defer exited-------------------------------\n")
		}()
		var ifc = 0
		for {
			select {
			case val := <- src.OutputChan:
				if val != nil {
					pkg := val.(*rtsp.RtspPacket)
					if pkg.NalType != 1 {
						//fmt.Printf("+++ Get one frame, type:%d \n", pkg.NalType)
					}
					if pkg.NalType == 5 {
						ifc++
						fmt.Printf("+++ runRtspSourceTestNetworkReconnect() Get one I frame, ifc:%d type:%d ts:%d \n", ifc, pkg.NalType, pkg.Timestamp)
					}

					if ifc >= 30 {
						okExit<-true
						return
					}
				} else {
					fmt.Printf("--- runRtspSourceTestNetworkReconnect() Get frame returned\n")
					return
				}
			case <-time.After(time.Second * 2):
				fmt.Printf("--- runRtspSourceTestNetworkReconnect() Get frame timeout, stream disconnected.\n")
				break
			case <-isFinished:
				return
			}
		}
	}()

	if spsinfo, err := src.Play(time.Second * 10); err != nil {
		fmt.Printf("--- rtsp source fail, err:%s \n", err.Error())
		fc++
	} else {
		if spsinfo.Width == 0 || spsinfo.Height == 0 {
			fmt.Printf("--- rtsp sps info incorrect, spsInfo:%+v \n", spsinfo)
		}

		select {
		case <-okExit:
			break
		case <-time.After(time.Second * 20):
			break
		}

		if err := src.Stop(); err != nil {
			fmt.Printf("+++ runRtspSourceTestNetworkReconnect() defer src.Stop() err:%+v ++++++++++++++++++\n", err)
		} else {
		}
		close(isFinished)
	}

	tc++
	fmt.Printf("========================== TC: %d fc:%d \n", tc, fc)
}

func runRtspSourceTestPlayStopLoop(url string, isExitOutCh chan <- bool)  {

	var tc, fc int
	defer func() {
		if p := recover(); p != nil {

			fmt.Printf("+++ runRtspSourceTestOnce() recover! p:%+v", p)

			debug.PrintStack()

		}
		close(isExitOutCh)
		fmt.Printf("+++ runRtspSourceTestOnce defer exited , tc:%d fc:%d ++++++++++++++++++\n", tc, fc)
		fmt.Printf("========================== TC: %d fc:%d \n", tc, fc)
	}()

	src := rtsp.NewRtspSource(context.Background(), "0", url, 30)
	defer func() {
		if err := src.Stop(); err != nil {
			fmt.Printf("+++ runRtspSourceTestOnce defer src.Stop() err:%+v ++++++++++++++++++\n", err)
		} else {
		}
	}()

	okExit := make(chan bool)
	isFinished := make(chan bool)
	go func() {
		defer func() {
			fmt.Printf("+++ runRtspSourceTestPlayStopLoop() Get frame defer exited-------------------------------\n")
		}()
		var ifc = 0
		for {
			select {
			case val := <- src.OutputChan:
				if val != nil {
					pkg := val.(*rtsp.RtspPacket)
					if pkg.NalType != 1 {
						//fmt.Printf("+++ Get one frame, type:%d \n", pkg.NalType)
					}
					if pkg.NalType == 5 {
						ifc++
					}
					if ifc == 2 {
						okExit <- true
						ifc = 0
					}
				} else {
					fmt.Printf("--- Get frame returned\n")
					return
				}
			case <-isFinished:
				fmt.Printf("--- Get frame returned\n")
				return
			}
		}
	}()

	for i := 0; i < 3600 ; i++ {
		tc++
		if spsinfo, err := src.Play(time.Second * 10); err != nil {
			fmt.Printf("--- rtsp source fail, err:%s \n", err.Error())
			fc++
		} else {
			if spsinfo.Width == 0 || spsinfo.Height == 0 {
				fmt.Printf("--- rtsp sps info incorrect, spsInfo:%+v \n", spsinfo)
			}

			select {
			case <-okExit:
				if err := src.Stop(); err != nil {
					fmt.Printf("+++ runRtspSourceTestPlayStopLoop() src.Stop() err:%+v ++++++++++++++++++\n", err)
				} else {
					fmt.Printf("+++ runRtspSourceTestPlayStopLoop() one loop success, tc:%d fc:%d ++++++++++++++++++\n", tc, fc)
				}
				break
			case <-time.After(time.Second * 15):
				if err := src.Stop(); err != nil {
					fmt.Printf("+++ runRtspSourceTestPlayStopLoop() get I frame timeout src.Stop() err:%+v ++++++++++++++++++\n", err)
				} else {
				}
				fc++
				fmt.Printf("+++ runRtspSourceTestPlayStopLoop() get I frame timeout, tc:%d fc:%d ++++++++++++++++++\n", tc, fc)
				break
			}
		}
	}

}



func main()  {

	url := flag.String("rtspUrl","rtsp://yt:grg123456@10.1.174.158:554/h264/ch1/main/av_stream", "rtsp source for connection.")
	flag.Parse()

	isExit := make(chan bool)
	//go runRtspSourceTestOnce(*url, isExit)
	go runRtspSourceTestPlayStopLoop(*url, isExit)

	// manual test
	//go runRtspSourceTestNetworkReconnect(*url, isExit)

	<-isExit
	fmt.Printf("+++ main() test finished\n")
	for {
		time.Sleep(time.Second * 1)
	}

	fmt.Printf("+++ main() exit\n")
}






