package opentaobao

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func RespDecode(httpResp *http.Response, v any) error {
	resp := &Response{
		Data: v,
	}
	err := json.NewDecoder(httpResp.Body).Decode(resp)
	if err != nil {
		return err
	}
	fmt.Printf("resp: %+v\n", resp)
	fmt.Printf("resp.Data: %+v\n", resp.Data)
	if resp.IsError() {
		return resp.Error()
	}
	if httpResp.StatusCode != http.StatusOK {
		return fmt.Errorf("http status code: %d", httpResp.StatusCode)
	}
	return nil
}

// Error HTTP响应错误项
type Error struct {
	Code    int    `json:"code" swaggo:"true,错误码"`
	Message string `json:"message" swaggo:"true,错误信息"`
}

type Response struct {
	Err  Error `json:"error" swaggo:"true,错误项"`
	Data any   `json:"data"`
}

func (e *Response) IsError() bool {
	return e.Err.Code != 200 && e.Err.Code != 0
}

func (e *Response) Error() error {
	return errors.New(e.Err.Message)
}
