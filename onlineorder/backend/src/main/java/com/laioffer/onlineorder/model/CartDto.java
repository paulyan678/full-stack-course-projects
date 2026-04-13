package com.laioffer.onlineorder.model;

import com.laioffer.onlineorder.entity.CartEntity;
import java.math.BigDecimal;
import java.util.List;

public record CartDto(
        Long id,
        BigDecimal totalPrice,
        List<OrderItemDto> orderItems) {

    public CartDto(CartEntity entity, List<OrderItemDto> orderItems) {
        this(entity.id(), entity.totalPrice(), List.copyOf(orderItems));
    }
}
