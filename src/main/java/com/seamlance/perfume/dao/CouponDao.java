package com.seamlance.perfume.dao;

import com.seamlance.perfume.entity.Coupon;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

@Repository
public interface CouponDao extends JpaRepository<Coupon, String> {
    Coupon findByCode(String code);
}
