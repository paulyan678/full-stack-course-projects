package com.laioffer.onlineorder.repository;

import com.laioffer.onlineorder.entity.CartEntity;
import java.math.BigDecimal;
import java.util.Optional;
import org.springframework.data.jdbc.repository.query.Modifying;
import org.springframework.data.jdbc.repository.query.Query;
import org.springframework.data.repository.ListCrudRepository;

public interface CartRepository extends ListCrudRepository<CartEntity, Long> {

    Optional<CartEntity> findByCustomerId(Long customerId);

    @Query("SELECT id, customer_id, total_price FROM carts WHERE customer_id = :customerId FOR UPDATE")
    Optional<CartEntity> findLockedByCustomerId(Long customerId);

    @Modifying
    @Query("UPDATE carts SET total_price = :totalPrice WHERE id = :cartId")
    int updateTotalPrice(Long cartId, BigDecimal totalPrice);
}
