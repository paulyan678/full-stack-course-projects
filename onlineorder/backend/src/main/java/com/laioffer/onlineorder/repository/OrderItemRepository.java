package com.laioffer.onlineorder.repository;

import com.laioffer.onlineorder.entity.OrderItemEntity;
import java.util.List;
import java.util.Optional;
import org.springframework.data.jdbc.repository.query.Modifying;
import org.springframework.data.jdbc.repository.query.Query;
import org.springframework.data.repository.ListCrudRepository;

public interface OrderItemRepository extends ListCrudRepository<OrderItemEntity, Long> {

    List<OrderItemEntity> findAllByCartIdOrderByIdAsc(Long cartId);

    Optional<OrderItemEntity> findByCartIdAndMenuItemId(Long cartId, Long menuItemId);

    @Modifying
    @Query("DELETE FROM order_items WHERE cart_id = :cartId")
    int deleteByCartId(Long cartId);
}
