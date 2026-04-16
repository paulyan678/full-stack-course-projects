package com.laioffer.onlineorder.model;

import jakarta.validation.constraints.Email;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.Size;

public record RegisterBody(
        @NotBlank(message = "email is required")
        @Email(message = "email must be valid")
        @Size(max = 254, message = "email is too long")
        String email,

        @NotBlank(message = "password is required")
        @Size(min = 8, max = 72, message = "password must contain 8 to 72 characters")
        String password,

        @NotBlank(message = "first_name is required")
        @Size(max = 80, message = "first_name is too long")
        String firstName,

        @NotBlank(message = "last_name is required")
        @Size(max = 80, message = "last_name is too long")
        String lastName) {
}
