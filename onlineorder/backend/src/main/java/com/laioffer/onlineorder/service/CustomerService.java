package com.laioffer.onlineorder.service;

import com.laioffer.onlineorder.entity.CartEntity;
import com.laioffer.onlineorder.entity.CustomerEntity;
import com.laioffer.onlineorder.exception.EmailAlreadyRegisteredException;
import com.laioffer.onlineorder.exception.ResourceNotFoundException;
import com.laioffer.onlineorder.repository.CartRepository;
import com.laioffer.onlineorder.repository.CustomerRepository;
import java.math.BigDecimal;
import java.util.Locale;
import org.springframework.dao.DataIntegrityViolationException;
import org.springframework.security.core.userdetails.User;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.security.crypto.password.PasswordEncoder;
import org.springframework.security.provisioning.UserDetailsManager;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

@Service
public class CustomerService {

    private final CartRepository cartRepository;
    private final CustomerRepository customerRepository;
    private final PasswordEncoder passwordEncoder;
    private final UserDetailsManager userDetailsManager;

    public CustomerService(
            CartRepository cartRepository,
            CustomerRepository customerRepository,
            PasswordEncoder passwordEncoder,
            UserDetailsManager userDetailsManager) {
        this.cartRepository = cartRepository;
        this.customerRepository = customerRepository;
        this.passwordEncoder = passwordEncoder;
        this.userDetailsManager = userDetailsManager;
    }

    @Transactional
    public void signUp(String email, String password, String firstName, String lastName) {
        String normalizedEmail = normalizeEmail(email);
        if (userDetailsManager.userExists(normalizedEmail)) {
            throw new EmailAlreadyRegisteredException();
        }

        UserDetails user = User.builder()
                .username(normalizedEmail)
                .password(passwordEncoder.encode(password))
                .roles("USER")
                .build();

        try {
            userDetailsManager.createUser(user);
            int updated = customerRepository.updateNameByEmail(
                    normalizedEmail, firstName.trim(), lastName.trim());
            if (updated != 1) {
                throw new IllegalStateException("Could not finish creating the customer profile");
            }
            CustomerEntity customer = getCustomerByEmail(normalizedEmail);
            cartRepository.save(new CartEntity(null, customer.id(), BigDecimal.ZERO));
        } catch (DataIntegrityViolationException exception) {
            throw new EmailAlreadyRegisteredException();
        }
    }

    @Transactional(readOnly = true)
    public CustomerEntity getCustomerByEmail(String email) {
        return customerRepository.findByEmail(normalizeEmail(email))
                .orElseThrow(() -> new ResourceNotFoundException("Customer was not found"));
    }

    private String normalizeEmail(String email) {
        return email.trim().toLowerCase(Locale.ROOT);
    }
}
