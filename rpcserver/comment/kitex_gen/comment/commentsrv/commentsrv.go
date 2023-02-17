// Code generated by Kitex v0.4.4. DO NOT EDIT.

package commentsrv

import (
	"context"
	comment "dousheng/rpcserver/comment/kitex_gen/comment"
	"fmt"
	client "github.com/cloudwego/kitex/client"
	kitex "github.com/cloudwego/kitex/pkg/serviceinfo"
	streaming "github.com/cloudwego/kitex/pkg/streaming"
	proto "google.golang.org/protobuf/proto"
)

func serviceInfo() *kitex.ServiceInfo {
	return commentSrvServiceInfo
}

var commentSrvServiceInfo = NewServiceInfo()

func NewServiceInfo() *kitex.ServiceInfo {
	serviceName := "CommentSrv"
	handlerType := (*comment.CommentSrv)(nil)
	methods := map[string]kitex.MethodInfo{
		"CommentAction": kitex.NewMethodInfo(commentActionHandler, newCommentActionArgs, newCommentActionResult, false),
		"CommentList":   kitex.NewMethodInfo(commentListHandler, newCommentListArgs, newCommentListResult, false),
	}
	extra := map[string]interface{}{
		"PackageName": "comment",
	}
	svcInfo := &kitex.ServiceInfo{
		ServiceName:     serviceName,
		HandlerType:     handlerType,
		Methods:         methods,
		PayloadCodec:    kitex.Protobuf,
		KiteXGenVersion: "v0.4.4",
		Extra:           extra,
	}
	return svcInfo
}

func commentActionHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	switch s := arg.(type) {
	case *streaming.Args:
		st := s.Stream
		req := new(comment.DouyinCommentActionRequest)
		if err := st.RecvMsg(req); err != nil {
			return err
		}
		resp, err := handler.(comment.CommentSrv).CommentAction(ctx, req)
		if err != nil {
			return err
		}
		if err := st.SendMsg(resp); err != nil {
			return err
		}
	case *CommentActionArgs:
		success, err := handler.(comment.CommentSrv).CommentAction(ctx, s.Req)
		if err != nil {
			return err
		}
		realResult := result.(*CommentActionResult)
		realResult.Success = success
	}
	return nil
}
func newCommentActionArgs() interface{} {
	return &CommentActionArgs{}
}

func newCommentActionResult() interface{} {
	return &CommentActionResult{}
}

type CommentActionArgs struct {
	Req *comment.DouyinCommentActionRequest
}

func (p *CommentActionArgs) FastRead(buf []byte, _type int8, number int32) (n int, err error) {
	if !p.IsSetReq() {
		p.Req = new(comment.DouyinCommentActionRequest)
	}
	return p.Req.FastRead(buf, _type, number)
}

func (p *CommentActionArgs) FastWrite(buf []byte) (n int) {
	if !p.IsSetReq() {
		return 0
	}
	return p.Req.FastWrite(buf)
}

func (p *CommentActionArgs) Size() (n int) {
	if !p.IsSetReq() {
		return 0
	}
	return p.Req.Size()
}

func (p *CommentActionArgs) Marshal(out []byte) ([]byte, error) {
	if !p.IsSetReq() {
		return out, fmt.Errorf("No req in CommentActionArgs")
	}
	return proto.Marshal(p.Req)
}

func (p *CommentActionArgs) Unmarshal(in []byte) error {
	msg := new(comment.DouyinCommentActionRequest)
	if err := proto.Unmarshal(in, msg); err != nil {
		return err
	}
	p.Req = msg
	return nil
}

var CommentActionArgs_Req_DEFAULT *comment.DouyinCommentActionRequest

func (p *CommentActionArgs) GetReq() *comment.DouyinCommentActionRequest {
	if !p.IsSetReq() {
		return CommentActionArgs_Req_DEFAULT
	}
	return p.Req
}

func (p *CommentActionArgs) IsSetReq() bool {
	return p.Req != nil
}

type CommentActionResult struct {
	Success *comment.DouyinCommentActionResponse
}

var CommentActionResult_Success_DEFAULT *comment.DouyinCommentActionResponse

func (p *CommentActionResult) FastRead(buf []byte, _type int8, number int32) (n int, err error) {
	if !p.IsSetSuccess() {
		p.Success = new(comment.DouyinCommentActionResponse)
	}
	return p.Success.FastRead(buf, _type, number)
}

func (p *CommentActionResult) FastWrite(buf []byte) (n int) {
	if !p.IsSetSuccess() {
		return 0
	}
	return p.Success.FastWrite(buf)
}

func (p *CommentActionResult) Size() (n int) {
	if !p.IsSetSuccess() {
		return 0
	}
	return p.Success.Size()
}

func (p *CommentActionResult) Marshal(out []byte) ([]byte, error) {
	if !p.IsSetSuccess() {
		return out, fmt.Errorf("No req in CommentActionResult")
	}
	return proto.Marshal(p.Success)
}

func (p *CommentActionResult) Unmarshal(in []byte) error {
	msg := new(comment.DouyinCommentActionResponse)
	if err := proto.Unmarshal(in, msg); err != nil {
		return err
	}
	p.Success = msg
	return nil
}

