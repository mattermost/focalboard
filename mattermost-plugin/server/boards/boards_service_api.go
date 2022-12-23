package boards

import (
	"github.com/mattermost/focalboard/server/app"
	"github.com/mattermost/focalboard/server/model"

	mm_model "github.com/mattermost/mattermost-server/v6/model"
	"github.com/mattermost/mattermost-server/v6/product"
)

// boardsServiceAPI provides a service API for other products such as Channels.
type boardsServiceAPI struct {
	app *app.App
}

func NewBoardsServiceAPI(app *BoardsApp) *boardsServiceAPI {
	return &boardsServiceAPI{
		app: app.server.App(),
	}
}

func (bs *boardsServiceAPI) GetTemplates(teamID string, userID string) ([]*model.Board, error) {
	return bs.app.GetTemplateBoards(teamID, userID)
}

func (bs *boardsServiceAPI) GetBoard(boardID string) (*model.Board, error) {
	return bs.app.GetBoard(boardID)
}

func (bs *boardsServiceAPI) CreateBoard(board *model.Board, userID string, addmember bool) (*model.Board, error) {
	return bs.app.CreateBoard(board, userID, addmember)
}

func (bs *boardsServiceAPI) PatchBoard(boardPatch *model.BoardPatch, boardID string, userID string) (*model.Board, error) {
	return bs.app.PatchBoard(boardPatch, boardID, userID)
}

func (bs *boardsServiceAPI) DeleteBoard(boardID string, userID string) error {
	return bs.app.DeleteBoard(boardID, userID)
}

func (bs *boardsServiceAPI) SearchBoards(searchTerm string, searchField model.BoardSearchField,
	userID string, includePublicBoards bool) ([]*model.Board, error) {
	return bs.app.SearchBoardsForUser(searchTerm, searchField, userID, includePublicBoards)
}

func (bs *boardsServiceAPI) LinkBoardToChannel(boardID string, channelID string, userID string) (*model.Board, error) {
	patch := &model.BoardPatch{
		ChannelID: &channelID,
	}
	return bs.app.PatchBoard(patch, boardID, userID)
}

func (bs *boardsServiceAPI) GetCards(boardID string) ([]*model.Card, error) {
	return bs.app.GetCardsForBoard(boardID, 0, 0)
}

func (bs *boardsServiceAPI) GetCard(cardID string) (*model.Card, error) {
	return bs.app.GetCardByID(cardID)
}

func (bs *boardsServiceAPI) CreateCard(card *model.Card, boardID string, userID string) (*model.Card, error) {
	return bs.app.CreateCard(card, boardID, userID, false)
}

func (bs *boardsServiceAPI) PatchCard(cardPatch *model.CardPatch, cardID string, userID string) (*model.Card, error) {
	return bs.app.PatchCard(cardPatch, cardID, userID, false)
}

func (bs *boardsServiceAPI) DeleteCard(cardID string, userID string) error {
	return bs.app.DeleteBlock(cardID, userID)
}

func (bs *boardsServiceAPI) HasPermissionToBoard(userID, boardID string, permission *mm_model.Permission) bool {
	return bs.app.HasPermissionToBoard(userID, boardID, permission)
}

// Ensure boardsServiceAPI implements product.BoardsService interface.
var _ product.BoardsService = (*boardsServiceAPI)(nil)