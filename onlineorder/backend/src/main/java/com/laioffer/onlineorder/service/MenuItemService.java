package com.laioffer.onlineorder.service;

import com.laioffer.onlineorder.exception.ResourceNotFoundException;
import com.laioffer.onlineorder.model.MenuItemDto;
import com.laioffer.onlineorder.repository.MenuItemRepository;
import com.laioffer.onlineorder.repository.RestaurantRepository;
import java.util.List;
import org.springframework.cache.annotation.Cacheable;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

@Service
public class MenuItemService {

    private final MenuItemRepository menuItemRepository;
    private final RestaurantRepository restaurantRepository;

    public MenuItemService(
            MenuItemRepository menuItemRepository,
            RestaurantRepository restaurantRepository) {
        this.menuItemRepository = menuItemRepository;
        this.restaurantRepository = restaurantRepository;
    }

    @Cacheable(cacheNames = "menus", key = "#restaurantId")
    @Transactional(readOnly = true)
    public List<MenuItemDto> getMenuItemsByRestaurantId(long restaurantId) {
        if (!restaurantRepository.existsById(restaurantId)) {
            throw new ResourceNotFoundException("Restaurant " + restaurantId + " was not found");
        }
        return menuItemRepository.findAllByRestaurantIdOrderByNameAsc(restaurantId).stream()
                .map(MenuItemDto::new)
                .toList();
    }
}
