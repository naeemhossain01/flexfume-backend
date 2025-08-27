package com.seamlance.perfume.dao;

import com.seamlance.perfume.entity.Address;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.Optional;

@Repository
public interface AddressDao extends JpaRepository<Address, String> {
    Optional<Address> findByUserId(String userId);

    Address findByIdAndUserId(String addressId, String userId);
}
