package com.laioffer.onlineorder.controller;

import com.laioffer.onlineorder.entity.CustomerEntity;
import com.laioffer.onlineorder.model.AddToCartBody;
import com.laioffer.onlineorder.model.CartDto;
import com.laioffer.onlineorder.service.CartService;
import com.laioffer.onlineorder.service.CustomerService;
import jakarta.validation.Valid;
import org.springframework.http.HttpStatus;
import org.springframework.security.core.Authentication;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.ResponseStatus;
import org.springframework.web.bind.annotation.RestController;

@RestController
public class CartController {

    private final CartService cartService;
    private final CustomerService customerService;

    public CartController(CartService cartService, CustomerService customerService) {
        this.cartService = cartService;
        this.customerService = customerService;
    }

    @GetMapping("/cart")
    public CartDto getCart(Authentication authentication) {
        CustomerEntity customer = currentCustomer(authentication);
        return cartService.getCart(customer.id());
    }

    @PostMapping("/cart")
    @ResponseStatus(HttpStatus.NO_CONTENT)
    public void addToCart(Authentication authentication, @Valid @RequestBody AddToCartBody body) {
        CustomerEntity customer = currentCustomer(authentication);
        cartService.addMenuItemToCart(customer.id(), body.menuId());
    }

    @PostMapping("/cart/checkout")
    @ResponseStatus(HttpStatus.NO_CONTENT)
    public void checkout(Authentication authentication) {
        CustomerEntity customer = currentCustomer(authentication);
        cartService.clearCart(customer.id());
    }

    private CustomerEntity currentCustomer(Authentication authentication) {
        return customerService.getCustomerByEmail(authentication.getName());
    }
}
