
package rtsp

import "fmt"

type ConnectError struct {
	Err error
}

func (e ConnectError) Error() string {
	return "rtsp connect err: " + e.Err.Error()
}

type SocketReadError struct {
	Info string
	Err error
}

func (e SocketReadError) Error() string {
	return fmt.Sprintf("rtsp socket read info:%s, err:%+v", e.Info, e.Err)
}

type SocketWriteError struct {
	Info string
	Err error
}

func (e SocketWriteError) Error() string {
	return fmt.Sprintf("rtsp socket write info:%s, err:%+v", e.Info, e.Err)
}

type StreamInfoError struct {
	Info string
}

func (e StreamInfoError) Error() string {
	return fmt.Sprintf("rtsp stream info err:%s", e.Info)
}

type NoVideoSizeError struct {

}

func (e NoVideoSizeError) Error() string  {
	return "rtsp stream have no video size"
}

type NonError struct {
}

func (e NonError) Error() string {
	return "no error."
}

type ContextClosed struct {
}

func (e ContextClosed) Error() string {
	return "context closed"
}

type TimeoutError struct {
	Err string
}

func (e TimeoutError) Error() string {
	return "timeout info:" + e.Err
}



