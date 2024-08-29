package v1

import (
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"time"
	"todolist_api/internal/service"
)

type taskRouter struct {
	task service.Task
}

func newTaskRouter(g *echo.Group, task service.Task) {
	r := &taskRouter{task: task}

	g.POST("", r.create)
	g.GET("", r.list)
	g.GET("/:id", r.findById)
	g.PUT("/:id", r.update)
	g.DELETE("/:id", r.delete)
}

type taskCreateInput struct {
	Title       string    `json:"title" validate:"required"`
	Description string    `json:"description" validate:"required"`
	DueDate     time.Time `json:"due_date" validate:"required"`
}

//	@Summary		Create task
//	@Description	Create task
//	@Tags			task
//	@Accept			json
//	@Produce		json
//	@Param			input	body		taskCreateInput	true	"input"
//	@Success		201		{object}	service.TaskOutput
//	@Failure		400		{object}	echo.HTTPError
//	@Failure		500		{object}	echo.HTTPError
//	@Security		JWT
//	@Router			/api/v1/tasks [post]
func (r *taskRouter) create(c echo.Context) error {
	var input taskCreateInput

	if err := c.Bind(&input); err != nil {
		errorResponse(c, http.StatusBadRequest, echo.ErrBadRequest)
		return nil
	}
	if err := c.Validate(&input); err != nil {
		errorResponse(c, http.StatusBadRequest, err)
		return nil
	}

	username, ok := c.Get(usernameCtx).(string)
	if !ok {
		errorResponse(c, http.StatusInternalServerError, echo.ErrInternalServerError)
		return nil
	}

	task, err := r.task.Create(c.Request().Context(), service.TaskCreateInput{
		Username:    username,
		Title:       input.Title,
		Description: input.Description,
		DueDate:     input.DueDate,
	})
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			errorResponse(c, http.StatusBadRequest, err)
			return nil
		}
		errorResponse(c, http.StatusInternalServerError, echo.ErrInternalServerError)
		return err
	}
	return c.JSON(http.StatusCreated, task)
}

//	@Summary		Get tasks
//	@Description	Get list of all user tasks
//	@Tags			task
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		service.TaskOutput
//	@Failure		400	{object}	echo.HTTPError
//	@Failure		500	{object}	echo.HTTPError
//	@Security		JWT
//	@Router			/api/v1/tasks [get]
func (r *taskRouter) list(c echo.Context) error {
	username, ok := c.Get(usernameCtx).(string)
	if !ok {
		errorResponse(c, http.StatusInternalServerError, echo.ErrInternalServerError)
		return nil
	}

	tasks, err := r.task.Find(c.Request().Context(), username)
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, echo.ErrInternalServerError)
		return err
	}
	return c.JSON(http.StatusOK, tasks)
}

//	@Summary		Get task
//	@Description	Get user task by id
//	@Tags			task
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"id"
//	@Success		200	{object}	service.TaskOutput
//	@Failure		400	{object}	echo.HTTPError
//	@Failure		404	{object}	echo.HTTPError
//	@Failure		500	{object}	echo.HTTPError
//	@Security		JWT
//	@Router			/api/v1/tasks/{id} [get]
func (r *taskRouter) findById(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errorResponse(c, http.StatusBadRequest, echo.ErrBadRequest)
		return nil
	}
	username, ok := c.Get(usernameCtx).(string)
	if !ok {
		errorResponse(c, http.StatusInternalServerError, echo.ErrInternalServerError)
		return nil
	}
	task, err := r.task.FindById(c.Request().Context(), id, username)
	if err != nil {
		if errors.Is(err, service.ErrTaskNotFound) {
			errorResponse(c, http.StatusNotFound, err)
			return nil
		}
		errorResponse(c, http.StatusInternalServerError, echo.ErrInternalServerError)
		return err
	}
	return c.JSON(http.StatusOK, task)
}

type taskUpdateInput struct {
	Title       string    `json:"title" validate:"required"`
	Description string    `json:"description" validate:"required"`
	DueDate     time.Time `json:"due_date" validate:"required"`
}

//	@Summary		Update task
//	@Description	Update user task by id
//	@Tags			task
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int				true	"id"
//	@Param			input	body		taskUpdateInput	true	"input"
//	@Success		200		{object}	service.TaskOutput
//	@Failure		400		{object}	echo.HTTPError
//	@Failure		404		{object}	echo.HTTPError
//	@Failure		500		{object}	echo.HTTPError
//	@Security		JWT
//	@Router			/api/v1/tasks/{id} [put]
func (r *taskRouter) update(c echo.Context) error {
	var input taskUpdateInput

	if err := c.Bind(&input); err != nil {
		errorResponse(c, http.StatusBadRequest, echo.ErrBadRequest)
		return nil
	}
	if err := c.Validate(&input); err != nil {
		errorResponse(c, http.StatusBadRequest, err)
		return nil
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errorResponse(c, http.StatusBadRequest, echo.ErrBadRequest)
		return nil
	}
	username, ok := c.Get(usernameCtx).(string)
	if !ok {
		errorResponse(c, http.StatusInternalServerError, echo.ErrInternalServerError)
		return nil
	}

	task, err := r.task.Update(c.Request().Context(), service.TaskUpdateInput{
		Id:          id,
		Username:    username,
		Title:       input.Title,
		Description: input.Description,
		DueDate:     input.DueDate,
	})
	if err != nil {
		if errors.Is(err, service.ErrTaskNotFound) {
			errorResponse(c, http.StatusNotFound, err)
			return nil
		}
		errorResponse(c, http.StatusInternalServerError, echo.ErrInternalServerError)
		return err
	}
	return c.JSON(http.StatusOK, task)
}

//	@Summary		Delete task
//	@Description	Delete user task by id
//	@Tags			task
//	@Accept			json
//	@Produce		json
//	@Param			id	path	int	true	"id"
//	@Success		204
//	@Failure		400	{object}	echo.HTTPError
//	@Failure		404	{object}	echo.HTTPError
//	@Failure		500	{object}	echo.HTTPError
//	@Security		JWT
//	@Router			/api/v1/tasks/{id} [delete]
func (r *taskRouter) delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		errorResponse(c, http.StatusBadRequest, echo.ErrBadRequest)
		return nil
	}
	username, ok := c.Get(usernameCtx).(string)
	if !ok {
		errorResponse(c, http.StatusInternalServerError, echo.ErrInternalServerError)
		return nil
	}
	if err = r.task.Delete(c.Request().Context(), id, username); err != nil {
		if errors.Is(err, service.ErrTaskNotFound) {
			errorResponse(c, http.StatusNotFound, err)
			return nil
		}
		errorResponse(c, http.StatusInternalServerError, echo.ErrInternalServerError)
		return err
	}
	return c.NoContent(http.StatusNoContent)
}
