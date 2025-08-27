package com.seamlance.perfume.dao;

import com.seamlance.perfume.entity.Order;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.JpaSpecificationExecutor;
import org.springframework.stereotype.Repository;

import java.util.List;
import java.util.Optional;

@Repository
public interface OrderDao extends JpaRepository<Order, String>, JpaSpecificationExecutor<Order> {
    List<Order> findByUserId(String userId);
}
