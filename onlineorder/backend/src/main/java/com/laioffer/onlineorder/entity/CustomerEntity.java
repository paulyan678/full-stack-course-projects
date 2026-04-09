package com.laioffer.onlineorder.entity;

import org.springframework.data.annotation.Id;
import org.springframework.data.relational.core.mapping.Column;
import org.springframework.data.relational.core.mapping.Table;

@Table("customers")
public record CustomerEntity(
        @Id @Column("id") Long id,
        @Column("email") String email,
        @Column("password") String password,
        @Column("enabled") boolean enabled,
        @Column("first_name") String firstName,
        @Column("last_name") String lastName) {
}
