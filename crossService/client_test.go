package crossservice

import (
	"context"
	fmt "fmt"
	"strings"
	"testing"
	"tryCockroachdbErr/customerr"

	"github.com/cockroachdb/errors"

	"github.com/cockroachdb/errors/exthttp"
	"github.com/cockroachdb/errors/grpc/status"
	"github.com/cockroachdb/errors/testutils"
	codes "google.golang.org/grpc/codes"
)

func TestGrpc(t *testing.T) {
	tt := testutils.T{T: t}

	resp, err := Client.Echo(context.Background(), &EchoRequest{Text: "hello"})
	tt.Assert(err == nil)
	tt.Assert(resp.Reply == "echoing: hello")
	tt.Assert(status.Code(err) == codes.OK)

	_, err = Client.Echo(context.Background(), &EchoRequest{Text: "noecho"})
	tt.Assert(err != nil)
	tt.Assert(errors.Is(err, ErrCantEcho))
	tt.Assert(status.Code(err) == codes.Unknown)

	_, err = Client.Echo(context.Background(), &EchoRequest{Text: "really_long_message"})
	tt.Assert(err != nil)
	tt.Assert(err.Error() == "really_long_message is too long: text is too long")
	tt.Assert(errors.Is(err, ErrTooLong))
	tt.Assert(errors.UnwrapAll(err).Error() == "text is too long")

	_, err = Client.Echo(context.Background(), &EchoRequest{Text: "internal"})
	tt.Assert(err != nil)
	tt.Assert(err.Error() == "there was a problem: internal error!")
	tt.Assert(status.Code(err) == codes.Internal)
	tt.Assert(errors.Is(err, ErrInternal))
	spv := fmt.Sprintf("%+v", err)
	t.Logf("spv:\n%s", spv)
	tt.Assert(strings.Contains(spv, "gRPC code: Internal"))
	fmt.Println("--------------------------")

	_, err = Client.Echo(context.Background(), &EchoRequest{Text: "internal2"})
	tt.Assert(err != nil)
	tt.Assert(errors.Is(err, ErrInternal))
	fmt.Println(errors.GetAllHints(err))
	// tt.Assert(errors.Is(err, ErrTooLong))

	_, err = Client.Echo(context.Background(), &EchoRequest{Text: "httpCode"})
	tt.Assert(exthttp.GetHTTPCode(err, 123) == 500)

	_, err = Client.Echo(context.Background(), &EchoRequest{Text: "custom"})
	tt.Assert(err != nil)
	tt.Assert(errors.Is(err, ErrInternal))
	tt.Assert(errors.As(err, new(*customerr.WithCustom)))
	tt.Assert(customerr.GetCustomStr(err, "not thing") == "123123")
}
