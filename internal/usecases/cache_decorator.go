package usecases

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/DKhorkov/libs/cache"
	"github.com/DKhorkov/libs/logging"
	"github.com/rxwycdh/rxhash"

	"github.com/DKhorkov/hmtm-bff/internal/entities"
	"github.com/DKhorkov/hmtm-bff/internal/interfaces"
)

const (
	// Prefixes:
	getUserByIDPrefix       = "get_user_by_id"
	getUserByEmailPrefix    = "get_user_by_email"
	getUsersPrefix          = "users"
	getToysPrefix           = "toys"
	countToysPrefix         = "toys_count"
	countMastersPrefix      = "masters_count"
	getMasterToysPrefix     = "master_toys"
	countMasterToysPrefix   = "master_toys_count"
	getToyByIDPrefix        = "get_toy_by_id"
	getMastersPrefix        = "get_masters"
	getMasterByIDPrefix     = "get_master_by_id"
	getMasterByUserIDPrefix = "get_master_by_user_id"
	getCategoriesPrefix     = "get_categories"
	getCategoryByIDPrefix   = "get_category_by_id"
	getTagsPrefix           = "get_tags"
	getTagByIDPrefix        = "get_tag_by_id"
	getTicketByIDPrefix     = "get_ticket_by_id"
	getTicketsPrefix        = "get_tickets"
	getUserTicketsPrefix    = "get_user_tickets"
	countTicketsPrefix      = "tickets_count"
	countUserTicketsPrefix  = "user_tickets_count"

	// TTls:
	getUserByIDTTL       = time.Hour * 24
	getUserByEmailTTL    = time.Hour * 24
	getUsersTTL          = time.Minute * 5
	getToysTTL           = time.Minute * 5
	countToysTTL         = time.Minute * 5
	countMastersTTL      = time.Minute * 5
	getMasterToysTTL     = time.Hour * 6
	countMasterToysTTL   = time.Hour * 6
	getToyByIDTTL        = time.Hour * 24
	getMastersTTL        = time.Hour * 6
	getMasterByIDTTL     = time.Hour * 24
	getMasterByUserIDTTL = time.Hour * 24
	getCategoriesTTL     = time.Hour * 24
	getCategoryByIDTTL   = time.Hour * 24
	getTagsTTL           = time.Hour * 24
	getTagByIDTTL        = time.Hour * 24
	getTicketByIDTTL     = time.Hour * 24
	getTicketsTTL        = time.Minute * 5
	getUserTicketsTTL    = time.Minute * 5
	countTicketsTTL      = time.Minute * 5
	countUserTicketsTTL  = time.Hour * 6
)

func NewCacheDecorator(
	useCases interfaces.UseCases,
	cacheProvider cache.Provider,
	logger logging.Logger,
) *CacheDecorator {
	return &CacheDecorator{
		UseCases:      useCases,
		logger:        logger,
		cacheProvider: cacheProvider,
	}
}

type CacheDecorator struct {
	interfaces.UseCases
	cacheProvider cache.Provider
	logger        logging.Logger
}

func (c *CacheDecorator) GetUserByID(ctx context.Context, id uint64) (*entities.User, error) {
	if _, err := c.cacheProvider.Ping(ctx); err != nil {
		logging.LogErrorContext(
			ctx,
			c.logger,
			"Cache provider error",
			err,
		)

		return c.UseCases.GetUserByID(ctx, id)
	}

	cacheKey := fmt.Sprintf("%s:%d", getUserByIDPrefix, id)
	userToDecode, err := c.cacheProvider.Get(ctx, cacheKey)
	if err == nil {
		var user entities.User
		if err = json.Unmarshal([]byte(userToDecode), &user); err == nil {
			return &user, nil
		}
	}

	logging.LogErrorContext(
		ctx,
		c.logger,
		fmt.Sprintf("Failed to get cached User with id=%d", id),
		err,
	)

	user, err := c.UseCases.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	encodedUser, err := json.Marshal(user)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			c.logger,
			fmt.Sprintf("Failed to cache User with id=%d", id),
			err,
		)

		return user, nil
	}

	if err = c.cacheProvider.Set(ctx, cacheKey, encodedUser, getUserByIDTTL); err != nil {
		logging.LogErrorContext(
			ctx,
			c.logger,
			fmt.Sprintf("Failed to cache User with id=%d", id),
			err,
		)
	}

	return user, nil
}

func (c *CacheDecorator) GetUserByEmail(
	ctx context.Context,
	email string,
) (*entities.User, error) {
	if _, err := c.cacheProvider.Ping(ctx); err != nil {
		logging.LogErrorContext(
			ctx,
			c.logger,
			"Cache provider error",
			err,
		)

		return c.UseCases.GetUserByEmail(ctx, email)
	}

	cacheKey := fmt.Sprintf("%s:%s", getUserByEmailPrefix, email)
	userToDecode, err := c.cacheProvider.Get(ctx, cacheKey)
	if err == nil {
		var user entities.User
		if err = json.Unmarshal([]byte(userToDecode), &user); err == nil {
			return &user, nil
		}
	}

	logging.LogErrorContext(
		ctx,
		c.logger,
		fmt.Sprintf("Failed to get cached User with email=%s", email),
		err,
	)

	user, err := c.UseCases.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	encodedUser, err := json.Marshal(user)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			c.logger,
			fmt.Sprintf("Failed to cache User with email=%s", email),
			err,
		)

		return user, nil
	}

	if err = c.cacheProvider.Set(ctx, cacheKey, encodedUser, getUserByEmailTTL); err != nil {
		logging.LogErrorContext(
			ctx,
			c.logger,
			fmt.Sprintf("Failed to cache User with email=%s", email),
			err,
		)
	}

	return user, nil
}

