package com.laioffer.onlineorder.service;

import com.laioffer.onlineorder.entity.CartEntity;
import com.laioffer.onlineorder.entity.MenuItemEntity;
import com.laioffer.onlineorder.entity.OrderItemEntity;
import com.laioffer.onlineorder.exception.ResourceNotFoundException;
import com.laioffer.onlineorder.model.CartDto;
import com.laioffer.onlineorder.model.OrderItemDto;
import com.laioffer.onlineorder.repository.CartRepository;
import com.laioffer.onlineorder.repository.MenuItemRepository;
import com.laioffer.onlineorder.repository.OrderItemRepository;
import java.math.BigDecimal;
import java.util.List;
import java.util.Map;
import java.util.function.Function;
import java.util.stream.Collectors;
import org.springframework.cache.annotation.CacheEvict;
import org.springframework.cache.annotation.Cacheable;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

@Service
public class CartService {

    private final CartRepository cartRepository;
    private final MenuItemRepository menuItemRepository;
    private final OrderItemRepository orderItemRepository;

    public CartService(
            CartRepository cartRepository,
            MenuItemRepository menuItemRepository,
            OrderItemRepository orderItemRepository) {
        this.cartRepository = cartRepository;
        this.menuItemRepository = menuItemRepository;
        this.orderItemRepository = orderItemRepository;
    }

    @CacheEvict(cacheNames = "carts", key = "#customerId")
    @Transactional
    public void addMenuItemToCart(long customerId, long menuItemId) {
        CartEntity cart = getLockedCart(customerId);
        MenuItemEntity menuItem = menuItemRepository.findById(menuItemId)
                .orElseThrow(() -> new ResourceNotFoundException(
                        "Menu item " + menuItemId + " was not found"));

        OrderItemEntity orderItem = orderItemRepository
                .findByCartIdAndMenuItemId(cart.id(), menuItemId)
                .map(existing -> new OrderItemEntity(
                        existing.id(), menuItemId, cart.id(), menuItem.price(),
                        existing.quantity() + 1))
                .orElseGet(() -> new OrderItemEntity(
                        null, menuItemId, cart.id(), menuItem.price(), 1));

        orderItemRepository.save(orderItem);
        cartRepository.updateTotalPrice(cart.id(), cart.totalPrice().add(menuItem.price()));
    }

    @Cacheable(cacheNames = "carts", key = "#customerId")
    @Transactional(readOnly = true)
    public CartDto getCart(Long customerId) {
        CartEntity cart = cartRepository.findByCustomerId(customerId)
                .orElseThrow(() -> new ResourceNotFoundException("Cart was not found"));
        List<OrderItemEntity> orderItems = orderItemRepository.findAllByCartIdOrderByIdAsc(cart.id());

        Map<Long, MenuItemEntity> menuItems = menuItemRepository
                .findAllById(orderItems.stream().map(OrderItemEntity::menuItemId).toList())
                .stream()
                .collect(Collectors.toMap(MenuItemEntity::id, Function.identity()));

        List<OrderItemDto> itemDtos = orderItems.stream()
                .map(orderItem -> {
                    MenuItemEntity menuItem = menuItems.get(orderItem.menuItemId());
                    if (menuItem == null) {
                        throw new ResourceNotFoundException(
                                "Menu item " + orderItem.menuItemId() + " was not found");
                    }
                    return new OrderItemDto(orderItem, menuItem);
                })
                .toList();
        return new CartDto(cart, itemDtos);
    }

    @CacheEvict(cacheNames = "carts", key = "#customerId")
    @Transactional
    public void clearCart(Long customerId) {
        CartEntity cart = getLockedCart(customerId);
        orderItemRepository.deleteByCartId(cart.id());
        cartRepository.updateTotalPrice(cart.id(), BigDecimal.ZERO);
    }

    private CartEntity getLockedCart(Long customerId) {
        return cartRepository.findLockedByCustomerId(customerId)
                .orElseThrow(() -> new ResourceNotFoundException("Cart was not found"));
    }
}
