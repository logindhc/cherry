package cherryHandler

import (
	cherryCode "github.com/cherry-game/cherry/code"
	cherryError "github.com/cherry-game/cherry/error"
	facade "github.com/cherry-game/cherry/facade"
	cherryLogger "github.com/cherry-game/cherry/logger"
	cherryProto "github.com/cherry-game/cherry/net/proto"
	"github.com/nats-io/nats.go"
	"reflect"
	"runtime/debug"
)

type (
	ExecutorRemote struct {
		facade.IApplication
		HandlerFn    *facade.HandlerFn
		RemotePacket *cherryProto.RemotePacket
		NatsMsg      *nats.Msg
	}
)

func (p *ExecutorRemote) Invoke() {
	defer func() {
		if rev := recover(); rev != nil {
			cherryLogger.Warnf("recover in Remote. %s", string(debug.Stack()))
			cherryLogger.Warnf("msg = [%+v]", p.RemotePacket)
		}
	}()

	argsLen := len(p.HandlerFn.InArgs)
	if argsLen < 0 || argsLen > 1 {
		cherryLogger.Warnf("[Route = %v] method in args error.", p.RemotePacket.Route)
		cherryLogger.Warnf("func() or func(request)")
		return
	}

	var ret []reflect.Value
	var params []reflect.Value

	switch argsLen {
	case 0:
		ret = p.HandlerFn.Value.Call(params)
		break
	case 1:
		val, err := p.unmarshalData()
		if err != nil {
			cherryLogger.Warnf("[Route = %s] unmarshal data error.error = %s", p.RemotePacket.Route, err)
			return
		}
		params = make([]reflect.Value, 1)
		params[0] = reflect.ValueOf(val)

		ret = p.HandlerFn.Value.Call(params)
		break
	}

	if p.NatsMsg.Reply == "" {
		return
	}

	rsp := &cherryProto.Response{
		Code: cherryCode.OK,
	}

	if len(ret) == 1 {
		if val := ret[0].Interface(); val != nil {
			if code, ok := val.(int32); ok {
				rsp.Code = code
			}
		}

		rspData, _ := p.Marshal(rsp)
		err := p.NatsMsg.Respond(rspData)
		if err != nil {
			cherryLogger.Warn(err)
		}

	} else if len(ret) == 2 {
		if val := ret[1].Interface(); val != nil {
			if code, ok := val.(int32); ok {
				rsp.Code = code
			}
		}

		if ret[0].IsNil() == false {

			data, err := p.Marshal(ret[0].Interface())
			if err != nil {
				rsp.Code = cherryCode.RPCRemoteExecuteError
				cherryLogger.Warn(err)
			} else {
				rsp.Data = data
			}
		}

		rspData, _ := p.Marshal(rsp)
		err := p.NatsMsg.Respond(rspData)
		if err != nil {
			cherryLogger.Warn(err)
		}
	}
}

func (p *ExecutorRemote) unmarshalData() (interface{}, error) {
	if len(p.HandlerFn.InArgs) != 1 {
		return nil, cherryError.Error("remote handler params len is error.")
	}

	in2 := p.HandlerFn.InArgs[0]

	var val interface{}
	val = reflect.New(in2.Elem()).Interface()
	err := p.Unmarshal(p.RemotePacket.Data, val)
	if err != nil {
		return nil, err
	}

	return val, err
}

func (p *ExecutorRemote) String() string {
	return p.RemotePacket.Route
}
