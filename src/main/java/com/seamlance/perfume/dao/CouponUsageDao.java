package com.seamlance.perfume.dao;

import java.math.BigDecimal;

import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.query.Param;
import org.springframework.stereotype.Repository;

import com.seamlance.perfume.entity.CouponUsage;

@Repository
public interface CouponUsageDao extends JpaRepository<CouponUsage, String> {

    @Query(nativeQuery = true, value = "SELECT cu.* FROM COUPON_USAGE cu LEFT JOIN COUPON c ON cu.COUPON_ID = c.COUPON_ID WHERE cu.USER_ID = :userId  and c.CODE = :code")
    CouponUsage findByCouponCodeAndUserId(@Param("code") String couponCode, @Param("userId") String userId);
    
    // Aggregation queries for usage statistics
    @Query("SELECT COALESCE(SUM(cu.usageCount), 0) FROM CouponUsage cu WHERE cu.coupon.id = :couponId")
    Integer getTotalUsageCountByCouponId(@Param("couponId") String couponId);
    
    @Query("SELECT COALESCE(SUM(cu.discountedAmount), 0) FROM CouponUsage cu WHERE cu.coupon.id = :couponId")
    BigDecimal getTotalSavingsAmountByCouponId(@Param("couponId") String couponId);
    
    @Query("SELECT COUNT(DISTINCT cu.user.id) FROM CouponUsage cu WHERE cu.coupon.id = :couponId")
    Integer getUniqueUsersCountByCouponId(@Param("couponId") String couponId);
    
    // Batch aggregation for multiple coupons (more efficient)
    @Query("SELECT cu.coupon.id, COALESCE(SUM(cu.usageCount), 0), COALESCE(SUM(cu.discountedAmount), 0), COUNT(DISTINCT cu.user.id) " +
           "FROM CouponUsage cu WHERE cu.coupon.id IN :couponIds GROUP BY cu.coupon.id")
    Object[][] getUsageStatisticsByCouponIds(@Param("couponIds") java.util.List<String> couponIds);
}
