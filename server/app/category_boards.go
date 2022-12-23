package app

import (
	"errors"
	"fmt"

	"github.com/mattermost/focalboard/server/model"
)

const defaultCategoryBoards = "Boards"
const defaultCategoryPages = "Pages"

var errCategoryBoardsLengthMismatch = errors.New("cannot update category boards order, passed list of categories boards different size than in database")
var errBoardNotFoundInCategory = errors.New("specified board ID not found in specified category ID")
var errBoardMembershipNotFound = errors.New("board membership not found for user's board")

func (a *App) GetUserCategoryBoards(userID, teamID string) ([]model.CategoryBoards, error) {
	categories, err := a.store.GetUserCategoryBoards(userID, teamID)
	if err != nil {
		return nil, err
	}

	categoryBoards := []model.CategoryBoards{}
	categoryPages := []model.CategoryBoards{}
	for _, category := range categories {
		if category.Type == model.CategoryTypeSystem || category.Type == model.CategoryTypeCustom {
			categoryBoards = append(categoryBoards, category)
		} else if category.Type == model.CategoryTypePagesSystem || category.Type == model.CategoryTypePagesCustom {
			categoryPages = append(categoryPages, category)
		}
	}

	createdCategoryBoards, err := a.createDefaultCategoriesIfRequired(categoryBoards, categoryPages, userID, teamID)
	if err != nil {
		return nil, err
	}

	resultCategories := append(categoryBoards, categoryPages...)
	resultCategories = append(resultCategories, createdCategoryBoards...)
	return resultCategories, nil
}

func (a *App) createDefaultCategoriesIfRequired(existingCategoryBoards []model.CategoryBoards, existingCategoryPages []model.CategoryBoards, userID, teamID string) ([]model.CategoryBoards, error) {
	createdCategories := []model.CategoryBoards{}
	fmt.Println("BOARDS", existingCategoryBoards)
	fmt.Println("PAGES", existingCategoryPages)

	boardsCategoryExist := false
	for _, categoryBoard := range existingCategoryBoards {
		if categoryBoard.Type == model.CategoryTypeSystem && categoryBoard.Name == defaultCategoryBoards {
			boardsCategoryExist = true
		}
	}

	if !boardsCategoryExist {
		createdCategoryBoards, err := a.createBoardsCategory(userID, teamID, existingCategoryBoards)
		if err != nil {
			return nil, err
		}

		createdCategories = append(createdCategories, *createdCategoryBoards)
	}

	pagesCategoryExist := false
	for _, categoryPage := range existingCategoryPages {
		if categoryPage.Type == model.CategoryTypePagesSystem && categoryPage.Name == defaultCategoryPages {
			pagesCategoryExist = true
		}
	}

	if !pagesCategoryExist {
		createdCategoryPages, err := a.createPagesCategory(userID, teamID, existingCategoryPages)
		if err != nil {
			return nil, err
		}

		createdCategories = append(createdCategories, *createdCategoryPages)
	}

	return createdCategories, nil
}

