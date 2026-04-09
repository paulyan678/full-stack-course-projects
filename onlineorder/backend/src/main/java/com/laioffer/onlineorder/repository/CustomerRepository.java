package com.laioffer.onlineorder.repository;

import com.laioffer.onlineorder.entity.CustomerEntity;
import java.util.Optional;
import org.springframework.data.jdbc.repository.query.Modifying;
import org.springframework.data.jdbc.repository.query.Query;
import org.springframework.data.repository.ListCrudRepository;

public interface CustomerRepository extends ListCrudRepository<CustomerEntity, Long> {

    Optional<CustomerEntity> findByEmail(String email);

    @Modifying
    @Query("UPDATE customers SET first_name = :firstName, last_name = :lastName WHERE email = :email")
    int updateNameByEmail(String email, String firstName, String lastName);
}