func (p *CommentActionResult) GetSuccess() *comment.DouyinCommentActionResponse {
	if !p.IsSetSuccess() {
		return CommentActionResult_Success_DEFAULT
	}
	return p.Success
}

func (p *CommentActionResult) SetSuccess(x interface{}) {
	p.Success = x.(*comment.DouyinCommentActionResponse)
}

func (p *CommentActionResult) IsSetSuccess() bool {
	return p.Success != nil
}

func commentListHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	switch s := arg.(type) {
	case *streaming.Args:
		st := s.Stream
		req := new(comment.DouyinCommentListRequest)
		if err := st.RecvMsg(req); err != nil {
			return err
		}
		resp, err := handler.(comment.CommentSrv).CommentList(ctx, req)
		if err != nil {
			return err
		}
		if err := st.SendMsg(resp); err != nil {
			return err
		}
	case *CommentListArgs:
		success, err := handler.(comment.CommentSrv).CommentList(ctx, s.Req)
		if err != nil {
			return err
		}
		realResult := result.(*CommentListResult)
		realResult.Success = success
	}
	return nil
}
func newCommentListArgs() interface{} {
	return &CommentListArgs{}
}

func newCommentListResult() interface{} {
	return &CommentListResult{}
}

type CommentListArgs struct {
	Req *comment.DouyinCommentListRequest
}

func (p *CommentListArgs) FastRead(buf []byte, _type int8, number int32) (n int, err error) {
	if !p.IsSetReq() {
		p.Req = new(comment.DouyinCommentListRequest)
	}
	return p.Req.FastRead(buf, _type, number)
}

func (p *CommentListArgs) FastWrite(buf []byte) (n int) {
	if !p.IsSetReq() {
		return 0
	}
	return p.Req.FastWrite(buf)
}

func (p *CommentListArgs) Size() (n int) {
	if !p.IsSetReq() {
		return 0
	}
	return p.Req.Size()
}

func (p *CommentListArgs) Marshal(out []byte) ([]byte, error) {
	if !p.IsSetReq() {
		return out, fmt.Errorf("No req in CommentListArgs")
	}
	return proto.Marshal(p.Req)
}

func (p *CommentListArgs) Unmarshal(in []byte) error {
	msg := new(comment.DouyinCommentListRequest)
	if err := proto.Unmarshal(in, msg); err != nil {
		return err
	}
	p.Req = msg
	return nil
}

var CommentListArgs_Req_DEFAULT *comment.DouyinCommentListRequest

func (p *CommentListArgs) GetReq() *comment.DouyinCommentListRequest {
	if !p.IsSetReq() {
		return CommentListArgs_Req_DEFAULT
	}
	return p.Req
}

func (p *CommentListArgs) IsSetReq() bool {
	return p.Req != nil
}

type CommentListResult struct {
	Success *comment.DouyinCommentListResponse
}

var CommentListResult_Success_DEFAULT *comment.DouyinCommentListResponse

func (p *CommentListResult) FastRead(buf []byte, _type int8, number int32) (n int, err error) {
	if !p.IsSetSuccess() {
		p.Success = new(comment.DouyinCommentListResponse)
	}
	return p.Success.FastRead(buf, _type, number)
}

func (p *CommentListResult) FastWrite(buf []byte) (n int) {
	if !p.IsSetSuccess() {
		return 0
	}
	return p.Success.FastWrite(buf)
}

func (p *CommentListResult) Size() (n int) {
	if !p.IsSetSuccess() {
		return 0
	}
	return p.Success.Size()
}

func (p *CommentListResult) Marshal(out []byte) ([]byte, error) {
	if !p.IsSetSuccess() {
		return out, fmt.Errorf("No req in CommentListResult")
	}
	return proto.Marshal(p.Success)
}

func (p *CommentListResult) Unmarshal(in []byte) error {
	msg := new(comment.DouyinCommentListResponse)
	if err := proto.Unmarshal(in, msg); err != nil {
		return err
	}
	p.Success = msg
	return nil
}

func (p *CommentListResult) GetSuccess() *comment.DouyinCommentListResponse {
	if !p.IsSetSuccess() {
		return CommentListResult_Success_DEFAULT
	}
	return p.Success
}

func (p *CommentListResult) SetSuccess(x interface{}) {
	p.Success = x.(*comment.DouyinCommentListResponse)
}

func (p *CommentListResult) IsSetSuccess() bool {
	return p.Success != nil
}

type kClient struct {
	c client.Client
}

func newServiceClient(c client.Client) *kClient {
	return &kClient{
		c: c,
	}
}

func (p *kClient) CommentAction(ctx context.Context, Req *comment.DouyinCommentActionRequest) (r *comment.DouyinCommentActionResponse, err error) {
	var _args CommentActionArgs
	_args.Req = Req
	var _result CommentActionResult
	if err = p.c.Call(ctx, "CommentAction", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) CommentList(ctx context.Context, Req *comment.DouyinCommentListRequest) (r *comment.DouyinCommentListResponse, err error) {
	var _args CommentListArgs
	_args.Req = Req
	var _result CommentListResult
	if err = p.c.Call(ctx, "CommentList", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}
