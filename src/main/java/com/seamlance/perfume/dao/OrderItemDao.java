package com.seamlance.perfume.dao;

import com.seamlance.perfume.entity.OrderItem;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

@Repository
public interface OrderItemDao extends JpaRepository<OrderItem, String> {
}