func (c *CacheDecorator) GetUsers(ctx context.Context, pagination *entities.Pagination) ([]entities.User, error) {
	if _, err := c.cacheProvider.Ping(ctx); err != nil {
		logging.LogErrorContext(
			ctx,
			c.logger,
			"Cache provider error",
			err,
		)

		return c.UseCases.GetUsers(ctx, pagination)
	}

	paginationHash, err := rxhash.HashStruct(pagination)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			c.logger,
			"Failed to get cached Users",
			err,
		)

		return c.UseCases.GetUsers(ctx, pagination)
	}

	cacheKey := fmt.Sprintf("%s:%s", getUsersPrefix, paginationHash)
	usersToDecode, err := c.cacheProvider.Get(ctx, cacheKey)
	if err == nil {
		var users []entities.User
		if err = json.Unmarshal([]byte(usersToDecode), &users); err == nil {
			return users, nil
		}
	}

	logging.LogErrorContext(
		ctx,
		c.logger,
		"Failed to get cached Users",
		err,
	)

	users, err := c.UseCases.GetUsers(ctx, pagination)
	if err != nil {
		return nil, err
	}

	encodedUsers, err := json.Marshal(users)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			c.logger,
			"Failed to cache Users",
			err,
		)

		return users, nil
	}

	if err = c.cacheProvider.Set(ctx, cacheKey, encodedUsers, getUsersTTL); err != nil {
		logging.LogErrorContext(
			ctx,
			c.logger,
			"Failed to cache Users",
			err,
		)
	}

	return users, nil
}

func (c *CacheDecorator) GetToys(
	ctx context.Context,
	pagination *entities.Pagination,
	filters *entities.ToysFilters,
) ([]entities.Toy, error) {
	if _, err := c.cacheProvider.Ping(ctx); err != nil {
		logging.LogErrorContext(
			ctx,
			c.logger,
			"Cache provider error",
			err,
		)

		return c.UseCases.GetToys(ctx, pagination, filters)
	}

	paginationHash, err := rxhash.HashStruct(pagination)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			c.logger,
			"Failed to get cached Toys",
			err,
		)

		return c.UseCases.GetToys(ctx, pagination, filters)
	}

	filtersHash, err := rxhash.HashStruct(filters)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			c.logger,
			"Failed to get cached Toys",
			err,
		)

		return c.UseCases.GetToys(ctx, pagination, filters)
	}

	cacheKey := fmt.Sprintf(
		"%s:%s_%s",
		getToysPrefix,
		paginationHash,
		filtersHash,
	)

	toysToDecode, err := c.cacheProvider.Get(ctx, cacheKey)
	if err == nil {
		var toys []entities.Toy
		if err = json.Unmarshal([]byte(toysToDecode), &toys); err == nil {
			return toys, nil
		}
	}

	logging.LogErrorContext(
		ctx,
		c.logger,
		"Failed to get cached Toys",
		err,
	)

	toys, err := c.UseCases.GetToys(ctx, pagination, filters)
	if err != nil {
		return nil, err
	}

	encodedToys, err := json.Marshal(toys)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			c.logger,
			"Failed to cache Toys",
			err,
		)

		return toys, nil
	}

	if err = c.cacheProvider.Set(ctx, cacheKey, encodedToys, getToysTTL); err != nil {
		logging.LogErrorContext(
			ctx,
			c.logger,
			"Failed to cache Toys",
			err,
		)
	}

	return toys, nil
}

func (c *CacheDecorator) CountToys(ctx context.Context, filters *entities.ToysFilters) (uint64, error) {
	if _, err := c.cacheProvider.Ping(ctx); err != nil {
		logging.LogErrorContext(
			ctx,
			c.logger,
			"Cache provider error",
			err,
		)

		return c.UseCases.CountToys(ctx, filters)
	}

	filtersHash, err := rxhash.HashStruct(filters)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			c.logger,
			"Failed to get cached Toys counter",
			err,
		)

		return c.UseCases.CountToys(ctx, filters)
	}

	cacheKey := fmt.Sprintf("%s:%s", countToysPrefix, filtersHash)
	strCounter, err := c.cacheProvider.Get(ctx, cacheKey)
	if err == nil {
		if counter, err := strconv.ParseUint(strCounter, 10, 64); err == nil {
			return counter, nil
		}
	}

	logging.LogErrorContext(
		ctx,
		c.logger,
		"Failed to get cached Toys counter",
		err,
	)

	counter, err := c.UseCases.CountToys(ctx, filters)
	if err != nil {
		return 0, err
	}

	if err = c.cacheProvider.Set(ctx, cacheKey, counter, countToysTTL); err != nil {
		logging.LogErrorContext(
			ctx,
			c.logger,
			"Failed to cache Toys counter",
			err,
		)
	}

	return counter, nil
}

