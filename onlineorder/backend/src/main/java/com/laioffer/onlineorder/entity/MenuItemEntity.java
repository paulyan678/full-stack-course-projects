package com.laioffer.onlineorder.entity;

import java.math.BigDecimal;
import org.springframework.data.annotation.Id;
import org.springframework.data.relational.core.mapping.Column;
import org.springframework.data.relational.core.mapping.Table;

@Table("menu_items")
public record MenuItemEntity(
        @Id @Column("id") Long id,
        @Column("restaurant_id") Long restaurantId,
        @Column("name") String name,
        @Column("description") String description,
        @Column("price") BigDecimal price,
        @Column("image_url") String imageUrl) {
}