func (a *App) createBoardsCategory(userID, teamID string, existingCategoryBoards []model.CategoryBoards) (*model.CategoryBoards, error) {
	// create the category
	category := model.Category{
		Name:      defaultCategoryBoards,
		UserID:    userID,
		TeamID:    teamID,
		Collapsed: false,
		Type:      model.CategoryTypeSystem,
		SortOrder: len(existingCategoryBoards) * model.CategoryBoardsSortOrderGap,
	}
	createdCategory, err := a.CreateCategory(&category)
	if err != nil {
		return nil, fmt.Errorf("createBoardsCategory default category creation failed: %w", err)
	}

	// once the category is created, we need to move all boards which do not
	// belong to any category, into this category.

	boardMembers, err := a.GetMembersForUser(userID)
	if err != nil {
		return nil, fmt.Errorf("createBoardsCategory error fetching user's board memberships: %w", err)
	}

	boardMemberByBoardID := map[string]*model.BoardMember{}
	for _, boardMember := range boardMembers {
		boardMemberByBoardID[boardMember.BoardID] = boardMember
	}

	createdCategoryBoards := &model.CategoryBoards{
		Category: *createdCategory,
		BoardIDs: []string{},
	}

	// get user's current team's baords
	userTeamBoards, err := a.GetBoardsForUserAndTeam(userID, teamID, false)
	if err != nil {
		return nil, fmt.Errorf("createBoardsCategory error fetching user's team's boards: %w", err)
	}

	for _, board := range userTeamBoards {
		boardMembership, ok := boardMemberByBoardID[board.ID]
		if !ok {
			return nil, fmt.Errorf("createBoardsCategory: %w", errBoardMembershipNotFound)
		}

		// boards with implicit access (aka synthetic membership),
		// should show up in LHS only when openign them explicitelly.
		// So we don't process any synthetic membership boards
		// and only add boards with explicit access to, to the the LHS,
		// for example, if a user explicitelly added another user to a board.
		if boardMembership.Synthetic {
			continue
		}

		if board.IsPagesFolder {
			continue
		}

		belongsToCategory := false

		for _, categoryBoard := range existingCategoryBoards {
			if categoryBoard.Type != model.CategoryTypeSystem && categoryBoard.Type != model.CategoryTypeCustom {
				for _, boardID := range categoryBoard.BoardIDs {
					if boardID == board.ID {
						belongsToCategory = true
						break
					}
				}
			}

			// stop looking into other categories if
			// the board was found in a category
			if belongsToCategory {
				break
			}
		}

		if !belongsToCategory {
			if err := a.AddUpdateUserCategoryBoard(teamID, userID, map[string]string{board.ID: createdCategory.ID}); err != nil {
				return nil, fmt.Errorf("createBoardsCategory failed to add category-less board to the default category, defaultCategoryID: %s, error: %w", createdCategory.ID, err)
			}

			createdCategoryBoards.BoardIDs = append(createdCategoryBoards.BoardIDs, board.ID)
		}
	}

	return createdCategoryBoards, nil
}

func (a *App) createPagesCategory(userID, teamID string, existingCategoryBoards []model.CategoryBoards) (*model.CategoryBoards, error) {
	// create the category
	category := model.Category{
		Name:      defaultCategoryPages,
		UserID:    userID,
		TeamID:    teamID,
		Collapsed: false,
		Type:      model.CategoryTypePagesSystem,
		SortOrder: len(existingCategoryBoards) * model.CategoryBoardsSortOrderGap,
	}
	createdCategory, err := a.CreateCategory(&category)
	if err != nil {
		return nil, fmt.Errorf("createBoardsCategory default category creation failed: %w", err)
	}

	// once the category is created, we need to move all boards which do not
	// belong to any category, into this category.

	boardMembers, err := a.GetMembersForUser(userID)
	if err != nil {
		return nil, fmt.Errorf("createBoardsCategory error fetching user's board memberships: %w", err)
	}

	boardMemberByBoardID := map[string]*model.BoardMember{}
	for _, boardMember := range boardMembers {
		boardMemberByBoardID[boardMember.BoardID] = boardMember
	}

	createdCategoryBoards := &model.CategoryBoards{
		Category: *createdCategory,
		BoardIDs: []string{},
	}

	// get user's current team's baords
	userTeamBoards, err := a.GetBoardsForUserAndTeam(userID, teamID, false)
	if err != nil {
		return nil, fmt.Errorf("createBoardsCategory error fetching user's team's boards: %w", err)
	}

	for _, board := range userTeamBoards {
		boardMembership, ok := boardMemberByBoardID[board.ID]
		if !ok {
			return nil, fmt.Errorf("createBoardsCategory: %w", errBoardMembershipNotFound)
		}

		// boards with implicit access (aka synthetic membership),
		// should show up in LHS only when openign them explicitelly.
		// So we don't process any synthetic membership boards
		// and only add boards with explicit access to, to the the LHS,
		// for example, if a user explicitelly added another user to a board.
		if boardMembership.Synthetic {
			continue
		}

		if !board.IsPagesFolder {
			continue
		}

		belongsToCategory := false

		for _, categoryBoard := range existingCategoryBoards {
			if categoryBoard.Type != model.CategoryTypePagesSystem && categoryBoard.Type != model.CategoryTypePagesCustom {
				for _, boardID := range categoryBoard.BoardIDs {
					if boardID == board.ID {
						belongsToCategory = true
						break
					}
				}
			}

			// stop looking into other categories if
			// the board was found in a category
			if belongsToCategory {
				break
			}
		}

		if !belongsToCategory {
			if err := a.AddUpdateUserCategoryBoard(teamID, userID, map[string]string{board.ID: createdCategory.ID}); err != nil {
				return nil, fmt.Errorf("createBoardsCategory failed to add category-less board to the default category, defaultCategoryID: %s, error: %w", createdCategory.ID, err)
			}

			createdCategoryBoards.BoardIDs = append(createdCategoryBoards.BoardIDs, board.ID)
		}
	}

	return createdCategoryBoards, nil
}

