package com.laioffer.onlineorder.entity;

import java.math.BigDecimal;
import org.springframework.data.annotation.Id;
import org.springframework.data.relational.core.mapping.Column;
import org.springframework.data.relational.core.mapping.Table;

@Table("carts")
public record CartEntity(
        @Id @Column("id") Long id,
        @Column("customer_id") Long customerId,
        @Column("total_price") BigDecimal totalPrice) {
}
