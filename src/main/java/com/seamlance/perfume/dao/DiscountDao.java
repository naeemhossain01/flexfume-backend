package com.seamlance.perfume.dao;

import com.seamlance.perfume.entity.Discount;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

@Repository
public interface DiscountDao extends JpaRepository<Discount, String> {
}