func (c *CacheDecorator) CountMasterToys(
	ctx context.Context,
	masterID uint64,
	filters *entities.ToysFilters,
) (uint64, error) {
	if _, err := c.cacheProvider.Ping(ctx); err != nil {
		logging.LogErrorContext(
			ctx,
			c.logger,
			"Cache provider error",
			err,
		)

		return c.UseCases.CountMasterToys(ctx, masterID, filters)
	}

	filtersHash, err := rxhash.HashStruct(filters)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			c.logger,
			fmt.Sprintf("Failed to get cached Toys counter for Master with id=%d", masterID),
			err,
		)

		return c.UseCases.CountMasterToys(ctx, masterID, filters)
	}

	cacheKey := fmt.Sprintf(
		"%s:%d_%s",
		countMasterToysPrefix,
		masterID,
		filtersHash,
	)

	strCounter, err := c.cacheProvider.Get(ctx, cacheKey)
	if err == nil {
		if counter, err := strconv.ParseUint(strCounter, 10, 64); err == nil {
			return counter, nil
		}
	}

	logging.LogErrorContext(
		ctx,
		c.logger,
		fmt.Sprintf("Failed to get cached Toys counter for Master with id=%d", masterID),
		err,
	)

	counter, err := c.UseCases.CountMasterToys(ctx, masterID, filters)
	if err != nil {
		return 0, err
	}

	if err = c.cacheProvider.Set(ctx, cacheKey, counter, countMasterToysTTL); err != nil {
		logging.LogErrorContext(
			ctx,
			c.logger,
			fmt.Sprintf("Failed to cache Toys counter for Master with id=%d", masterID),
			err,
		)
	}

	return counter, nil
}

func (c *CacheDecorator) GetMasterToys(
	ctx context.Context,
	masterID uint64,
	pagination *entities.Pagination,
	filters *entities.ToysFilters,
) ([]entities.Toy, error) {
	if _, err := c.cacheProvider.Ping(ctx); err != nil {
		logging.LogErrorContext(
			ctx,
			c.logger,
			"Cache provider error",
			err,
		)

		return c.UseCases.GetMasterToys(ctx, masterID, pagination, filters)
	}

	paginationHash, err := rxhash.HashStruct(pagination)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			c.logger,
			fmt.Sprintf("Failed to get cached Toys for Master with id=%d", masterID),
			err,
		)

		return c.UseCases.GetMasterToys(ctx, masterID, pagination, filters)
	}

	filtersHash, err := rxhash.HashStruct(filters)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			c.logger,
			fmt.Sprintf("Failed to get cached Toys for Master with id=%d", masterID),
			err,
		)

		return c.UseCases.GetMasterToys(ctx, masterID, pagination, filters)
	}

	cacheKey := fmt.Sprintf(
		"%s:%d_%s_%s",
		getMasterToysPrefix,
		masterID,
		paginationHash,
		filtersHash,
	)

	toysToDecode, err := c.cacheProvider.Get(ctx, cacheKey)
	if err == nil {
		var toys []entities.Toy
		if err = json.Unmarshal([]byte(toysToDecode), &toys); err == nil {
			return toys, nil
		}
	}

	logging.LogErrorContext(
		ctx,
		c.logger,
		fmt.Sprintf("Failed to get cached Toys for Master with id=%d", masterID),
		err,
	)

	toys, err := c.UseCases.GetMasterToys(ctx, masterID, pagination, filters)
	if err != nil {
		return nil, err
	}

	encodedToys, err := json.Marshal(toys)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			c.logger,
			fmt.Sprintf("Failed to cache Toys for Master with id=%d", masterID),
			err,
		)

		return toys, nil
	}

	if err = c.cacheProvider.Set(ctx, cacheKey, encodedToys, getMasterToysTTL); err != nil {
		logging.LogErrorContext(
			ctx,
			c.logger,
			fmt.Sprintf("Failed to cache Toys for Master with id=%d", masterID),
			err,
		)
	}

	return toys, nil
}

func (c *CacheDecorator) GetToyByID(ctx context.Context, id uint64) (*entities.Toy, error) {
	if _, err := c.cacheProvider.Ping(ctx); err != nil {
		logging.LogErrorContext(
			ctx,
			c.logger,
			"Cache provider error",
			err,
		)

		return c.UseCases.GetToyByID(ctx, id)
	}

	cacheKey := fmt.Sprintf("%s:%d", getToyByIDPrefix, id)
	toyToDecode, err := c.cacheProvider.Get(ctx, cacheKey)
	if err == nil {
		var toy entities.Toy
		if err = json.Unmarshal([]byte(toyToDecode), &toy); err == nil {
			return &toy, nil
		}
	}

	logging.LogErrorContext(
		ctx,
		c.logger,
		fmt.Sprintf("Failed to get cached Toy with id=%d", id),
		err,
	)

	toy, err := c.UseCases.GetToyByID(ctx, id)
	if err != nil {
		return nil, err
	}

	encodedToy, err := json.Marshal(toy)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			c.logger,
			fmt.Sprintf("Failed to cache Toy with id=%d", id),
			err,
		)

		return toy, nil
	}

	if err = c.cacheProvider.Set(ctx, cacheKey, encodedToy, getToyByIDTTL); err != nil {
		logging.LogErrorContext(
			ctx,
			c.logger,
			fmt.Sprintf("Failed to cache Toy with id=%d", id),
			err,
		)
	}

	return toy, nil
}

func (c *CacheDecorator) GetMasters(
	ctx context.Context,
	pagination *entities.Pagination,
	filters *entities.MastersFilters,
) ([]entities.Master, error) {
	if _, err := c.cacheProvider.Ping(ctx); err != nil {
		logging.LogErrorContext(
			ctx,
			c.logger,
			"Cache provider error",
			err,
		)

		return c.UseCases.GetMasters(ctx, pagination, filters)
	}

	paginationHash, err := rxhash.HashStruct(pagination)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			c.logger,
			"Failed to get cached Masters",
			err,
		)

		return c.UseCases.GetMasters(ctx, pagination, filters)
	}

	cacheKey := fmt.Sprintf("%s:%s", getMastersPrefix, paginationHash)
	mastersToDecode, err := c.cacheProvider.Get(ctx, cacheKey)
	if err == nil {
		var masters []entities.Master
		if err = json.Unmarshal([]byte(mastersToDecode), &masters); err == nil {
			return masters, nil
		}
	}

	logging.LogErrorContext(
		ctx,
		c.logger,
		"Failed to get cached Masters",
		err,
	)

	masters, err := c.UseCases.GetMasters(ctx, pagination, filters)
	if err != nil {
		return nil, err
	}

	encodedMasters, err := json.Marshal(masters)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			c.logger,
			"Failed to cache Masters",
			err,
		)

		return masters, nil
	}

	if err = c.cacheProvider.Set(ctx, cacheKey, encodedMasters, getMastersTTL); err != nil {
		logging.LogErrorContext(
			ctx,
			c.logger,
			"Failed to cache Masters",
			err,
		)
	}

	return masters, nil
}

