package com.laioffer.onlineorder.entity;

import org.springframework.data.annotation.Id;
import org.springframework.data.relational.core.mapping.Column;
import org.springframework.data.relational.core.mapping.Table;

@Table("restaurants")
public record RestaurantEntity(
        @Id @Column("id") Long id,
        @Column("name") String name,
        @Column("address") String address,
        @Column("phone") String phone,
        @Column("image_url") String imageUrl) {
}
