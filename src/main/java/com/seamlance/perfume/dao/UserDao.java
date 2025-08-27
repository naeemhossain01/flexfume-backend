package com.seamlance.perfume.dao;

import com.seamlance.perfume.entity.User;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.Optional;

@Repository
public interface UserDao extends JpaRepository<User, String> {
    Optional<User> findByEmail(String email);
    Optional<User> findUserByPhoneNumber(String phoneNumber);
}