func (c *CacheDecorator) CountMasters(ctx context.Context, filters *entities.MastersFilters) (uint64, error) {
	if _, err := c.cacheProvider.Ping(ctx); err != nil {
		logging.LogErrorContext(
			ctx,
			c.logger,
			"Cache provider error",
			err,
		)

		return c.UseCases.CountMasters(ctx, filters)
	}

	filtersHash, err := rxhash.HashStruct(filters)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			c.logger,
			"Failed to get cached Masters counter",
			err,
		)

		return c.UseCases.CountMasters(ctx, filters)
	}

	cacheKey := fmt.Sprintf("%s:%s", countMastersPrefix, filtersHash)
	strCounter, err := c.cacheProvider.Get(ctx, cacheKey)
	if err == nil {
		if counter, err := strconv.ParseUint(strCounter, 10, 64); err == nil {
			return counter, nil
		}
	}

	logging.LogErrorContext(
		ctx,
		c.logger,
		"Failed to get cached Masters counter",
		err,
	)

	counter, err := c.UseCases.CountMasters(ctx, filters)
	if err != nil {
		return 0, err
	}

	if err = c.cacheProvider.Set(ctx, cacheKey, counter, countMastersTTL); err != nil {
		logging.LogErrorContext(
			ctx,
			c.logger,
			"Failed to cache Masters counter",
			err,
		)
	}

	return counter, nil
}

func (c *CacheDecorator) GetMasterByID(ctx context.Context, id uint64) (*entities.Master, error) {
	if _, err := c.cacheProvider.Ping(ctx); err != nil {
		logging.LogErrorContext(
			ctx,
			c.logger,
			"Cache provider error",
			err,
		)

		return c.UseCases.GetMasterByID(ctx, id)
	}

	cacheKey := fmt.Sprintf("%s:%d", getMasterByIDPrefix, id)
	masterToDecode, err := c.cacheProvider.Get(ctx, cacheKey)
	if err == nil {
		var master entities.Master
		if err = json.Unmarshal([]byte(masterToDecode), &master); err == nil {
			return &master, nil
		}
	}

	logging.LogErrorContext(
		ctx,
		c.logger,
		fmt.Sprintf("Failed to get cached Master with id=%d", id),
		err,
	)

	master, err := c.UseCases.GetMasterByID(ctx, id)
	if err != nil {
		return nil, err
	}

	encodedMaster, err := json.Marshal(master)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			c.logger,
			fmt.Sprintf("Failed to cache Master with id=%d", id),
			err,
		)

		return master, nil
	}

	if err = c.cacheProvider.Set(ctx, cacheKey, encodedMaster, getMasterByIDTTL); err != nil {
		logging.LogErrorContext(
			ctx,
			c.logger,
			fmt.Sprintf("Failed to cache Master with id=%d", id),
			err,
		)
	}

	return master, nil
}

func (c *CacheDecorator) GetMasterByUserID(ctx context.Context, userID uint64) (*entities.Master, error) {
	if _, err := c.cacheProvider.Ping(ctx); err != nil {
		logging.LogErrorContext(
			ctx,
			c.logger,
			"Cache provider error",
			err,
		)

		return c.UseCases.GetMasterByUserID(ctx, userID)
	}

	cacheKey := fmt.Sprintf("%s:%d", getMasterByUserIDPrefix, userID)
	masterToDecode, err := c.cacheProvider.Get(ctx, cacheKey)
	if err == nil {
		var master entities.Master
		if err = json.Unmarshal([]byte(masterToDecode), &master); err == nil {
			return &master, nil
		}
	}

	logging.LogErrorContext(
		ctx,
		c.logger,
		fmt.Sprintf("Failed to get cached Master with userID=%d", userID),
		err,
	)

	master, err := c.UseCases.GetMasterByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	encodedMaster, err := json.Marshal(master)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			c.logger,
			fmt.Sprintf("Failed to cache Master with userID=%d", userID),
			err,
		)

		return master, nil
	}

	if err = c.cacheProvider.Set(ctx, cacheKey, encodedMaster, getMasterByUserIDTTL); err != nil {
		logging.LogErrorContext(
			ctx,
			c.logger,
			fmt.Sprintf("Failed to cache Master with userID=%d", userID),
			err,
		)
	}

	return master, nil
}

func (c *CacheDecorator) GetAllCategories(ctx context.Context) ([]entities.Category, error) {
	if _, err := c.cacheProvider.Ping(ctx); err != nil {
		logging.LogErrorContext(
			ctx,
			c.logger,
			"Cache provider error",
			err,
		)

		return c.UseCases.GetAllCategories(ctx)
	}

	cacheKey := getCategoriesPrefix
	categoriesToDecode, err := c.cacheProvider.Get(ctx, cacheKey)
	if err == nil {
		var categories []entities.Category
		if err = json.Unmarshal([]byte(categoriesToDecode), &categories); err == nil {
			return categories, nil
		}
	}

	logging.LogErrorContext(
		ctx,
		c.logger,
		"Failed to get cached Categories",
		err,
	)

	categories, err := c.UseCases.GetAllCategories(ctx)
	if err != nil {
		return nil, err
	}

	encodedCategories, err := json.Marshal(categories)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			c.logger,
			"Failed to cache Categories",
			err,
		)

		return categories, nil
	}

	if err = c.cacheProvider.Set(ctx, cacheKey, encodedCategories, getCategoriesTTL); err != nil {
		logging.LogErrorContext(
			ctx,
			c.logger,
			"Failed to cache Categories",
			err,
		)
	}

	return categories, nil
}

