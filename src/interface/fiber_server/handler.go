package fiber_server

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/wayla99/go_gorm.git/src/entity/staff"
	"github.com/wayla99/go_gorm.git/src/use_case"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const (
	OK             = "OK"
	xContentLength = "X-Content-Length"
)

type ErrorResponse struct {
	Error     string `json:"error"`
	ErrorCode int    `json:"error_code"`
	IssueId   string `json:"issue_id"`
} //@Name ErrorResponse

type SuccessResp struct {
	Status string      `json:"status"`
	Code   int         `json:"code"`
	Data   interface{} `json:"data"`
} //@Name SuccessResp

func (f *FiberServer) sendError(c *fiber.Ctx, status int, err error, errCode int, issueId string) error {
	if issueId == "" {
		span := trace.SpanFromContext(getSpanContext(c))
		issueId = span.SpanContext().TraceID().String()
	}

	return c.Status(status).JSON(ErrorResponse{
		Error:     err.Error(),
		ErrorCode: errCode,
		IssueId:   issueId,
	})
}

func (f *FiberServer) errorHandler(c *fiber.Ctx, err error, issueIds ...string) error {
	unwrapErr := errors.Unwrap(err)
	if unwrapErr == nil {
		unwrapErr = err
	}

	var issueId string
	if len(issueIds) > 0 {
		issueId = issueIds[0]
	}

	_, span := tracer.Start(getSpanContext(c), "errorHandler")
	defer span.End()
	span.SetStatus(codes.Error, err.Error())

	switch unwrapErr {
	case ErrUnauthenticated:
		return f.sendError(c, 401, err, 1, issueId)
	case ErrInvalidPayload:
		return f.sendError(c, 400, err, 2, issueId)
	case ErrInvalidParameter:
		return f.sendError(c, 400, err, 3, issueId)
	case staff.ErrInvalidStaff:
		return f.sendError(c, 400, err, 4, issueId)
	case use_case.ErrStaffNotFound:
		return f.sendError(c, 404, err, 5, issueId)
	case nil:
		return f.sendError(c, 500, errors.New("nil error"), -1, issueId)
	default:
		// Unknown error
		return f.sendError(c, 500, err, 0, issueId)
	}
}

func (f *FiberServer) successHandler(c *fiber.Ctx, status int, data interface{}) error {
	return c.Status(status).JSON(SuccessResp{
		Status: OK,
		Code:   status,
		Data:   data,
	})
}

func (f *FiberServer) paginatorHandler(c *fiber.Ctx, total int, items interface{}) error {
	status := http.StatusOK
	if total < 1 {
		status = http.StatusNoContent
	}
	c.GetRespHeader(xContentLength, strconv.Itoa(total))
	return f.successHandler(c, status, items)
}
