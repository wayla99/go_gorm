package fiber_server

import (
	"net/http"

	"github.com/wayla99/go_gorm.git/src/use_case"

	"github.com/gofiber/fiber/v2"
)

func (f *FiberServer) addRouteStaff(r fiber.Router) {
	r.Get("/", f.getStaffs)
	r.Post("/", f.createStaff)
	r.Get("/:staff_id", f.getStaffById)
	r.Put("/:staff_id", f.updateStaffById)
	r.Delete("/:staff_id", f.deleteStaffById)
}

// getStaffs godoc
// @Summary get staffs
// @Description return rows of staff
// @Tags Staffs
// @Security X-User-Headers
// @Accept  json
// @Produce  json
// @param offset query number false "offset number"
// @param limit query number false "limit number"
// @Param sorts query []string false "Sort for staff data, ex.`created_at:asc`"
// @Param filters query []string false "Filters for staff data, ex. `first_name:like:test`" collectionFormat(multi)
// @Success 200 {object} SuccessResp{data=[]fiber_server.Staff}
// @Failure 400 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /staffs [get]
func (f *FiberServer) getStaffs(ctx *fiber.Ctx) error {
	input := &use_case.List{}
	if err := ctx.QueryParser(input); err != nil {
		return f.errorHandler(ctx, err)
	}
	if input.Offset == 0 {
		input.Offset = 1
	}
	if input.Limit == 0 {
		input.Limit = 10
	}

	staff, total, err := f.useCase.GetStaffs(getSpanContext(ctx), input)
	if err != nil {
		return f.errorHandler(ctx, err)
	}

	return f.paginatorHandler(ctx, total, staff)

}

// createStaff godoc
// @Summary create staffs
// @Description return array of created id
// @Tags Staffs
// @Security X-User-Headers
// @Accept  json
// @Produce  json
// @Param data body Staff true "The input staff struct"
// @Success 201 {string} SuccessResp{data=[]fiber_server.Staff}
// @Failure 400 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /staffs [post]
func (f *FiberServer) createStaff(ctx *fiber.Ctx) error {
	var staff Staff
	if err := ctx.BodyParser(&staff); err != nil {
		return f.errorHandler(ctx, ErrInvalidPayload)
	}

	sf, err := f.useCase.CreateStaff(getSpanContext(ctx), staff.toUseCase())
	if err != nil {
		return f.errorHandler(ctx, err)
	}

	return f.successHandler(ctx, http.StatusCreated, sf)
}

// getStaffById godoc
// @Summary get staff by id
// @Description return a row of staff
// @Tags Staffs
// @Security X-User-Headers
// @Accept  json
// @Produce  json
// @Param staff_id path string true "staff id of staff to be fetched"
// @Success 200 {object} SuccessResp{data=[]fiber_server.Staff}
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /staffs/{staff_id} [get]
func (f *FiberServer) getStaffById(ctx *fiber.Ctx) error {
	sf, err := f.useCase.GetStaffById(
		getSpanContext(ctx),
		ctx.Params("staff_id"),
	)

	if err != nil {
		return f.errorHandler(ctx, err)
	}

	return f.successHandler(ctx, http.StatusOK, sf)
}

// updateStaffById godoc
// @Summary update staff
// @Description return OK
// @Tags Staffs
// @Security X-User-Headers
// @Accept  json
// @Produce  json
// @Param staff_id path string true "staff id of staff to be updated"
// @Param data body Staff true "The input staff struct"
// @Success 200 {string} SuccessResp{data=[]fiber_server.Staff}
// @Failure 400 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /staffs/{staff_id} [put]
func (f *FiberServer) updateStaffById(ctx *fiber.Ctx) error {
	var staff Staff
	if err := ctx.BodyParser(&staff); err != nil {
		return f.errorHandler(ctx, ErrInvalidPayload)
	}

	sf, err := f.useCase.UpdateStaffById(
		getSpanContext(ctx),
		ctx.Params("staff_id"),
		staff.toUseCase(),
	)

	if err != nil {
		return f.errorHandler(ctx, err)
	}

	return f.successHandler(ctx, http.StatusOK, sf)
}

// deleteStaffById godoc
// @Summary delete staff
// @Description return OK
// @Tags Staffs
// @Security X-User-Headers
// @Accept  json
// @Produce  json
// @Param staff_id path string true "staff id of staff to be deleted"
// @Success 200 {string} SuccessResp{data=[]fiber_server.Staff}
// @Failure 400 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /staffs/{staff_id} [delete]
func (f *FiberServer) deleteStaffById(ctx *fiber.Ctx) error {
	if err := f.useCase.DeleteStaffById(getSpanContext(ctx), ctx.Params("staff_id")); err != nil {
		return f.errorHandler(ctx, err)
	}

	return f.successHandler(ctx, http.StatusOK, OK)
}