func (c *CacheDecorator) GetCategoryByID(ctx context.Context, id uint32) (*entities.Category, error) {
	if _, err := c.cacheProvider.Ping(ctx); err != nil {
		logging.LogErrorContext(
			ctx,
			c.logger,
			"Cache provider error",
			err,
		)

		return c.UseCases.GetCategoryByID(ctx, id)
	}

	cacheKey := fmt.Sprintf("%s:%d", getCategoryByIDPrefix, id)
	categoryToDecode, err := c.cacheProvider.Get(ctx, cacheKey)
	if err == nil {
		var category entities.Category
		if err = json.Unmarshal([]byte(categoryToDecode), &category); err == nil {
			return &category, nil
		}
	}

	logging.LogErrorContext(
		ctx,
		c.logger,
		fmt.Sprintf("Failed to get cached Category with id=%d", id),
		err,
	)

	category, err := c.UseCases.GetCategoryByID(ctx, id)
	if err != nil {
		return nil, err
	}

	encodedCategory, err := json.Marshal(category)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			c.logger,
			fmt.Sprintf("Failed to cache Category with id=%d", id),
			err,
		)

		return category, nil
	}

	if err = c.cacheProvider.Set(ctx, cacheKey, encodedCategory, getCategoryByIDTTL); err != nil {
		logging.LogErrorContext(
			ctx,
			c.logger,
			fmt.Sprintf("Failed to cache Category with id=%d", id),
			err,
		)
	}

	return category, nil
}

func (c *CacheDecorator) GetAllTags(ctx context.Context) ([]entities.Tag, error) {
	if _, err := c.cacheProvider.Ping(ctx); err != nil {
		logging.LogErrorContext(
			ctx,
			c.logger,
			"Cache provider error",
			err,
		)

		return c.UseCases.GetAllTags(ctx)
	}

	cacheKey := getTagsPrefix
	tagsToDecode, err := c.cacheProvider.Get(ctx, cacheKey)
	if err == nil {
		var tags []entities.Tag
		if err = json.Unmarshal([]byte(tagsToDecode), &tags); err == nil {
			return tags, nil
		}
	}

	logging.LogErrorContext(
		ctx,
		c.logger,
		"Failed to get cached Tags",
		err,
	)

	tags, err := c.UseCases.GetAllTags(ctx)
	if err != nil {
		return nil, err
	}

	encodedTags, err := json.Marshal(tags)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			c.logger,
			"Failed to cache Tags",
			err,
		)

		return tags, nil
	}

	if err = c.cacheProvider.Set(ctx, cacheKey, encodedTags, getTagsTTL); err != nil {
		logging.LogErrorContext(
			ctx,
			c.logger,
			"Failed to cache Tags",
			err,
		)
	}

	return tags, nil
}

func (c *CacheDecorator) GetTagByID(ctx context.Context, id uint32) (*entities.Tag, error) {
	if _, err := c.cacheProvider.Ping(ctx); err != nil {
		logging.LogErrorContext(
			ctx,
			c.logger,
			"Cache provider error",
			err,
		)

		return c.UseCases.GetTagByID(ctx, id)
	}

	cacheKey := fmt.Sprintf("%s:%d", getTagByIDPrefix, id)
	tagToDecode, err := c.cacheProvider.Get(ctx, cacheKey)
	if err == nil {
		var tag entities.Tag
		if err = json.Unmarshal([]byte(tagToDecode), &tag); err == nil {
			return &tag, nil
		}
	}

	logging.LogErrorContext(
		ctx,
		c.logger,
		fmt.Sprintf("Failed to get cached Tag with id=%d", id),
		err,
	)

	tag, err := c.UseCases.GetTagByID(ctx, id)
	if err != nil {
		return nil, err
	}

	encodedTag, err := json.Marshal(tag)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			c.logger,
			fmt.Sprintf("Failed to cache Tag with id=%d", id),
			err,
		)

		return tag, nil
	}

	if err = c.cacheProvider.Set(ctx, cacheKey, encodedTag, getTagByIDTTL); err != nil {
		logging.LogErrorContext(
			ctx,
			c.logger,
			fmt.Sprintf("Failed to cache Tag with id=%d", id),
			err,
		)
	}

	return tag, nil
}

func (c *CacheDecorator) GetTicketByID(ctx context.Context, id uint64) (*entities.Ticket, error) {
	if _, err := c.cacheProvider.Ping(ctx); err != nil {
		logging.LogErrorContext(
			ctx,
			c.logger,
			"Cache provider error",
			err,
		)

		return c.UseCases.GetTicketByID(ctx, id)
	}

	cacheKey := fmt.Sprintf("%s:%d", getTicketByIDPrefix, id)
	ticketToDecode, err := c.cacheProvider.Get(ctx, cacheKey)
	if err == nil {
		var ticket entities.Ticket
		if err = json.Unmarshal([]byte(ticketToDecode), &ticket); err == nil {
			return &ticket, nil
		}
	}

	logging.LogErrorContext(
		ctx,
		c.logger,
		fmt.Sprintf("Failed to get cached Ticket with id=%d", id),
		err,
	)

	ticket, err := c.UseCases.GetTicketByID(ctx, id)
	if err != nil {
		return nil, err
	}

	encodedTicket, err := json.Marshal(ticket)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			c.logger,
			fmt.Sprintf("Failed to cache Ticket with id=%d", id),
			err,
		)

		return ticket, nil
	}

	if err = c.cacheProvider.Set(ctx, cacheKey, encodedTicket, getTicketByIDTTL); err != nil {
		logging.LogErrorContext(
			ctx,
			c.logger,
			fmt.Sprintf("Failed to cache Ticket with id=%d", id),
			err,
		)
	}

	return ticket, nil
}