func (a *App) AddUpdateUserCategoryBoard(teamID, userID string, boardCategoryMapping map[string]string) error {
	err := a.store.AddUpdateCategoryBoard(userID, boardCategoryMapping)
	if err != nil {
		return err
	}

	wsPayload := make([]*model.BoardCategoryWebsocketData, len(boardCategoryMapping))
	i := 0
	for boardID, categoryID := range boardCategoryMapping {
		wsPayload[i] = &model.BoardCategoryWebsocketData{
			BoardID:    boardID,
			CategoryID: categoryID,
		}
		i++
	}

	a.blockChangeNotifier.Enqueue(func() error {
		a.wsAdapter.BroadcastCategoryBoardChange(
			teamID,
			userID,
			wsPayload,
		)
		return nil
	})

	return nil
}

func (a *App) ReorderCategoryBoards(userID, teamID, categoryID string, newBoardsOrder []string) ([]string, error) {
	if err := a.verifyNewCategoryBoardsMatchExisting(userID, teamID, categoryID, newBoardsOrder); err != nil {
		return nil, err
	}

	newOrder, err := a.store.ReorderCategoryBoards(categoryID, newBoardsOrder)
	if err != nil {
		return nil, err
	}

	go func() {
		a.wsAdapter.BroadcastCategoryBoardsReorder(teamID, userID, categoryID, newOrder)
	}()

	return newOrder, nil
}

func (a *App) verifyNewCategoryBoardsMatchExisting(userID, teamID, categoryID string, newBoardsOrder []string) error {
	// this function is to ensure that we don't miss specifying
	// all boards of the category while reordering.
	existingCategoryBoards, err := a.GetUserCategoryBoards(userID, teamID)
	if err != nil {
		return err
	}

	var targetCategoryBoards *model.CategoryBoards
	for i := range existingCategoryBoards {
		if existingCategoryBoards[i].Category.ID == categoryID {
			targetCategoryBoards = &existingCategoryBoards[i]
			break
		}
	}

	if targetCategoryBoards == nil {
		return fmt.Errorf("%w categoryID: %s", errCategoryNotFound, categoryID)
	}

	if len(targetCategoryBoards.BoardIDs) != len(newBoardsOrder) {
		return fmt.Errorf(
			"%w length new category boards: %d, length existing category boards: %d, userID: %s, teamID: %s, categoryID: %s",
			errCategoryBoardsLengthMismatch,
			len(newBoardsOrder),
			len(targetCategoryBoards.BoardIDs),
			userID,
			teamID,
			categoryID,
		)
	}

	existingBoardMap := map[string]bool{}
	for _, boardID := range targetCategoryBoards.BoardIDs {
		existingBoardMap[boardID] = true
	}

	for _, boardID := range newBoardsOrder {
		if _, found := existingBoardMap[boardID]; !found {
			return fmt.Errorf(
				"%w board ID: %s, category ID: %s, userID: %s, teamID: %s",
				errBoardNotFoundInCategory,
				boardID,
				categoryID,
				userID,
				teamID,
			)
		}
	}

	return nil
}
