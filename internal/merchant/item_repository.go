package merchant

import localError "belimang/pkg/error"

type IItemRepository interface {
	FindAll(params GetItemQueryParam) ([]Item, *localError.GlobalError)
}
