package customerr

import (
	"context"
	"fmt"
	"log"

	"github.com/cockroachdb/errors"
	"github.com/cockroachdb/errors/errbase"
	"github.com/cockroachdb/errors/markers"
	"github.com/gogo/protobuf/proto"
	officalproto "google.golang.org/protobuf/proto"
)

type WithCustom struct {
	cause  error
	custom string
}

func WrapWithCustom(err error, custom string) error {
	if err == nil {
		return nil
	}
	return &WithCustom{
		cause:  err,
		custom: custom,
	}
}

func GetCustomStr(err error, defaultStr string) string {
	if v, ok := markers.If(err, func(err error) (interface{}, bool) {
		if w, ok := err.(*WithCustom); ok {
			return w.custom, true
		}
		return nil, false
	}); ok {
		return v.(string)
	}
	return defaultStr
}

func (w *WithCustom) Error() string { return w.cause.Error() }

func (w *WithCustom) Cause() error  { return w.cause }
func (w *WithCustom) Unwrap() error { return w.cause }

func (w *WithCustom) Format(s fmt.State, verb rune) { errors.FormatError(w, s, verb) }

func (w *WithCustom) SafeFormatError(p errors.Printer) (next error) {
	if p.Detail() {
		p.Print("custom string %s", w.custom)
	}
	return w.cause
}

func encodeWithCustom(_ context.Context, err error) (string, []string, proto.Message) {
	w := err.(*WithCustom)
	details := []string{fmt.Sprintf("Custom string %s", w.custom)}
	payload := &EncodedCustomErr{Custom: w.custom}

	m, _ := officalproto.Marshal(payload)
	log.Println(m)
	log.Println(string(m))
	newPayload := &EncodedCustomErr{}
	officalproto.Unmarshal(m, newPayload)
	log.Printf("%+v", newPayload)
	return "", details, payload
}

func decodeWithCustom(
	_ context.Context, cause error, _ string, _ []string, payload proto.Message,
) error {
	wp, ok := payload.(*EncodedCustomErr)
	if !ok {
		return &WithCustom{cause: cause, custom: "123123"}
	}
	return &WithCustom{cause: cause, custom: wp.Custom}
}

func init() {
	errbase.RegisterWrapperEncoder(errbase.GetTypeKey((*WithCustom)(nil)), encodeWithCustom)
	errbase.RegisterWrapperDecoder(errbase.GetTypeKey((*WithCustom)(nil)), decodeWithCustom)
}
