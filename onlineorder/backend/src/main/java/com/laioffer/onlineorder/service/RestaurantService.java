package com.laioffer.onlineorder.service;

import com.laioffer.onlineorder.entity.MenuItemEntity;
import com.laioffer.onlineorder.model.MenuItemDto;
import com.laioffer.onlineorder.model.RestaurantDto;
import com.laioffer.onlineorder.repository.MenuItemRepository;
import com.laioffer.onlineorder.repository.RestaurantRepository;
import java.util.List;
import java.util.Map;
import java.util.stream.Collectors;
import org.springframework.cache.annotation.Cacheable;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

@Service
public class RestaurantService {

    private final MenuItemRepository menuItemRepository;
    private final RestaurantRepository restaurantRepository;

    public RestaurantService(
            RestaurantRepository restaurantRepository,
            MenuItemRepository menuItemRepository) {
        this.restaurantRepository = restaurantRepository;
        this.menuItemRepository = menuItemRepository;
    }

    @Cacheable("restaurants")
    @Transactional(readOnly = true)
    public List<RestaurantDto> getRestaurants() {
        Map<Long, List<MenuItemDto>> menuItemsByRestaurant = menuItemRepository.findAll().stream()
                .collect(Collectors.groupingBy(
                        MenuItemEntity::restaurantId,
                        Collectors.mapping(MenuItemDto::new, Collectors.toList())));

        return restaurantRepository.findAll().stream()
                .map(restaurant -> new RestaurantDto(
                        restaurant,
                        menuItemsByRestaurant.getOrDefault(restaurant.id(), List.of())))
                .toList();
    }
}