func (c *CacheDecorator) GetTickets(
	ctx context.Context,
	pagination *entities.Pagination,
	filters *entities.TicketsFilters,
) ([]entities.Ticket, error) {
	if _, err := c.cacheProvider.Ping(ctx); err != nil {
		logging.LogErrorContext(
			ctx,
			c.logger,
			"Cache provider error",
			err,
		)

		return c.UseCases.GetTickets(ctx, pagination, filters)
	}

	paginationHash, err := rxhash.HashStruct(pagination)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			c.logger,
			"Failed to get cached Tickets",
			err,
		)

		return c.UseCases.GetTickets(ctx, pagination, filters)
	}

	filtersHash, err := rxhash.HashStruct(filters)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			c.logger,
			"Failed to get cached Tickets",
			err,
		)

		return c.UseCases.GetTickets(ctx, pagination, filters)
	}

	cacheKey := fmt.Sprintf(
		"%s:%s_%s",
		getTicketsPrefix,
		paginationHash,
		filtersHash,
	)

	ticketsToDecode, err := c.cacheProvider.Get(ctx, cacheKey)
	if err == nil {
		var tickets []entities.Ticket
		if err = json.Unmarshal([]byte(ticketsToDecode), &tickets); err == nil {
			return tickets, nil
		}
	}

	logging.LogErrorContext(
		ctx,
		c.logger,
		"Failed to get cached Tickets",
		err,
	)

	tickets, err := c.UseCases.GetTickets(ctx, pagination, filters)
	if err != nil {
		return nil, err
	}

	encodedTickets, err := json.Marshal(tickets)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			c.logger,
			"Failed to cache Tickets",
			err,
		)

		return tickets, nil
	}

	if err = c.cacheProvider.Set(ctx, cacheKey, encodedTickets, getTicketsTTL); err != nil {
		logging.LogErrorContext(
			ctx,
			c.logger,
			"Failed to cache Tickets",
			err,
		)
	}

	return tickets, nil
}

func (c *CacheDecorator) GetUserTickets(
	ctx context.Context,
	userID uint64,
	pagination *entities.Pagination,
	filters *entities.TicketsFilters,
) ([]entities.Ticket, error) {
	if _, err := c.cacheProvider.Ping(ctx); err != nil {
		logging.LogErrorContext(
			ctx,
			c.logger,
			"Cache provider error",
			err,
		)

		return c.UseCases.GetUserTickets(ctx, userID, pagination, filters)
	}

	paginationHash, err := rxhash.HashStruct(pagination)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			c.logger,
			fmt.Sprintf("Failed to get cached Tickets for User with id=%d", userID),
			err,
		)

		return c.UseCases.GetUserTickets(ctx, userID, pagination, filters)
	}

	filtersHash, err := rxhash.HashStruct(filters)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			c.logger,
			fmt.Sprintf("Failed to get cached Tickets for User with id=%d", userID),
			err,
		)

		return c.UseCases.GetUserTickets(ctx, userID, pagination, filters)
	}

	cacheKey := fmt.Sprintf(
		"%s:%d_%s_%s",
		getUserTicketsPrefix,
		userID,
		paginationHash,
		filtersHash,
	)

	ticketsToDecode, err := c.cacheProvider.Get(ctx, cacheKey)
	if err == nil {
		var tickets []entities.Ticket
		if err = json.Unmarshal([]byte(ticketsToDecode), &tickets); err == nil {
			return tickets, nil
		}
	}

	logging.LogErrorContext(
		ctx,
		c.logger,
		fmt.Sprintf("Failed to get cached Tickets for User with id=%d", userID),
		err,
	)

	tickets, err := c.UseCases.GetUserTickets(ctx, userID, pagination, filters)
	if err != nil {
		return nil, err
	}

	encodedTickets, err := json.Marshal(tickets)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			c.logger,
			fmt.Sprintf("Failed to cache Tickets for User with id=%d", userID),
			err,
		)

		return tickets, nil
	}

	if err = c.cacheProvider.Set(ctx, cacheKey, encodedTickets, getUserTicketsTTL); err != nil {
		logging.LogErrorContext(
			ctx,
			c.logger,
			fmt.Sprintf("Failed to cache Tickets for User with id=%d", userID),
			err,
		)
	}

	return tickets, nil
}

func (c *CacheDecorator) CountTickets(ctx context.Context, filters *entities.TicketsFilters) (uint64, error) {
	if _, err := c.cacheProvider.Ping(ctx); err != nil {
		logging.LogErrorContext(
			ctx,
			c.logger,
			"Cache provider error",
			err,
		)

		return c.UseCases.CountTickets(ctx, filters)
	}

	filtersHash, err := rxhash.HashStruct(filters)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			c.logger,
			"Failed to get cached Tickets counter",
			err,
		)

		return c.UseCases.CountTickets(ctx, filters)
	}

	cacheKey := fmt.Sprintf("%s:%s", countTicketsPrefix, filtersHash)
	strCounter, err := c.cacheProvider.Get(ctx, cacheKey)
	if err == nil {
		if counter, err := strconv.ParseUint(strCounter, 10, 64); err == nil {
			return counter, nil
		}
	}

	logging.LogErrorContext(
		ctx,
		c.logger,
		"Failed to get cached Tickets counter",
		err,
	)

	counter, err := c.UseCases.CountTickets(ctx, filters)
	if err != nil {
		return 0, err
	}

	if err = c.cacheProvider.Set(ctx, cacheKey, counter, countTicketsTTL); err != nil {
		logging.LogErrorContext(
			ctx,
			c.logger,
			"Failed to cache Tickets counter",
			err,
		)
	}

	return counter, nil
}

