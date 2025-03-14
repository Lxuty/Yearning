package manage

import (
	"Yearning-go/src/handler/commom"
	"Yearning-go/src/model"
	"github.com/cookieY/yee"
	"net/http"
)

type board struct {
	Board string `json:"board"`
}

const BOARD_MESSAGE_SAVE = "公告已保存"

func GeneralPostBoard(c yee.Context) (err error) {
	req := new(board)
	if err = c.Bind(req); err != nil {
		c.Logger().Error(err.Error())
		return c.JSON(http.StatusOK, err.Error())
	}
	model.DB().Model(model.CoreGlobalConfiguration{}).Update(&model.CoreGlobalConfiguration{Board: req.Board})
	return c.JSON(http.StatusOK, commom.SuccessPayLoadToMessage(BOARD_MESSAGE_SAVE))
}

func GeneralGetBoard(c yee.Context) (err error) {
	var board model.CoreGlobalConfiguration
	model.DB().Select("board").First(&board)
	return c.JSON(http.StatusOK, commom.SuccessPayload(board.Board))
}
