package com.laioffer.onlineorder.repository;

import com.laioffer.onlineorder.entity.MenuItemEntity;
import java.util.List;
import org.springframework.data.repository.ListCrudRepository;

public interface MenuItemRepository extends ListCrudRepository<MenuItemEntity, Long> {

    List<MenuItemEntity> findAllByRestaurantIdOrderByNameAsc(Long restaurantId);
}