func (c *CacheDecorator) CountUserTickets(
	ctx context.Context,
	userID uint64,
	filters *entities.TicketsFilters,
) (uint64, error) {
	if _, err := c.cacheProvider.Ping(ctx); err != nil {
		logging.LogErrorContext(
			ctx,
			c.logger,
			"Cache provider error",
			err,
		)

		return c.UseCases.CountUserTickets(ctx, userID, filters)
	}

	filtersHash, err := rxhash.HashStruct(filters)
	if err != nil {
		logging.LogErrorContext(
			ctx,
			c.logger,
			fmt.Sprintf("Failed to get cached Tickets counter for User with id=%d", userID),
			err,
		)

		return c.UseCases.CountUserTickets(ctx, userID, filters)
	}

	cacheKey := fmt.Sprintf(
		"%s:%d_%s",
		countUserTicketsPrefix,
		userID,
		filtersHash,
	)

	strCounter, err := c.cacheProvider.Get(ctx, cacheKey)
	if err == nil {
		if counter, err := strconv.ParseUint(strCounter, 10, 64); err == nil {
			return counter, nil
		}
	}

	logging.LogErrorContext(
		ctx,
		c.logger,
		fmt.Sprintf("Failed to get cached Tickets counter for User with id=%d", userID),
		err,
	)

	counter, err := c.UseCases.CountUserTickets(ctx, userID, filters)
	if err != nil {
		return 0, err
	}

	if err = c.cacheProvider.Set(ctx, cacheKey, counter, countUserTicketsTTL); err != nil {
		logging.LogErrorContext(
			ctx,
			c.logger,
			fmt.Sprintf("Failed to cache Tickets counter for User with id=%d", userID),
			err,
		)
	}

	return counter, nil
}

func (c *CacheDecorator) UpdateUserProfile(
	ctx context.Context,
	userToDecodeProfileData entities.RawUpdateUserProfileDTO,
) error {
	if _, err := c.cacheProvider.Ping(ctx); err != nil {
		logging.LogErrorContext(
			ctx,
			c.logger,
			"Cache provider error",
			err,
		)

		return c.UseCases.UpdateUserProfile(ctx, userToDecodeProfileData)
	}

	err := c.UseCases.UpdateUserProfile(ctx, userToDecodeProfileData)
	if err == nil {
		var user *entities.User
		user, err = c.UseCases.GetMe(ctx, userToDecodeProfileData.AccessToken)
		if err != nil {
			logging.LogErrorContext(
				ctx,
				c.logger,
				fmt.Sprintf("Failed to get User with AccessToken=%s to delete cache", userToDecodeProfileData.AccessToken),
				err,
			)

			return nil
		}

		patterns := []string{
			fmt.Sprintf("%s:%d", getUserByIDPrefix, user.ID),
			fmt.Sprintf("%s:%s", getUserByEmailPrefix, user.Email),
		}

		for _, pattern := range patterns {
			if cacheErr := c.cacheProvider.DelByPattern(ctx, pattern, nil); cacheErr != nil {
				logging.LogErrorContext(
					ctx,
					c.logger,
					fmt.Sprintf("Failed to delete cache by pattern %s", pattern),
					err,
				)
			}
		}
	}

	return err
}

func (c *CacheDecorator) UpdateToy(
	ctx context.Context,
	rawToyData entities.RawUpdateToyDTO,
) error {
	if _, err := c.cacheProvider.Ping(ctx); err != nil {
		logging.LogErrorContext(
			ctx,
			c.logger,
			"Cache provider error",
			err,
		)

		return c.UseCases.UpdateToy(ctx, rawToyData)
	}

	err := c.UseCases.UpdateToy(ctx, rawToyData)
	if err == nil {
		var toy *entities.Toy
		toy, err = c.UseCases.GetToyByID(ctx, rawToyData.ID)
		if err != nil {
			logging.LogErrorContext(
				ctx,
				c.logger,
				fmt.Sprintf("Failed to get Toy with ID=%d to delete cache", rawToyData.ID),
				err,
			)

			return nil
		}

		patterns := []string{
			fmt.Sprintf("%s:%d", getToyByIDPrefix, toy.ID),
			fmt.Sprintf("%s:%d*", getMasterToysPrefix, toy.MasterID),
			fmt.Sprintf("%s:%d*", countMasterToysPrefix, toy.MasterID),
		}

		for _, pattern := range patterns {
			if cacheErr := c.cacheProvider.DelByPattern(ctx, pattern, nil); cacheErr != nil {
				logging.LogErrorContext(
					ctx,
					c.logger,
					fmt.Sprintf("Failed to delete cache by pattern %s", pattern),
					err,
				)
			}
		}
	}

	return err
}

