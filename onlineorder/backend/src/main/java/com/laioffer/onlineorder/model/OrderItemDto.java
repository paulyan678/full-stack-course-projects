package com.laioffer.onlineorder.model;

import com.laioffer.onlineorder.entity.MenuItemEntity;
import com.laioffer.onlineorder.entity.OrderItemEntity;
import java.math.BigDecimal;

public record OrderItemDto(
        Long orderItemId,
        Long menuItemId,
        Long restaurantId,
        BigDecimal price,
        Integer quantity,
        BigDecimal lineTotal,
        String menuItemName,
        String menuItemDescription,
        String menuItemImageUrl) {

    public OrderItemDto(OrderItemEntity orderItem, MenuItemEntity menuItem) {
        this(orderItem.id(), orderItem.menuItemId(), menuItem.restaurantId(), orderItem.price(),
                orderItem.quantity(),
                orderItem.price().multiply(BigDecimal.valueOf(orderItem.quantity())),
                menuItem.name(), menuItem.description(), menuItem.imageUrl());
    }
}
