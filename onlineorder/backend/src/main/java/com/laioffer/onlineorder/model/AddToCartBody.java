package com.laioffer.onlineorder.model;

import jakarta.validation.constraints.NotNull;
import jakarta.validation.constraints.Positive;

public record AddToCartBody(
        @NotNull(message = "menu_id is required")
        @Positive(message = "menu_id must be positive")
        Long menuId) {
}