func (c *CacheDecorator) DeleteToy(ctx context.Context, accessToken string, id uint64) error {
	if _, err := c.cacheProvider.Ping(ctx); err != nil {
		logging.LogErrorContext(
			ctx,
			c.logger,
			"Cache provider error",
			err,
		)

		return c.UseCases.DeleteToy(ctx, accessToken, id)
	}

	err := c.UseCases.DeleteToy(ctx, accessToken, id)
	if err == nil {
		var toy *entities.Toy
		toy, err = c.UseCases.GetToyByID(ctx, id)
		if err != nil {
			logging.LogErrorContext(
				ctx,
				c.logger,
				fmt.Sprintf("Failed to get Toy with ID=%d to delete cache", id),
				err,
			)

			return nil
		}

		patterns := []string{
			fmt.Sprintf("%s:%d", getToyByIDPrefix, toy.ID),
			fmt.Sprintf("%s:%d*", getMasterToysPrefix, toy.MasterID),
			fmt.Sprintf("%s:%d*", countMasterToysPrefix, toy.MasterID),
		}

		for _, pattern := range patterns {
			if cacheErr := c.cacheProvider.DelByPattern(ctx, pattern, nil); cacheErr != nil {
				logging.LogErrorContext(
					ctx,
					c.logger,
					fmt.Sprintf("Failed to delete cache by pattern %s", pattern),
					err,
				)
			}
		}
	}

	return err
}

func (c *CacheDecorator) UpdateMaster(
	ctx context.Context,
	rawMasterData entities.RawUpdateMasterDTO,
) error {
	if _, err := c.cacheProvider.Ping(ctx); err != nil {
		logging.LogErrorContext(
			ctx,
			c.logger,
			"Cache provider error",
			err,
		)

		return c.UseCases.UpdateMaster(ctx, rawMasterData)
	}

	err := c.UseCases.UpdateMaster(ctx, rawMasterData)
	if err == nil {
		var master *entities.Master
		master, err = c.UseCases.GetMasterByID(ctx, rawMasterData.ID)
		if err != nil {
			logging.LogErrorContext(
				ctx,
				c.logger,
				fmt.Sprintf("Failed to get Master with ID=%d to delete cache", rawMasterData.ID),
				err,
			)

			return nil
		}

		patterns := []string{
			fmt.Sprintf("%s:%d", getMasterByIDPrefix, master.ID),
			fmt.Sprintf("%s:%d", getMasterByUserIDPrefix, master.UserID),
			fmt.Sprintf("%s:%d*", getMasterToysPrefix, master.ID),
			fmt.Sprintf("%s:%d*", countMasterToysPrefix, master.ID),
			fmt.Sprintf("%s*", getMastersPrefix),
		}

		for _, pattern := range patterns {
			if cacheErr := c.cacheProvider.DelByPattern(ctx, pattern, nil); cacheErr != nil {
				logging.LogErrorContext(
					ctx,
					c.logger,
					fmt.Sprintf("Failed to delete cache by pattern %s", pattern),
					err,
				)
			}
		}
	}

	return err
}

func (c *CacheDecorator) UpdateTicket(
	ctx context.Context,
	rawTicketData entities.RawUpdateTicketDTO,
) error {
	if _, err := c.cacheProvider.Ping(ctx); err != nil {
		logging.LogErrorContext(
			ctx,
			c.logger,
			"Cache provider error",
			err,
		)

		return c.UseCases.UpdateTicket(ctx, rawTicketData)
	}

	err := c.UseCases.UpdateTicket(ctx, rawTicketData)
	if err == nil {
		var ticket *entities.Ticket
		ticket, err = c.UseCases.GetTicketByID(ctx, rawTicketData.ID)
		if err != nil {
			logging.LogErrorContext(
				ctx,
				c.logger,
				fmt.Sprintf("Failed to get Ticket with ID=%d to delete cache", rawTicketData.ID),
				err,
			)

			return nil
		}

		patterns := []string{
			fmt.Sprintf("%s:%d", getTicketByIDPrefix, ticket.ID),
			fmt.Sprintf("%s:%d", getUserTicketsPrefix, ticket.UserID),
			fmt.Sprintf("%s:%d*", countUserTicketsPrefix, ticket.UserID),
		}

		for _, pattern := range patterns {
			if cacheErr := c.cacheProvider.DelByPattern(ctx, pattern, nil); cacheErr != nil {
				logging.LogErrorContext(
					ctx,
					c.logger,
					fmt.Sprintf("Failed to delete cache by pattern %s", pattern),
					err,
				)
			}
		}
	}

	return err
}

func (c *CacheDecorator) DeleteTicket(ctx context.Context, accessToken string, id uint64) error {
	if _, err := c.cacheProvider.Ping(ctx); err != nil {
		logging.LogErrorContext(
			ctx,
			c.logger,
			"Cache provider error",
			err,
		)

		return c.UseCases.DeleteTicket(ctx, accessToken, id)
	}

	err := c.UseCases.DeleteTicket(ctx, accessToken, id)
	if err == nil {
		var ticket *entities.Ticket
		ticket, err = c.UseCases.GetTicketByID(ctx, id)
		if err != nil {
			logging.LogErrorContext(
				ctx,
				c.logger,
				fmt.Sprintf("Failed to get Ticket with ID=%d to delete cache", id),
				err,
			)

			return nil
		}

		patterns := []string{
			fmt.Sprintf("%s:%d", getTicketByIDPrefix, ticket.ID),
			fmt.Sprintf("%s:%d", getUserTicketsPrefix, ticket.UserID),
			fmt.Sprintf("%s:%d*", countUserTicketsPrefix, ticket.UserID),
		}

		for _, pattern := range patterns {
			if cacheErr := c.cacheProvider.DelByPattern(ctx, pattern, nil); cacheErr != nil {
				logging.LogErrorContext(
					ctx,
					c.logger,
					fmt.Sprintf("Failed to delete cache by pattern %s", pattern),
					err,
				)
			}
		}
	}

	return err
}
