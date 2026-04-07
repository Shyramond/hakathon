package service

import "errors"

var (
    // Auth / login
    ErrLoginCooldown = errors.New("login reward already claimed, try again later")

    // Shop / purchase
    ErrInsufficientBalance = errors.New("insufficient token balance")
    ErrBenefitNotActive    = errors.New("benefit is not active or available")
    ErrAlreadyOwned        = errors.New("benefit already owned")

    // Inventory
    ErrItemNotFound    = errors.New("inventory item not found")
    ErrNotItemOwner    = errors.New("you do not own this item")
    ErrAlreadyEquipped = errors.New("item is already in the requested state")

    // Generic
    ErrUserNotFound    = errors.New("user not found")
    ErrBenefitNotFound = errors.New("benefit not found")
)